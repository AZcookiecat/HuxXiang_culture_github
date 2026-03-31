package community

import (
	"context"
	"testing"

	"huxiang/backend/go_post_service/internal/app"
)

type stubRepository struct{}
type captureRepository struct {
	lastComment AddCommentRequest
}

func (stubRepository) Health(context.Context) error { return nil }
func (stubRepository) ListPosts(context.Context, ListPostsParams) ([]PostSummary, Pagination, error) {
	return nil, Pagination{}, nil
}
func (stubRepository) ListUserPosts(context.Context, int64, UserPostsParams) ([]UserPostSummary, Pagination, error) {
	return nil, Pagination{}, nil
}
func (stubRepository) ListCategories(context.Context) ([]CategoryStat, error)      { return nil, nil }
func (stubRepository) GetPost(context.Context, int64, *int64) (*PostDetail, error) { return nil, nil }
func (stubRepository) CreatePost(context.Context, int64, UpsertPostRequest) (*CreatedPost, error) {
	return &CreatedPost{ID: 1, Title: "ok"}, nil
}
func (stubRepository) UpdatePost(context.Context, int64, int64, UpsertPostRequest) (*CreatedPost, error) {
	return &CreatedPost{ID: 1, Title: "ok"}, nil
}
func (stubRepository) DeletePost(context.Context, int64, int64) error                { return nil }
func (stubRepository) ToggleLike(context.Context, int64, int64) (*LikeResult, error) { return nil, nil }
func (stubRepository) ListComments(context.Context, int64) ([]CommentItem, error)    { return nil, nil }
func (stubRepository) AddComment(context.Context, int64, int64, AddCommentRequest) (*CreatedComment, error) {
	return &CreatedComment{ID: 1, Content: "ok"}, nil
}
func (stubRepository) ListRelatedPosts(context.Context, int64, int) ([]PostSummary, error) {
	return nil, nil
}
func (stubRepository) DeleteComment(context.Context, int64, int64) error { return nil }

func (r *captureRepository) Health(context.Context) error { return nil }
func (r *captureRepository) ListPosts(context.Context, ListPostsParams) ([]PostSummary, Pagination, error) {
	return nil, Pagination{}, nil
}
func (r *captureRepository) ListUserPosts(context.Context, int64, UserPostsParams) ([]UserPostSummary, Pagination, error) {
	return nil, Pagination{}, nil
}
func (r *captureRepository) ListCategories(context.Context) ([]CategoryStat, error) { return nil, nil }
func (r *captureRepository) GetPost(context.Context, int64, *int64) (*PostDetail, error) {
	return nil, nil
}
func (r *captureRepository) CreatePost(context.Context, int64, UpsertPostRequest) (*CreatedPost, error) {
	return &CreatedPost{ID: 1, Title: "ok"}, nil
}
func (r *captureRepository) UpdatePost(context.Context, int64, int64, UpsertPostRequest) (*CreatedPost, error) {
	return &CreatedPost{ID: 1, Title: "ok"}, nil
}
func (r *captureRepository) DeletePost(context.Context, int64, int64) error { return nil }
func (r *captureRepository) ToggleLike(context.Context, int64, int64) (*LikeResult, error) {
	return nil, nil
}
func (r *captureRepository) ListComments(context.Context, int64) ([]CommentItem, error) {
	return nil, nil
}
func (r *captureRepository) AddComment(_ context.Context, _ int64, _ int64, req AddCommentRequest) (*CreatedComment, error) {
	r.lastComment = req
	return &CreatedComment{ID: 1, Content: req.Content}, nil
}
func (r *captureRepository) ListRelatedPosts(context.Context, int64, int) ([]PostSummary, error) {
	return nil, nil
}
func (r *captureRepository) DeleteComment(context.Context, int64, int64) error { return nil }

func TestServiceCreatePostValidatesRequiredFields(t *testing.T) {
	service := NewService(stubRepository{}, app.NewInMemoryCache(), nil, 0)

	_, err := service.CreatePost(context.Background(), 1, UpsertPostRequest{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Status != 400 {
		t.Fatalf("expected 400, got %d", apiErr.Status)
	}
}

func TestServiceUpdatePostRequiresAtLeastOneField(t *testing.T) {
	service := NewService(stubRepository{}, app.NewInMemoryCache(), nil, 0)

	_, err := service.UpdatePost(context.Background(), 1, 1, UpsertPostRequest{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "至少提供一个可更新字段" {
		t.Fatalf("unexpected message: %s", apiErr.Message)
	}
}

func TestServiceAddCommentTrimsContent(t *testing.T) {
	repo := &captureRepository{}
	service := NewService(repo, app.NewInMemoryCache(), nil, 0)

	comment, err := service.AddComment(context.Background(), 1, 1, AddCommentRequest{Content: "  hi  "})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comment.Content != "hi" {
		t.Fatalf("unexpected comment response: %+v", comment)
	}
	if repo.lastComment.Content != "hi" {
		t.Fatalf("expected trimmed content, got %q", repo.lastComment.Content)
	}
}
