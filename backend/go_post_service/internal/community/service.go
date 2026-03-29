package community

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"huxiang/backend/go_post_service/internal/app"

	"github.com/sony/gobreaker"
)

const (
	defaultPageLimit  = 10
	maxPageLimit      = 50
	defaultRelatedNum = 2
	maxRelatedNum     = 10
	maxTitleLength    = 120
	maxCategoryLength = 32
	maxPostLength     = 10000
	maxCommentLength  = 1000
)

type APIError struct {
	Status  int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func (e *APIError) StatusCode() int {
	if e.Status <= 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e *APIError) PublicMessage() string {
	if e.Message == "" {
		return "服务器开小差了，请稍后重试"
	}
	return e.Message
}

type Author struct {
	ID       *int64  `json:"id"`
	Username string  `json:"username"`
	Avatar   string  `json:"avatar"`
	Bio      *string `json:"bio,omitempty"`
}

type PostSummary struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Author       Author    `json:"author"`
	Category     string    `json:"category"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type CommentItem struct {
	ID        int64         `json:"id"`
	Content   string        `json:"content"`
	Author    Author        `json:"author"`
	Replies   []CommentItem `json:"replies,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type PostDetail struct {
	ID                 int64         `json:"id"`
	Title              string        `json:"title"`
	Content            string        `json:"content"`
	Author             Author        `json:"author"`
	Category           string        `json:"category"`
	ViewCount          int64         `json:"view_count"`
	LikeCount          int64         `json:"like_count"`
	CommentCount       int64         `json:"comment_count"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	LikedByCurrentUser bool          `json:"liked_by_current_user"`
	Comments           []CommentItem `json:"comments"`
}

type Pagination struct {
	Cursor  *int64 `json:"cursor"`
	Limit   int    `json:"limit"`
	HasMore bool   `json:"has_more"`
}

type ListPostsParams struct {
	Cursor   *int64
	Limit    int
	Category string
	SortBy   string
	Keyword  string
}

type UserPostsParams struct {
	Cursor *int64
	Limit  int
	Status string
}

type UpsertPostRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
}

type AddCommentRequest struct {
	Content  string `json:"content"`
	ParentID *int64 `json:"parent_id"`
}

type CreatedPost struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type CreatedComment struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Author  Author `json:"author"`
}

type UserPostSummary struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Category     string    `json:"category"`
	Status       string    `json:"status"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CategoryStat struct {
	Name      string `json:"name"`
	PostCount int64  `json:"post_count"`
}

type LikeResult struct {
	Message   string `json:"message"`
	LikeCount int64  `json:"like_count"`
	Liked     bool   `json:"liked"`
}

type Repository interface {
	Health(ctx context.Context) error
	ListPosts(ctx context.Context, params ListPostsParams) ([]PostSummary, Pagination, error)
	ListUserPosts(ctx context.Context, userID int64, params UserPostsParams) ([]UserPostSummary, Pagination, error)
	ListCategories(ctx context.Context) ([]CategoryStat, error)
	GetPost(ctx context.Context, id int64, currentUserID *int64) (*PostDetail, error)
	CreatePost(ctx context.Context, userID int64, req UpsertPostRequest) (*CreatedPost, error)
	UpdatePost(ctx context.Context, id, userID int64, req UpsertPostRequest) (*CreatedPost, error)
	DeletePost(ctx context.Context, id, userID int64) error
	ToggleLike(ctx context.Context, id, userID int64) (*LikeResult, error)
	ListComments(ctx context.Context, postID int64) ([]CommentItem, error)
	AddComment(ctx context.Context, postID, userID int64, req AddCommentRequest) (*CreatedComment, error)
	ListRelatedPosts(ctx context.Context, postID int64, limit int) ([]PostSummary, error)
	DeleteComment(ctx context.Context, commentID, userID int64) error
}

type Service struct {
	repo    Repository
	cache   app.Cache
	events  *app.EventBus
	breaker *gobreaker.CircuitBreaker
	ttl     time.Duration
}

func NewService(repo Repository, cache app.Cache, events *app.EventBus, ttl time.Duration) *Service {
	return &Service{
		repo:   repo,
		cache:  cache,
		events: events,
		breaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "community-repository",
			MaxRequests: 10,
			Interval:    30 * time.Second,
			Timeout:     10 * time.Second,
		}),
		ttl: ttl,
	}
}

func (s *Service) Health(ctx context.Context) error {
	_, err := s.execute(func() (any, error) { return nil, s.repo.Health(ctx) })
	return err
}

func (s *Service) ListPosts(ctx context.Context, params ListPostsParams) ([]PostSummary, Pagination, error) {
	params = normalizeListPostsParams(params)
	cacheKey := fmt.Sprintf("community:posts:%s:%s:%s:%d:%v", params.SortBy, params.Category, params.Keyword, params.Limit, params.Cursor)
	if cached, ok := s.cache.Get(cacheKey); ok {
		result := cached.(struct {
			Posts      []PostSummary
			Pagination Pagination
		})
		return result.Posts, result.Pagination, nil
	}

	value, err := s.execute(func() (any, error) {
		posts, pagination, err := s.repo.ListPosts(ctx, params)
		if err != nil {
			return nil, err
		}
		return struct {
			Posts      []PostSummary
			Pagination Pagination
		}{Posts: posts, Pagination: pagination}, nil
	})
	if err != nil {
		return nil, Pagination{}, err
	}

	result := value.(struct {
		Posts      []PostSummary
		Pagination Pagination
	})
	s.cache.Set(cacheKey, result, s.ttl)
	return result.Posts, result.Pagination, nil
}

func (s *Service) ListUserPosts(ctx context.Context, userID int64, params UserPostsParams) ([]UserPostSummary, Pagination, error) {
	params = normalizeUserPostsParams(params)
	value, err := s.execute(func() (any, error) {
		posts, pagination, err := s.repo.ListUserPosts(ctx, userID, params)
		if err != nil {
			return nil, err
		}
		return struct {
			Posts      []UserPostSummary
			Pagination Pagination
		}{Posts: posts, Pagination: pagination}, nil
	})
	if err != nil {
		return nil, Pagination{}, err
	}

	result := value.(struct {
		Posts      []UserPostSummary
		Pagination Pagination
	})
	return result.Posts, result.Pagination, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]CategoryStat, error) {
	const cacheKey = "community:categories"
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]CategoryStat), nil
	}

	value, err := s.execute(func() (any, error) { return s.repo.ListCategories(ctx) })
	if err != nil {
		return nil, err
	}

	stats := value.([]CategoryStat)
	s.cache.Set(cacheKey, stats, s.ttl)
	return stats, nil
}

func (s *Service) GetPost(ctx context.Context, id int64, currentUserID *int64) (*PostDetail, error) {
	value, err := s.execute(func() (any, error) { return s.repo.GetPost(ctx, id, currentUserID) })
	if err != nil {
		return nil, err
	}
	return value.(*PostDetail), nil
}

func (s *Service) CreatePost(ctx context.Context, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
	req, err := validateCreatePostRequest(req)
	if err != nil {
		return nil, err
	}

	value, err := s.execute(func() (any, error) { return s.repo.CreatePost(ctx, userID, req) })
	if err != nil {
		return nil, err
	}

	s.invalidateCommunityCache()
	return value.(*CreatedPost), nil
}

func (s *Service) UpdatePost(ctx context.Context, id, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
	req, err := validateUpdatePostRequest(req)
	if err != nil {
		return nil, err
	}

	value, err := s.execute(func() (any, error) { return s.repo.UpdatePost(ctx, id, userID, req) })
	if err != nil {
		return nil, err
	}

	s.invalidateCommunityCache()
	return value.(*CreatedPost), nil
}

func (s *Service) DeletePost(ctx context.Context, id, userID int64) error {
	_, err := s.execute(func() (any, error) { return nil, s.repo.DeletePost(ctx, id, userID) })
	if err == nil {
		s.invalidateCommunityCache()
	}
	return err
}

func (s *Service) ToggleLike(ctx context.Context, id, userID int64) (*LikeResult, error) {
	value, err := s.execute(func() (any, error) { return s.repo.ToggleLike(ctx, id, userID) })
	if err != nil {
		return nil, err
	}

	s.invalidateCommunityCache()
	return value.(*LikeResult), nil
}

func (s *Service) ListComments(ctx context.Context, postID int64) ([]CommentItem, error) {
	cacheKey := fmt.Sprintf("community:comments:%d", postID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]CommentItem), nil
	}

	value, err := s.execute(func() (any, error) { return s.repo.ListComments(ctx, postID) })
	if err != nil {
		return nil, err
	}

	comments := value.([]CommentItem)
	s.cache.Set(cacheKey, comments, s.ttl)
	return comments, nil
}

func (s *Service) AddComment(ctx context.Context, postID, userID int64, req AddCommentRequest) (*CreatedComment, error) {
	req, err := validateCommentRequest(req)
	if err != nil {
		return nil, err
	}

	value, err := s.execute(func() (any, error) { return s.repo.AddComment(ctx, postID, userID, req) })
	if err != nil {
		return nil, err
	}

	s.invalidateCommunityCache()
	return value.(*CreatedComment), nil
}

func (s *Service) ListRelatedPosts(ctx context.Context, postID int64, limit int) ([]PostSummary, error) {
	limit = clamp(limit, defaultRelatedNum, maxRelatedNum)
	cacheKey := fmt.Sprintf("community:related:%d:%d", postID, limit)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]PostSummary), nil
	}

	value, err := s.execute(func() (any, error) { return s.repo.ListRelatedPosts(ctx, postID, limit) })
	if err != nil {
		return nil, err
	}

	posts := value.([]PostSummary)
	s.cache.Set(cacheKey, posts, s.ttl)
	return posts, nil
}

func (s *Service) DeleteComment(ctx context.Context, commentID, userID int64) error {
	_, err := s.execute(func() (any, error) { return nil, s.repo.DeleteComment(ctx, commentID, userID) })
	if err == nil {
		s.invalidateCommunityCache()
	}
	return err
}

func (s *Service) execute(fn func() (any, error)) (any, error) {
	value, err := s.breaker.Execute(fn)
	if err != nil {
		var apiErr *APIError
		if ok := errorAs(err, &apiErr); ok {
			return nil, apiErr
		}
		return nil, &APIError{
			Status:  http.StatusServiceUnavailable,
			Message: "服务暂时不可用，请稍后重试",
			Err:     err,
		}
	}
	return value, nil
}

func (s *Service) invalidateCommunityCache() {
	if s.events != nil {
		s.events.Publish(app.Event{Topic: "community.cache.invalidate"})
		return
	}
	s.cache.DeletePrefix("community:")
}

func normalizeListPostsParams(params ListPostsParams) ListPostsParams {
	params.Limit = clamp(params.Limit, defaultPageLimit, maxPageLimit)
	params.SortBy = normalizeSortBy(params.SortBy)
	params.Category = strings.TrimSpace(params.Category)
	params.Keyword = strings.TrimSpace(params.Keyword)
	return params
}

func normalizeUserPostsParams(params UserPostsParams) UserPostsParams {
	params.Limit = clamp(params.Limit, defaultPageLimit, maxPageLimit)
	params.Status = normalizePostStatus(params.Status)
	return params
}

func validateCreatePostRequest(req UpsertPostRequest) (UpsertPostRequest, error) {
	title, err := normalizeRequired(req.Title, "标题不能为空", maxTitleLength)
	if err != nil {
		return UpsertPostRequest{}, err
	}
	content, err := normalizeRequired(req.Content, "帖子内容不能为空", maxPostLength)
	if err != nil {
		return UpsertPostRequest{}, err
	}
	category, err := normalizeRequired(req.Category, "分类不能为空", maxCategoryLength)
	if err != nil {
		return UpsertPostRequest{}, err
	}

	return UpsertPostRequest{Title: &title, Content: &content, Category: &category}, nil
}

func validateUpdatePostRequest(req UpsertPostRequest) (UpsertPostRequest, error) {
	var (
		title    *string
		content  *string
		category *string
		count    int
	)

	if req.Title != nil {
		value, err := normalizeRequired(req.Title, "标题不能为空", maxTitleLength)
		if err != nil {
			return UpsertPostRequest{}, err
		}
		title = &value
		count++
	}
	if req.Content != nil {
		value, err := normalizeRequired(req.Content, "帖子内容不能为空", maxPostLength)
		if err != nil {
			return UpsertPostRequest{}, err
		}
		content = &value
		count++
	}
	if req.Category != nil {
		value, err := normalizeRequired(req.Category, "分类不能为空", maxCategoryLength)
		if err != nil {
			return UpsertPostRequest{}, err
		}
		category = &value
		count++
	}

	if count == 0 {
		return UpsertPostRequest{}, &APIError{Status: http.StatusBadRequest, Message: "至少提供一个可更新字段"}
	}

	return UpsertPostRequest{Title: title, Content: content, Category: category}, nil
}

func validateCommentRequest(req AddCommentRequest) (AddCommentRequest, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return AddCommentRequest{}, &APIError{Status: http.StatusBadRequest, Message: "评论内容不能为空"}
	}
	if len([]rune(content)) > maxCommentLength {
		return AddCommentRequest{}, &APIError{Status: http.StatusBadRequest, Message: "评论内容不能超过 1000 个字符"}
	}
	req.Content = content
	return req, nil
}

func normalizeRequired(value *string, emptyMessage string, maxLength int) (string, error) {
	if value == nil {
		return "", &APIError{Status: http.StatusBadRequest, Message: emptyMessage}
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return "", &APIError{Status: http.StatusBadRequest, Message: emptyMessage}
	}
	if len([]rune(trimmed)) > maxLength {
		return "", &APIError{Status: http.StatusBadRequest, Message: fmt.Sprintf("字段长度不能超过 %d 个字符", maxLength)}
	}
	return trimmed, nil
}

func normalizeSortBy(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "popular":
		return "popular"
	case "comments":
		return "comments"
	default:
		return "latest"
	}
}

func normalizePostStatus(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "published":
		return "published"
	case "draft":
		return "draft"
	case "archived":
		return "archived"
	default:
		return "all"
	}
}

func clamp(value, fallback, max int) int {
	if value <= 0 {
		return fallback
	}
	if value > max {
		return max
	}
	return value
}

func errorAs(err error, target **APIError) bool {
	if err == nil {
		return false
	}
	if value, ok := err.(*APIError); ok {
		*target = value
		return true
	}
	type unwrapper interface {
		Unwrap() error
	}
	if wrapped, ok := err.(unwrapper); ok {
		return errorAs(wrapped.Unwrap(), target)
	}
	return false
}
