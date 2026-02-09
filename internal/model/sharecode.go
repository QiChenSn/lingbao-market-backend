package model

type CreateShareCodeRequest struct {
	Code  string `json:"code" binding:"required"`
	Price int    `json:"price" binding:"required"` // 前端传字符串，我们在 Service 层转为 int
}

// ShareCode 存入 Redis Hash 中的完整数据模型
type ShareCode struct {
	Code       string `json:"code"`
	Price      int    `json:"price"`
	CreateTime int64  `json:"create_time"` // 时间戳
}

// ShareCodeResponse 用于返回给前端的结构
type ShareCodeResponse struct {
	Code       string `json:"code"`
	Price      int    `json:"price"`
	CreateTime int64  `json:"create_time"`
}
