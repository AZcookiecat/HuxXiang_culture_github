package community

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"huxiang/backend/go_post_service/internal/app"

	"github.com/sony/gobreaker"
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

type LikeResult struct {
	Message   string `json:"message"`
	LikeCount int64  `json:"like_count"`
	Liked     bool   `json:"liked"`
}

type Repository interface {
	Health(ctx context.Context) error
	ListPosts(ctx context.Context, params ListPostsParams) ([]PostSummary, Pagination, error)
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
	cacheKey := fmt.Sprintf("community:posts:%s:%s:%d:%v", params.SortBy, params.Category, params.Limit, params.Cursor)
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

func (s *Service) GetPost(ctx context.Context, id int64, currentUserID *int64) (*PostDetail, error) {
	value, err := s.execute(func() (any, error) { return s.repo.GetPost(ctx, id, currentUserID) })
	if err != nil {
		return nil, err
	}
	return value.(*PostDetail), nil
}

func (s *Service) CreatePost(ctx context.Context, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
	value, err := s.execute(func() (any, error) { return s.repo.CreatePost(ctx, userID, req) })
	if err != nil {
		return nil, err
	}
	s.invalidateCommunityCache()
	return value.(*CreatedPost), nil
}

func (s *Service) UpdatePost(ctx context.Context, id, userID int64, req UpsertPostRequest) (*CreatedPost, error) {
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
	value, err := s.execute(func() (any, error) { return s.repo.AddComment(ctx, postID, userID, req) })
	if err != nil {
		return nil, err
	}
	s.invalidateCommunityCache()
	return value.(*CreatedComment), nil
}

func (s *Service) ListRelatedPosts(ctx context.Context, postID int64, limit int) ([]PostSummary, error) {
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
		if apiErr, ok := err.(*APIError); ok {
			return nil, apiErr
		}
		return nil, &APIError{Status: http.StatusServiceUnavailable, Message: "服务暂时不可用", Err: err}
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
