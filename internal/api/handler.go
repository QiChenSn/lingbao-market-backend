package api

import (
	"net/http"

	"lingbao-market-backend/internal/model"
	"lingbao-market-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.ShareCodeService
}

func NewHandler(svc *service.ShareCodeService) *Handler {
	return &Handler{svc: svc}
}

// CreateShareCode POST /sharecode
func (h *Handler) CreateShareCode(c *gin.Context) {
	var req model.CreateShareCodeRequest
	// 绑定 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 调用 Service
	if err := h.svc.AddShareCode(c.Request.Context(), req); err != nil {
		// 区分一下是业务错误还是系统错误（这里简单处理）
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功"})
}

// ListShareCodes GET /sharecode?sort=price|time
func (h *Handler) ListShareCodes(c *gin.Context) {
	// 获取 query 参数，默认按价格排序
	sort := c.DefaultQuery("sort", "price")

	list, err := h.svc.GetRanking(c.Request.Context(), sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
