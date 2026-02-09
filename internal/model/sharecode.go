package model

type CreateShareCodeRequest struct {
	Code  string `json:"code" binding:"required,min=10,max=20"`
	Price int    `json:"price" binding:"required,min=1,max=999"`
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
