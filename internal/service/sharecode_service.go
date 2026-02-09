package service

import (
	"context"
	"time"

	"lingbao-market-backend/internal/model"
	"lingbao-market-backend/internal/repository"
)

type ShareCodeService struct {
	repo *repository.ShareCodeRepo
}

func NewShareCodeService(repo *repository.ShareCodeRepo) *ShareCodeService {
	return &ShareCodeService{repo: repo}
}

func (s *ShareCodeService) AddShareCode(ctx context.Context, req model.CreateShareCodeRequest) error {
	price := req.Price

	// 构建内部模型
	code := &model.ShareCode{
		Code:       req.Code,
		Price:      price,
		CreateTime: time.Now().Unix(),
	}

	// 落库
	return s.repo.SaveCode(ctx, code)
}

func (s *ShareCodeService) GetRanking(ctx context.Context, sortBy string, limit int64) ([]*model.ShareCode, error) {
	return s.repo.GetList(ctx, sortBy, limit)
}
