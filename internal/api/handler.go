package api

import (
	"strconv"

	"lingbao-market-backend/internal/model"
	"lingbao-market-backend/internal/service"
	"lingbao-market-backend/pkg/response"

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
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 调用 Service
	if err := h.svc.AddShareCode(c.Request.Context(), req); err != nil {
		// 区分一下是业务错误还是系统错误（这里简单处理）
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "发布成功", nil)
}

// ListShareCodes GET /sharecode?sort=price|time&limit=100
func (h *Handler) ListShareCodes(c *gin.Context) {
	sort := c.DefaultQuery("sort", "price")
	limit_str := c.DefaultQuery("limit", "100")

	limit, err := strconv.ParseInt(limit_str, 10, 64)
	if err != nil {
		response.BadRequest(c, "limit 参数错误")
		return
	}

	list, err := h.svc.GetRanking(c.Request.Context(), sort, limit)
	if err != nil {
		response.InternalError(c, "获取列表失败")
		return
	}

	response.Success(c, list)
}

// ListShareCodesPage GET /sharecode/page?page=1&pageSize=10&sort=price|time
func (h *Handler) ListShareCodesPage(c *gin.Context) {
	var req model.PaginationRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 调用Service
	result, err := h.svc.GetShareCodesPaginated(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, "获取列表失败")
		return
	}

	response.Success(c, result)
}
