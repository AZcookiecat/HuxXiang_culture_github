package community

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestListCommentsReturnsThreadedReplies(t *testing.T) {
	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	now := time.Now()
	postRows := sqlmock.NewRows([]string{
		"id", "title", "content", "author_id", "author_username", "author_avatar", "author_bio",
		"category", "status", "view_count", "like_count", "comment_count", "created_at", "updated_at",
	}).AddRow(1, "title", "content", 7, "author", "", "", "history", "published", 0, 0, 2, now, now)

	mock.ExpectQuery(`(?s)SELECT p\.id, p\.title, p\.content, p\.author_id.*WHERE p\.id = \?`).
		WithArgs(int64(1)).
		WillReturnRows(postRows)

	topRows := sqlmock.NewRows([]string{
		"id", "content", "author_id", "author_username", "author_avatar", "post_id", "parent_id", "created_at", "updated_at",
	}).AddRow(10, "top-level", 7, "author", "", 1, nil, now, now)

	mock.ExpectQuery(`(?s)SELECT c\.id, c\.content, c\.author_id.*WHERE c\.post_id = \? AND c\.parent_id IS NULL.*`).
		WithArgs(int64(1)).
		WillReturnRows(topRows)

	replyRows := sqlmock.NewRows([]string{
		"id", "content", "author_id", "author_username", "author_avatar", "post_id", "parent_id", "created_at", "updated_at",
	}).AddRow(11, "reply", 8, "replier", "", 1, 10, now, now)

	mock.ExpectQuery(`(?s)SELECT c\.id, c\.content, c\.author_id.*WHERE c\.parent_id = \?`).
		WithArgs(int64(10)).
		WillReturnRows(replyRows)

	comments, err := repo.ListComments(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListComments returned error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 top-level comment, got %d", len(comments))
	}
	if len(comments[0].Replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(comments[0].Replies))
	}
	if comments[0].Replies[0].ID != 11 {
		t.Fatalf("expected reply id 11, got %d", comments[0].Replies[0].ID)
	}

	assertExpectations(t, mock)
}

func TestListPostsAppliesKeywordFilter(t *testing.T) {
	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "content", "author_id", "author_username", "author_avatar", "author_bio",
		"category", "status", "view_count", "like_count", "comment_count", "created_at", "updated_at",
	}).AddRow(1, "go post", "clean code", 7, "author", "", "", "tech", "published", 1, 2, 3, now, now)

	mock.ExpectQuery(`(?s)SELECT.*FROM community_posts p.*WHERE p.status = 'published'.*p.title LIKE \? OR p.content LIKE \?.*ORDER BY p.created_at DESC, p.id DESC.*LIMIT \?`).
		WithArgs("%go%", "%go%", 11).
		WillReturnRows(rows)

	posts, pagination, err := repo.ListPosts(context.Background(), ListPostsParams{
		Limit:   10,
		SortBy:  "latest",
		Keyword: "go",
	})
	if err != nil {
		t.Fatalf("ListPosts returned error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if pagination.HasMore {
		t.Fatalf("expected no more pages")
	}

	assertExpectations(t, mock)
}

func TestListUserPostsFiltersByStatus(t *testing.T) {
	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "content", "author_id", "author_username", "author_avatar", "author_bio",
		"category", "status", "view_count", "like_count", "comment_count", "created_at", "updated_at",
	}).AddRow(2, "draft post", "content", 7, "author", "", "", "tech", "draft", 0, 0, 0, now, now)

	mock.ExpectQuery(`(?s)SELECT.*FROM community_posts p.*WHERE p.author_id = \?.*AND p.status = \?.*ORDER BY p.created_at DESC, p.id DESC.*LIMIT \?`).
		WithArgs(int64(7), "draft", 11).
		WillReturnRows(rows)

	posts, _, err := repo.ListUserPosts(context.Background(), 7, UserPostsParams{
		Limit:  10,
		Status: "draft",
	})
	if err != nil {
		t.Fatalf("ListUserPosts returned error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].Status != "draft" {
		t.Fatalf("expected draft status, got %s", posts[0].Status)
	}

	assertExpectations(t, mock)
}

func TestListCategoriesReturnsCounts(t *testing.T) {
	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"name", "post_count"}).
		AddRow("tech", 5).
		AddRow("history", 3)

	mock.ExpectQuery(`(?s)SELECT.*COUNT\(1\) AS post_count.*FROM community_posts p.*GROUP BY p.category.*ORDER BY post_count DESC, name ASC`).
		WillReturnRows(rows)

	stats, err := repo.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("ListCategories returned error: %v", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(stats))
	}
	if stats[0].Name != "tech" || stats[0].PostCount != 5 {
		t.Fatalf("unexpected first category: %+v", stats[0])
	}

	assertExpectations(t, mock)
}

func TestDeleteCommentRecomputesTotalCommentCount(t *testing.T) {
	repo, mock, cleanup := newMockRepository(t)
	defer cleanup()

	now := time.Now()
	commentRows := sqlmock.NewRows([]string{
		"id", "content", "author_id", "author_username", "author_avatar", "post_id", "parent_id", "created_at", "updated_at",
	}).AddRow(10, "reply", 7, "author", "", 1, nil, now, now)

	mock.ExpectQuery(`(?s)SELECT c\.id, c\.content, c\.author_id.*WHERE c\.id = \?`).
		WithArgs(int64(10)).
		WillReturnRows(commentRows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM comments WHERE parent_id = ?`)).
		WithArgs(int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM comments WHERE id = ?`)).
		WithArgs(int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(1) FROM comments WHERE post_id = ?`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE community_posts SET comment_count = ?, updated_at = NOW() WHERE id = ?`)).
		WithArgs(int64(5), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.DeleteComment(context.Background(), 10, 7); err != nil {
		t.Fatalf("DeleteComment returned error: %v", err)
	}

	assertExpectations(t, mock)
}

func newMockRepository(t *testing.T) (*MySQLRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewMySQLRepository(sqlxDB, sqlxDB)
	cleanup := func() {
		_ = sqlxDB.Close()
	}
	return repo, mock, cleanup
}

func assertExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
