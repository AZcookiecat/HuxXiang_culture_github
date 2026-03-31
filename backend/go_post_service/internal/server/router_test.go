package server

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"huxiang/backend/go_post_service/internal/app"
	"huxiang/backend/go_post_service/internal/community"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewRouterRegistersCommunityRoutes(t *testing.T) {
	handler := community.NewHandler(nil, "test-secret")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected router registration without panic, got %v", r)
		}
	}()

	router := NewRouter(testConfig(), testLogger(), app.NewMetrics(), handler, func() error { return nil })
	if router == nil {
		t.Fatal("expected router")
	}
}

func TestCommunityListPostsEndpoint(t *testing.T) {
	repo := &testCommunityRepository{}
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/community/posts?limit=5&sortBy=popular", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", rec.Code, rec.Body.String())
	}
	if repo.lastListPosts.Limit != 5 {
		t.Fatalf("expected limit 5, got %d", repo.lastListPosts.Limit)
	}
	if repo.lastListPosts.SortBy != "popular" {
		t.Fatalf("expected sortBy popular, got %q", repo.lastListPosts.SortBy)
	}
}

func TestCommunityListCommentsEndpointUsesPostIDPath(t *testing.T) {
	repo := &testCommunityRepository{}
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/community/posts/42/comments", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", rec.Code, rec.Body.String())
	}
	if repo.lastCommentsPostID != 42 {
		t.Fatalf("expected comments to load for post 42, got %d", repo.lastCommentsPostID)
	}
}

func TestCommunityMyPostsRouteIsRegistered(t *testing.T) {
	router := newTestRouter(&testCommunityRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/community/posts/mine", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated /posts/mine, got %d with body %s", rec.Code, rec.Body.String())
	}
}

func TestCommunityAddCommentEndpointAcceptsJWT(t *testing.T) {
	repo := &testCommunityRepository{}
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/community/posts/42/comments", bytes.NewBufferString(`{"content":"  hello  "}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testJWT(t, 7))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d with body %s", rec.Code, rec.Body.String())
	}
	if repo.lastCommentPostID != 42 {
		t.Fatalf("expected comment to target post 42, got %d", repo.lastCommentPostID)
	}
	if repo.lastCommentUserID != 7 {
		t.Fatalf("expected comment user 7, got %d", repo.lastCommentUserID)
	}
	if repo.lastComment.Content != "hello" {
		t.Fatalf("expected trimmed comment content, got %q", repo.lastComment.Content)
	}
}

func TestCommunityGetPostRejectsBadID(t *testing.T) {
	router := newTestRouter(&testCommunityRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/community/posts/not-a-number", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d with body %s", rec.Code, rec.Body.String())
	}
}

type testCommunityRepository struct {
	lastListPosts      community.ListPostsParams
	lastCommentsPostID int64
	lastCommentPostID  int64
	lastCommentUserID  int64
	lastComment        community.AddCommentRequest
}

func (r *testCommunityRepository) Health(context.Context) error { return nil }

func (r *testCommunityRepository) ListPosts(_ context.Context, params community.ListPostsParams) ([]community.PostSummary, community.Pagination, error) {
	r.lastListPosts = params
	return []community.PostSummary{{ID: 1, Title: "post"}}, community.Pagination{Limit: params.Limit}, nil
}

func (r *testCommunityRepository) ListUserPosts(context.Context, int64, community.UserPostsParams) ([]community.UserPostSummary, community.Pagination, error) {
	return []community.UserPostSummary{}, community.Pagination{Limit: 10}, nil
}

func (r *testCommunityRepository) ListCategories(context.Context) ([]community.CategoryStat, error) {
	return []community.CategoryStat{{Name: "history", PostCount: 1}}, nil
}

func (r *testCommunityRepository) GetPost(_ context.Context, id int64, _ *int64) (*community.PostDetail, error) {
	return &community.PostDetail{ID: id, Title: "post"}, nil
}

func (r *testCommunityRepository) CreatePost(context.Context, int64, community.UpsertPostRequest) (*community.CreatedPost, error) {
	return &community.CreatedPost{ID: 1, Title: "created"}, nil
}

func (r *testCommunityRepository) UpdatePost(context.Context, int64, int64, community.UpsertPostRequest) (*community.CreatedPost, error) {
	return &community.CreatedPost{ID: 1, Title: "updated"}, nil
}

func (r *testCommunityRepository) DeletePost(context.Context, int64, int64) error { return nil }

func (r *testCommunityRepository) ToggleLike(context.Context, int64, int64) (*community.LikeResult, error) {
	return &community.LikeResult{Message: "ok", LikeCount: 1, Liked: true}, nil
}

func (r *testCommunityRepository) ListComments(_ context.Context, postID int64) ([]community.CommentItem, error) {
	r.lastCommentsPostID = postID
	return []community.CommentItem{{ID: 1, Content: "comment"}}, nil
}

func (r *testCommunityRepository) AddComment(_ context.Context, postID, userID int64, req community.AddCommentRequest) (*community.CreatedComment, error) {
	r.lastCommentPostID = postID
	r.lastCommentUserID = userID
	r.lastComment = req
	return &community.CreatedComment{ID: 1, Content: req.Content}, nil
}

func (r *testCommunityRepository) ListRelatedPosts(context.Context, int64, int) ([]community.PostSummary, error) {
	return []community.PostSummary{{ID: 2, Title: "related"}}, nil
}

func (r *testCommunityRepository) DeleteComment(context.Context, int64, int64) error { return nil }

func newTestRouter(repo community.Repository) http.Handler {
	service := community.NewService(repo, app.NewInMemoryCache(), nil, 0)
	handler := community.NewHandler(service, "test-secret")
	return NewRouter(testConfig(), testLogger(), app.NewMetrics(), handler, func() error { return nil })
}

func testConfig() app.Config {
	return app.Config{
		ReadTimeout:      3 * time.Second,
		RateLimitRPS:     100,
		RateLimitBurst:   100,
		CORSAllowOrigins: []string{"*"},
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testJWT(t *testing.T, userID int64) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
	})
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return signed
}
