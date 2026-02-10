package model

import "errors"

// 业务错误定义
var (
	ErrPriceOutOfRange = errors.New("价格必须在700-980范围内")
)

type CreateShareCodeRequest struct {
	Code  string `json:"code" binding:"required"`
	Price int    `json:"price" binding:"required,min=700,max=980"`
}

// ShareCode 存入 Redis Hash 中的完整数据模型
type ShareCode struct {
	Code       string `json:"code"`
	Price      int    `json:"price"`
	CreateTime int64  `json:"create_time"` // 时间戳
	Used       int    `json:"used"`        // 使用次数
}

// ShareCodeResponse 用于返回给前端的结构
type ShareCodeResponse struct {
	Code       string `json:"code"`
	Price      int    `json:"price"`
	CreateTime int64  `json:"create_time"`
	Used       int    `json:"used"` // 使用次数
}

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`             // 页码，从1开始
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"` // 每页条数
	Sort     string `form:"sort"`                                       // 排序方式
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Data       interface{} `json:"data"`       // 数据列表
	Total      int64       `json:"total"`      // 总条数
	Page       int         `json:"page"`       // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页条数
	TotalPages int         `json:"totalPages"` // 总页数
}
