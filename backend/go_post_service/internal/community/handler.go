package community

import (
	"net/http"
	"strconv"
	"strings"

	"huxiang/backend/go_post_service/internal/app"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *Service
	jwtSecret string
}

func NewHandler(service *Service, jwtSecret string) *Handler {
	return &Handler{service: service, jwtSecret: jwtSecret}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	group.GET("/posts", h.ListPosts)
	group.GET("/posts/mine", h.ListMyPosts)
	group.GET("/posts/:id", h.GetPost)
	group.POST("/posts", h.CreatePost)
	group.PUT("/posts/:id", h.UpdatePost)
	group.DELETE("/posts/:id", h.DeletePost)
	group.POST("/posts/:id/like", h.ToggleLike)
	group.GET("/posts/:postID/comments", h.ListComments)
	group.POST("/posts/:postID/comments", h.AddComment)
	group.GET("/posts/related/:postID", h.ListRelatedPosts)
	group.GET("/categories", h.ListCategories)
	group.DELETE("/comments/:commentID", h.DeleteComment)
}

func (h *Handler) ListPosts(c *gin.Context) {
	cursor, err := parseOptionalInt64(c.Query("cursor"))
	if err != nil {
		app.AbortError(c, badRequest("cursor 参数格式错误"))
		return
	}

	posts, pagination, err := h.service.ListPosts(c.Request.Context(), ListPostsParams{
		Cursor:   cursor,
		Limit:    parseIntDefault(c.Query("limit"), defaultPageLimit),
		Category: strings.TrimSpace(c.Query("category")),
		SortBy:   c.Query("sortBy"),
		Keyword:  strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithPagination(c, http.StatusOK, posts, pagination)
}

func (h *Handler) ListMyPosts(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	cursor, err := parseOptionalInt64(c.Query("cursor"))
	if err != nil {
		app.AbortError(c, badRequest("cursor 参数格式错误"))
		return
	}

	status, err := parsePostStatus(c.Query("status"))
	if err != nil {
		app.AbortError(c, err)
		return
	}

	posts, pagination, err := h.service.ListUserPosts(c.Request.Context(), *userID, UserPostsParams{
		Cursor: cursor,
		Limit:  parseIntDefault(c.Query("limit"), defaultPageLimit),
		Status: status,
	})
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithPagination(c, http.StatusOK, posts, pagination)
}

func (h *Handler) GetPost(c *gin.Context) {
	postID, err := parsePathID(c, "id", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	userID, err := app.ParseUserIDFromRequest(c.Request, h.jwtSecret)
	if err != nil {
		app.AbortError(c, &APIError{Status: http.StatusUnauthorized, Message: "登录信息无效", Err: err})
		return
	}

	post, err := h.service.GetPost(c.Request.Context(), postID, userID)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.Success(c, http.StatusOK, post)
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	var req UpsertPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.AbortError(c, badRequest("请求体格式错误"))
		return
	}

	post, err := h.service.CreatePost(c.Request.Context(), *userID, req)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusCreated, "帖子发布成功", post)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	postID, err := parsePathID(c, "id", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	var req UpsertPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.AbortError(c, badRequest("请求体格式错误"))
		return
	}

	post, err := h.service.UpdatePost(c.Request.Context(), postID, *userID, req)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusOK, "帖子更新成功", post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	postID, err := parsePathID(c, "id", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	if err := h.service.DeletePost(c.Request.Context(), postID, *userID); err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusOK, "帖子删除成功", nil)
}

func (h *Handler) ToggleLike(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	postID, err := parsePathID(c, "id", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	result, err := h.service.ToggleLike(c.Request.Context(), postID, *userID)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusOK, result.Message, gin.H{
		"like_count": result.LikeCount,
		"liked":      result.Liked,
	})
}

func (h *Handler) ListComments(c *gin.Context) {
	postID, err := parsePathID(c, "postID", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	comments, err := h.service.ListComments(c.Request.Context(), postID)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithCount(c, http.StatusOK, comments, len(comments))
}

func (h *Handler) AddComment(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	postID, err := parsePathID(c, "postID", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.AbortError(c, badRequest("请求体格式错误"))
		return
	}

	comment, err := h.service.AddComment(c.Request.Context(), postID, *userID, req)
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusCreated, "评论发布成功", comment)
}

func (h *Handler) ListRelatedPosts(c *gin.Context) {
	postID, err := parsePathID(c, "postID", "帖子 ID 格式错误")
	if err != nil {
		return
	}

	posts, err := h.service.ListRelatedPosts(c.Request.Context(), postID, parseIntDefault(c.Query("limit"), defaultRelatedNum))
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.Success(c, http.StatusOK, posts)
}

func (h *Handler) ListCategories(c *gin.Context) {
	stats, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		app.AbortError(c, err)
		return
	}

	app.Success(c, http.StatusOK, stats)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}

	commentID, err := parsePathID(c, "commentID", "评论 ID 格式错误")
	if err != nil {
		return
	}

	if err := h.service.DeleteComment(c.Request.Context(), commentID, *userID); err != nil {
		app.AbortError(c, err)
		return
	}

	app.SuccessWithMessage(c, http.StatusOK, "评论删除成功", nil)
}

func requireUserID(c *gin.Context, secret string) (*int64, bool) {
	userID, err := app.ParseUserIDFromRequest(c.Request, secret)
	if err != nil {
		app.AbortError(c, &APIError{Status: http.StatusUnauthorized, Message: "登录信息无效", Err: err})
		return nil, false
	}
	if userID == nil {
		app.AbortError(c, &APIError{Status: http.StatusUnauthorized, Message: "请先登录"})
		return nil, false
	}
	return userID, true
}

func parsePathID(c *gin.Context, name, message string) (int64, error) {
	value, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil {
		err = badRequest(message)
		app.AbortError(c, err)
		return 0, err
	}
	return value, nil
}

func parseOptionalInt64(raw string) (*int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func parseIntDefault(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func parsePostStatus(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "all":
		return "all", nil
	case "published":
		return "published", nil
	case "draft":
		return "draft", nil
	case "archived":
		return "archived", nil
	default:
		return "", badRequest("status 只支持 all、published、draft、archived")
	}
}

func badRequest(message string) error {
	return &APIError{Status: http.StatusBadRequest, Message: message}
}
