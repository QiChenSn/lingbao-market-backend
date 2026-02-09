package service

import (
	"context"
	"errors"
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

	if price < 1 || price > 999 {
		return errors.New("价格无效，必须在 1-999 之间")
	}

	// 构建内部模型
	code := &model.ShareCode{
		Code:       req.Code,
		Price:      price,
		CreateTime: time.Now().Unix(),
	}

	// 落库
	return s.repo.SaveCode(ctx, code)
}

func (s *ShareCodeService) GetRanking(ctx context.Context, sortBy string) ([]*model.ShareCode, error) {
	// 限制返回前 100 条，防止大量数据拖慢接口
	return s.repo.GetList(ctx, sortBy, 100)
}
