package community

import (
	"net/http"
	"strconv"

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
	group.GET("/posts/:id", h.GetPost)
	group.POST("/posts", h.CreatePost)
	group.PUT("/posts/:id", h.UpdatePost)
	group.DELETE("/posts/:id", h.DeletePost)
	group.POST("/posts/:id/like", h.ToggleLike)
	group.GET("/posts/:postID/comments", h.ListComments)
	group.POST("/posts/:postID/comments", h.AddComment)
	group.GET("/posts/related/:postID", h.ListRelatedPosts)
	group.DELETE("/comments/:commentID", h.DeleteComment)
}

func (h *Handler) ListPosts(c *gin.Context) {
	posts, pagination, err := h.service.ListPosts(c.Request.Context(), ListPostsParams{
		Cursor:   parseIntPointer(c.Query("cursor")),
		Limit:    parseIntDefault(c.Query("limit"), 10),
		Category: c.Query("category"),
		SortBy:   c.DefaultQuery("sortBy", "latest"),
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": posts, "pagination": pagination})
}

func (h *Handler) GetPost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	userID, err := app.ParseUserIDFromRequest(c.Request, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token"})
		return
	}
	post, err := h.service.GetPost(c.Request.Context(), postID, userID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": post})
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	var req UpsertPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}
	post, err := h.service.CreatePost(c.Request.Context(), *userID, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "帖子发布成功", "data": post})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	var req UpsertPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}
	post, err := h.service.UpdatePost(c.Request.Context(), postID, *userID, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "帖子更新成功", "data": post})
}

func (h *Handler) DeletePost(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	if err := h.service.DeletePost(c.Request.Context(), postID, *userID); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "帖子删除成功"})
}

func (h *Handler) ToggleLike(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	result, err := h.service.ToggleLike(c.Request.Context(), postID, *userID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": result.Message, "like_count": result.LikeCount, "liked": result.Liked})
}

func (h *Handler) ListComments(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	comments, err := h.service.ListComments(c.Request.Context(), postID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": comments, "count": len(comments)})
}

func (h *Handler) AddComment(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	postID, err := strconv.ParseInt(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}
	comment, err := h.service.AddComment(c.Request.Context(), postID, *userID, req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "评论发布成功", "data": comment})
}

func (h *Handler) ListRelatedPosts(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帖子ID无效"})
		return
	}
	posts, err := h.service.ListRelatedPosts(c.Request.Context(), postID, parseIntDefault(c.Query("limit"), 2))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": posts})
}

func (h *Handler) DeleteComment(c *gin.Context) {
	userID, ok := requireUserID(c, h.jwtSecret)
	if !ok {
		return
	}
	commentID, err := strconv.ParseInt(c.Param("commentID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "评论ID无效"})
		return
	}
	if err := h.service.DeleteComment(c.Request.Context(), commentID, *userID); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "评论删除成功"})
}

func requireUserID(c *gin.Context, secret string) (*int64, bool) {
	userID, err := app.ParseUserIDFromRequest(c.Request, secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token"})
		return nil, false
	}
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing Authorization Header"})
		return nil, false
	}
	return userID, true
}

func writeError(c *gin.Context, err error) {
	apiErr, ok := err.(*APIError)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(apiErr.Status, gin.H{"message": apiErr.Message})
}

func parseIntPointer(raw string) *int64 {
	if raw == "" {
		return nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil
	}
	return &value
}

func parseIntDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
