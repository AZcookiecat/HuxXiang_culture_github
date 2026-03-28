package community

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MySQLRepository struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

type postRow struct {
	ID           int64     `db:"id"`
	Title        string    `db:"title"`
	Content      string    `db:"content"`
	AuthorID     int64     `db:"author_id"`
	AuthorName   string    `db:"author_username"`
	AuthorAvatar string    `db:"author_avatar"`
	AuthorBio    string    `db:"author_bio"`
	Category     string    `db:"category"`
	Status       string    `db:"status"`
	ViewCount    int64     `db:"view_count"`
	LikeCount    int64     `db:"like_count"`
	CommentCount int64     `db:"comment_count"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type commentRow struct {
	ID           int64         `db:"id"`
	Content      string        `db:"content"`
	AuthorID     int64         `db:"author_id"`
	AuthorName   string        `db:"author_username"`
	AuthorAvatar string        `db:"author_avatar"`
	PostID       int64         `db:"post_id"`
	ParentID     sql.NullInt64 `db:"parent_id"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

func NewMySQLRepository(writer, reader *sqlx.DB) *MySQLRepository {
	return &MySQLRepository{writer: writer, reader: reader}
}

func (r *MySQLRepository) Health(ctx context.Context) error {
	var one int
	return r.reader.GetContext(ctx, &one, `SELECT 1`)
}

func (r *MySQLRepository) ListPosts(ctx context.Context, params ListPostsParams) ([]PostSummary, Pagination, error) {
	base := `
SELECT p.id, p.title, p.content, p.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, '' AS author_bio, p.category, p.status,
p.view_count, p.like_count, p.comment_count, p.created_at, p.updated_at
FROM community_posts p
LEFT JOIN users u ON u.id = p.author_id
WHERE p.status = 'published'`
	args := make([]any, 0, 8)

	if params.Category != "" {
		base += ` AND p.category = ?`
		args = append(args, params.Category)
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "latest"
	}

	if params.Cursor != nil {
		cursorRow, err := r.cursorRow(ctx, *params.Cursor)
		if err == nil {
			switch sortBy {
			case "popular":
				score := cursorRow.LikeCount + cursorRow.CommentCount
				base += ` AND ((p.like_count + p.comment_count) < ? OR ((p.like_count + p.comment_count) = ? AND p.id < ?))`
				args = append(args, score, score, cursorRow.ID)
			case "comments":
				base += ` AND (p.comment_count < ? OR (p.comment_count = ? AND p.id < ?))`
				args = append(args, cursorRow.CommentCount, cursorRow.CommentCount, cursorRow.ID)
			default:
				base += ` AND (p.created_at < ? OR (p.created_at = ? AND p.id < ?))`
				args = append(args, cursorRow.CreatedAt, cursorRow.CreatedAt, cursorRow.ID)
			}
		}
	}

	switch sortBy {
	case "popular":
		base += ` ORDER BY (p.like_count + p.comment_count) DESC, p.id DESC`
	case "comments":
		base += ` ORDER BY p.comment_count DESC, p.id DESC`
	default:
		base += ` ORDER BY p.created_at DESC, p.id DESC`
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}
	base += ` LIMIT ?`
	args = append(args, limit+1)

	rows := []postRow{}
	if err := r.reader.SelectContext(ctx, &rows, base, args...); err != nil {
		return nil, Pagination{}, &APIError{Status: 500, Message: "获取帖子列表失败", Err: err}
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	posts := make([]PostSummary, 0, len(rows))
	for _, row := range rows {
		authorID := row.AuthorID
		posts = append(posts, PostSummary{
			ID:       row.ID,
			Title:    row.Title,
			Summary:  shorten(row.Content, 100),
			Category: row.Category,
			Author: Author{
				ID:       &authorID,
				Username: row.AuthorName,
				Avatar:   row.AuthorAvatar,
			},
			ViewCount:    row.ViewCount,
			LikeCount:    row.LikeCount,
			CommentCount: row.CommentCount,
			CreatedAt:    row.CreatedAt,
		})
	}

	var next *int64
	if hasMore && len(posts) > 0 {
		next = &posts[len(posts)-1].ID
	}

	return posts, Pagination{Cursor: next, Limit: limit, HasMore: hasMore}, nil
}

func (r *MySQLRepository) GetPost(ctx context.Context, id int64, currentUserID *int64) (*PostDetail, error) {
	row, err := r.findPost(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status != "published" && (currentUserID == nil || *currentUserID != row.AuthorID) {
		return nil, &APIError{Status: 403, Message: "没有权限查看此帖子"}
	}

	if _, err := r.writer.ExecContext(ctx, `UPDATE community_posts SET view_count = view_count + 1 WHERE id = ?`, id); err != nil {
		return nil, &APIError{Status: 500, Message: "获取帖子详情失败", Err: err}
	}
	row.ViewCount++

	liked := false
	if currentUserID != nil {
		var likedCount int
		if err := r.reader.GetContext(ctx, &likedCount, `SELECT COUNT(1) FROM user_post_likes WHERE user_id = ? AND post_id = ?`, *currentUserID, id); err == nil {
			liked = likedCount > 0
		}
	}

	comments, err := r.threadedComments(ctx, id)
	if err != nil {
		return nil, err
	}

	authorID := row.AuthorID
	bio := row.AuthorBio
	return &PostDetail{
		ID:       row.ID,
		Title:    row.Title,
		Content:  row.Content,
		Category: row.Category,
		Author: Author{
			ID:       &authorID,
			Username: row.AuthorName,
			Avatar:   avatarFallback(row.AuthorName, row.AuthorID, row.AuthorAvatar),
			Bio:      &bio,
		},
		ViewCount:          row.ViewCount,
		LikeCount:          row.LikeCount,
		CommentCount:       row.CommentCount,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
		LikedByCurrentUser: liked,
		Comments:           comments,
	}, nil
}

func (r *MySQLRepository) CreatePost(ctx context.Context, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
	if req.Title == nil || strings.TrimSpace(*req.Title) == "" || req.Content == nil || strings.TrimSpace(*req.Content) == "" || req.Category == nil || strings.TrimSpace(*req.Category) == "" {
		return nil, &APIError{Status: 400, Message: "title、content、category 为必填项"}
	}

	result, err := r.writer.ExecContext(ctx,
		`INSERT INTO community_posts (title, content, author_id, category, status, view_count, like_count, comment_count, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'published', 0, 0, 0, NOW(), NOW())`,
		strings.TrimSpace(*req.Title), strings.TrimSpace(*req.Content), userID, strings.TrimSpace(*req.Category),
	)
	if err != nil {
		return nil, &APIError{Status: 500, Message: "发布帖子失败", Err: err}
	}
	id, _ := result.LastInsertId()
	return &CreatedPost{ID: id, Title: strings.TrimSpace(*req.Title)}, nil
}

func (r *MySQLRepository) UpdatePost(ctx context.Context, id, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
	row, err := r.findPost(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := r.requireOwnerOrAdmin(ctx, userID, row.AuthorID); err != nil {
		return nil, err
	}

	title, content, category := row.Title, row.Content, row.Category
	if req.Title != nil {
		title = *req.Title
	}
	if req.Content != nil {
		content = *req.Content
	}
	if req.Category != nil {
		category = *req.Category
	}

	if _, err := r.writer.ExecContext(ctx, `UPDATE community_posts SET title = ?, content = ?, category = ?, updated_at = NOW() WHERE id = ?`, title, content, category, id); err != nil {
		return nil, &APIError{Status: 500, Message: "更新帖子失败", Err: err}
	}
	return &CreatedPost{ID: id, Title: title}, nil
}

func (r *MySQLRepository) DeletePost(ctx context.Context, id, userID int64) error {
	row, err := r.findPost(ctx, id)
	if err != nil {
		return err
	}
	if err := r.requireOwnerOrAdmin(ctx, userID, row.AuthorID); err != nil {
		return err
	}

	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return &APIError{Status: 500, Message: "删除帖子失败", Err: err}
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM comments WHERE post_id = ?`, id); err != nil {
		return &APIError{Status: 500, Message: "删除帖子失败", Err: err}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_post_likes WHERE post_id = ?`, id); err != nil {
		return &APIError{Status: 500, Message: "删除帖子失败", Err: err}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM community_posts WHERE id = ?`, id); err != nil {
		return &APIError{Status: 500, Message: "删除帖子失败", Err: err}
	}
	if err := tx.Commit(); err != nil {
		return &APIError{Status: 500, Message: "删除帖子失败", Err: err}
	}
	return nil
}

func (r *MySQLRepository) ToggleLike(ctx context.Context, id, userID int64) (*LikeResult, error) {
	if _, err := r.findPost(ctx, id); err != nil {
		return nil, err
	}

	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
	}
	defer tx.Rollback()

	var existing int
	if err := tx.GetContext(ctx, &existing, `SELECT COUNT(1) FROM user_post_likes WHERE user_id = ? AND post_id = ?`, userID, id); err != nil {
		return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
	}

	result := &LikeResult{Message: "点赞成功", Liked: true}
	if existing > 0 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_post_likes WHERE user_id = ? AND post_id = ?`, userID, id); err != nil {
			return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE community_posts SET like_count = GREATEST(like_count - 1, 0), updated_at = NOW() WHERE id = ?`, id); err != nil {
			return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
		}
		result.Message = "取消点赞成功"
		result.Liked = false
	} else {
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_post_likes (user_id, post_id) VALUES (?, ?)`, userID, id); err != nil {
			return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE community_posts SET like_count = like_count + 1, updated_at = NOW() WHERE id = ?`, id); err != nil {
			return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
		}
	}

	if err := tx.GetContext(ctx, &result.LikeCount, `SELECT like_count FROM community_posts WHERE id = ?`, id); err != nil {
		return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
	}
	if err := tx.Commit(); err != nil {
		return nil, &APIError{Status: 500, Message: "操作失败", Err: err}
	}
	return result, nil
}

func (r *MySQLRepository) ListComments(ctx context.Context, postID int64) ([]CommentItem, error) {
	if _, err := r.findPost(ctx, postID); err != nil {
		return nil, err
	}

	rows := []commentRow{}
	if err := r.reader.SelectContext(ctx, &rows, `
SELECT c.id, c.content, c.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, c.post_id, c.parent_id, c.created_at, c.updated_at
FROM comments c
LEFT JOIN users u ON u.id = c.author_id
WHERE c.post_id = ? AND c.parent_id IS NULL
ORDER BY c.created_at DESC`, postID); err != nil {
		return nil, &APIError{Status: 500, Message: "获取评论失败", Err: err}
	}

	items := make([]CommentItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toCommentItem(row))
	}
	return items, nil
}

func (r *MySQLRepository) AddComment(ctx context.Context, postID, userID int64, req AddCommentRequest) (*CreatedComment, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, &APIError{Status: 400, Message: "评论内容不能为空"}
	}
	if _, err := r.findPost(ctx, postID); err != nil {
		return nil, err
	}
	if req.ParentID != nil {
		var parentPostID int64
		err := r.reader.GetContext(ctx, &parentPostID, `SELECT post_id FROM comments WHERE id = ?`, *req.ParentID)
		if errors.Is(err, sql.ErrNoRows) || parentPostID != postID {
			return nil, &APIError{Status: 404, Message: "回复的评论不存在"}
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, &APIError{Status: 500, Message: "发布评论失败", Err: err}
		}
	}

	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return nil, &APIError{Status: 500, Message: "发布评论失败", Err: err}
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `INSERT INTO comments (content, author_id, post_id, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`, strings.TrimSpace(req.Content), userID, postID, nullableID(req.ParentID))
	if err != nil {
		return nil, &APIError{Status: 500, Message: "发布评论失败", Err: err}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE community_posts SET comment_count = comment_count + 1, updated_at = NOW() WHERE id = ?`, postID); err != nil {
		return nil, &APIError{Status: 500, Message: "发布评论失败", Err: err}
	}
	if err := tx.Commit(); err != nil {
		return nil, &APIError{Status: 500, Message: "发布评论失败", Err: err}
	}

	id, _ := result.LastInsertId()
	author, err := r.loadAuthor(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &CreatedComment{ID: id, Content: strings.TrimSpace(req.Content), Author: author}, nil
}

func (r *MySQLRepository) ListRelatedPosts(ctx context.Context, postID int64, limit int) ([]PostSummary, error) {
	current, err := r.findPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 2
	}

	rows := []postRow{}
	if err := r.reader.SelectContext(ctx, &rows, `
SELECT p.id, p.title, p.content, p.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, '' AS author_bio, p.category, p.status,
p.view_count, p.like_count, p.comment_count, p.created_at, p.updated_at
FROM community_posts p
LEFT JOIN users u ON u.id = p.author_id
WHERE p.id <> ? AND p.status = 'published' AND p.category = ?
ORDER BY (p.like_count + p.comment_count) DESC
LIMIT ?`, postID, current.Category, limit); err != nil {
		return nil, &APIError{Status: 500, Message: "获取相关推荐失败", Err: err}
	}

	posts := make([]PostSummary, 0, len(rows))
	for _, row := range rows {
		authorID := row.AuthorID
		posts = append(posts, PostSummary{
			ID:       row.ID,
			Title:    row.Title,
			Summary:  shorten(row.Content, 150),
			Category: row.Category,
			Author: Author{
				ID:       &authorID,
				Username: row.AuthorName,
				Avatar:   row.AuthorAvatar,
			},
			ViewCount:    row.ViewCount,
			LikeCount:    row.LikeCount,
			CommentCount: row.CommentCount,
			CreatedAt:    row.CreatedAt,
		})
	}
	return posts, nil
}

func (r *MySQLRepository) DeleteComment(ctx context.Context, commentID, userID int64) error {
	comment, err := r.findComment(ctx, commentID)
	if err != nil {
		return err
	}
	if err := r.requireOwnerOrAdmin(ctx, userID, comment.AuthorID); err != nil {
		return err
	}

	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM comments WHERE parent_id = ?`, commentID); err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM comments WHERE id = ?`, commentID); err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}

	var remaining int64
	if err := tx.GetContext(ctx, &remaining, `SELECT COUNT(1) FROM comments WHERE post_id = ? AND parent_id IS NULL`, comment.PostID); err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE community_posts SET comment_count = ?, updated_at = NOW() WHERE id = ?`, remaining, comment.PostID); err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}

	if err := tx.Commit(); err != nil {
		return &APIError{Status: 500, Message: "删除评论失败", Err: err}
	}
	return nil
}

func (r *MySQLRepository) findPost(ctx context.Context, id int64) (*postRow, error) {
	row := postRow{}
	err := r.reader.GetContext(ctx, &row, `
SELECT p.id, p.title, p.content, p.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, COALESCE(u.bio, '') AS author_bio, p.category, p.status,
p.view_count, p.like_count, p.comment_count, p.created_at, p.updated_at
FROM community_posts p
LEFT JOIN users u ON u.id = p.author_id
WHERE p.id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &APIError{Status: 404, Message: "帖子不存在"}
	}
	if err != nil {
		return nil, &APIError{Status: 500, Message: "获取帖子详情失败", Err: err}
	}
	return &row, nil
}

func (r *MySQLRepository) findComment(ctx context.Context, id int64) (*commentRow, error) {
	row := commentRow{}
	err := r.reader.GetContext(ctx, &row, `
SELECT c.id, c.content, c.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, c.post_id, c.parent_id, c.created_at, c.updated_at
FROM comments c
LEFT JOIN users u ON u.id = c.author_id
WHERE c.id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &APIError{Status: 404, Message: "评论不存在"}
	}
	if err != nil {
		return nil, &APIError{Status: 500, Message: "获取评论失败", Err: err}
	}
	return &row, nil
}

func (r *MySQLRepository) cursorRow(ctx context.Context, id int64) (*postRow, error) {
	row := postRow{}
	if err := r.reader.GetContext(ctx, &row, `SELECT id, created_at, like_count, comment_count FROM community_posts WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *MySQLRepository) requireOwnerOrAdmin(ctx context.Context, userID, ownerID int64) error {
	if userID == ownerID {
		return nil
	}

	var role string
	if err := r.reader.GetContext(ctx, &role, `SELECT role FROM users WHERE id = ?`, userID); err != nil {
		return &APIError{Status: 500, Message: "权限校验失败", Err: err}
	}
	if role != "admin" {
		return &APIError{Status: 403, Message: "没有权限执行该操作"}
	}
	return nil
}

func (r *MySQLRepository) threadedComments(ctx context.Context, postID int64) ([]CommentItem, error) {
	topRows := []commentRow{}
	if err := r.reader.SelectContext(ctx, &topRows, `
SELECT c.id, c.content, c.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, c.post_id, c.parent_id, c.created_at, c.updated_at
FROM comments c
LEFT JOIN users u ON u.id = c.author_id
WHERE c.post_id = ? AND c.parent_id IS NULL
ORDER BY c.created_at DESC`, postID); err != nil {
		return nil, &APIError{Status: 500, Message: "获取帖子详情失败", Err: err}
	}

	items := make([]CommentItem, 0, len(topRows))
	for _, row := range topRows {
		repliesRows := []commentRow{}
		if err := r.reader.SelectContext(ctx, &repliesRows, `
SELECT c.id, c.content, c.author_id, COALESCE(u.username, '匿名用户') AS author_username,
COALESCE(u.avatar, '') AS author_avatar, c.post_id, c.parent_id, c.created_at, c.updated_at
FROM comments c
LEFT JOIN users u ON u.id = c.author_id
WHERE c.parent_id = ?`, row.ID); err != nil {
			return nil, &APIError{Status: 500, Message: "获取帖子详情失败", Err: err}
		}

		item := toCommentItem(row)
		for _, reply := range repliesRows {
			item.Replies = append(item.Replies, toCommentItem(reply))
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) loadAuthor(ctx context.Context, userID int64) (Author, error) {
	type authorRow struct {
		ID       int64  `db:"id"`
		Username string `db:"username"`
		Avatar   string `db:"avatar"`
	}

	row := authorRow{}
	if err := r.reader.GetContext(ctx, &row, `SELECT id, username, COALESCE(avatar, '') AS avatar FROM users WHERE id = ?`, userID); err != nil {
		return Author{}, &APIError{Status: 500, Message: "发布评论失败", Err: err}
	}
	id := row.ID
	return Author{
		ID:       &id,
		Username: row.Username,
		Avatar:   avatarFallback(row.Username, row.ID, row.Avatar),
	}, nil
}

func toCommentItem(row commentRow) CommentItem {
	authorID := row.AuthorID
	return CommentItem{
		ID:      row.ID,
		Content: row.Content,
		Author: Author{
			ID:       &authorID,
			Username: row.AuthorName,
			Avatar:   avatarFallback(row.AuthorName, row.AuthorID, row.AuthorAvatar),
		},
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func nullableID(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func avatarFallback(username string, id int64, avatar string) string {
	if strings.TrimSpace(avatar) != "" {
		return avatar
	}
	initial := "U"
	if username != "" {
		initial = strings.ToUpper(string([]rune(username)[0]))
	}
	return fmt.Sprintf("https://picsum.photos/seed/%s%d/100", initial, id)
}

func shorten(content string, max int) string {
	runes := []rune(content)
	if len(runes) <= max {
		return content
	}
	return string(runes[:max]) + "..."
}
