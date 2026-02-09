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
		Used:       0, // 初始化使用次数为0
	}

	// 落库
	return s.repo.SaveCode(ctx, code)
}

func (s *ShareCodeService) GetRanking(ctx context.Context, sortBy string, limit int64) ([]*model.ShareCode, error) {
	return s.repo.GetList(ctx, sortBy, limit)
}

// GetShareCodesPaginated 分页查询分享码列表
func (s *ShareCodeService) GetShareCodesPaginated(ctx context.Context, req model.PaginationRequest) (*model.PaginationResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Sort == "" {
		req.Sort = "price"
	}

	// 调用repository获取分页数据
	data, total, err := s.repo.GetListWithPagination(ctx, req.Sort, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &model.PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// IncrementUsage 增加分享码使用次数
func (s *ShareCodeService) IncrementUsage(ctx context.Context, code string) error {
	return s.repo.IncrementUsage(ctx, code)
}

// GetShareCodeByCode 根据分享码获取详情
func (s *ShareCodeService) GetShareCodeByCode(ctx context.Context, code string) (*model.ShareCode, error) {
	return s.repo.GetByCode(ctx, code)
}
