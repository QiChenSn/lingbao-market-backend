package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lingbao-market-backend/internal/model"

	"github.com/redis/go-redis/v9"
)

type ShareCodeRepo struct {
	rdb *redis.Client
}

func NewShareCodeRepo(rdb *redis.Client) *ShareCodeRepo {
	return &ShareCodeRepo{rdb: rdb}
}

// SaveCode 使用 Pipeline 原子性地写入数据和索引
func (r *ShareCodeRepo) SaveCode(ctx context.Context, data *model.ShareCode) error {
	// 按天生成 Key，每天自动隔离
	dateStr := time.Now().Format("2006-01-02")
	keyData := fmt.Sprintf("market:data:%s", dateStr)          // Hash: 存详情
	keyIdxPrice := fmt.Sprintf("market:idx:price:%s", dateStr) // ZSet: 按价格排
	keyIdxTime := fmt.Sprintf("market:idx:time:%s", dateStr)   // ZSet: 按时间排

	jsonData, _ := json.Marshal(data)
	expireTime := 26 * time.Hour // 设置略多于1天的过期时间

	// 开启 Pipeline
	_, err := r.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// 1. 存入详情
		pipe.HSet(ctx, keyData, data.Code, jsonData)

		// 2. 更新价格索引 (Score = Price)
		pipe.ZAdd(ctx, keyIdxPrice, redis.Z{
			Score:  float64(data.Price), // ZSet Score 必须是 float
			Member: data.Code,
		})

		// 3. 更新时间索引 (Score = Timestamp)
		pipe.ZAdd(ctx, keyIdxTime, redis.Z{
			Score:  float64(data.CreateTime),
			Member: data.Code,
		})

		// 4. 设置过期
		pipe.Expire(ctx, keyData, expireTime)
		pipe.Expire(ctx, keyIdxPrice, expireTime)
		pipe.Expire(ctx, keyIdxTime, expireTime)
		return nil
	})

	return err
}

// GetList 根据 sortType 获取列表
func (r *ShareCodeRepo) GetList(ctx context.Context, sortType string, limit int64) ([]*model.ShareCode, error) {
	dateStr := time.Now().Format("2006-01-02")
	keyData := fmt.Sprintf("market:data:%s", dateStr)

	// 决定查询哪个索引
	var targetIdx string
	if sortType == "time" {
		targetIdx = fmt.Sprintf("market:idx:time:%s", dateStr)
	} else {
		targetIdx = fmt.Sprintf("market:idx:price:%s", dateStr)
	}

	// 1. 从 ZSet 倒序取前 N 个 Code
	codes, err := r.rdb.ZRevRange(ctx, targetIdx, 0, limit-1).Result()
	if err != nil || len(codes) == 0 {
		return []*model.ShareCode{}, nil
	}

	// 2. 拿着 Code 去 Hash 取详情
	jsonStrings, err := r.rdb.HMGet(ctx, keyData, codes...).Result()
	if err != nil {
		return nil, err
	}

	// 3. 反序列化
	var result []*model.ShareCode
	for _, s := range jsonStrings {
		if s == nil {
			continue
		}
		if str, ok := s.(string); ok {
			var item model.ShareCode
			_ = json.Unmarshal([]byte(str), &item)
			result = append(result, &item)
		}
	}
	return result, nil
}

// GetListWithPagination 分页查询
func (r *ShareCodeRepo) GetListWithPagination(ctx context.Context, sortType string, page, pageSize int) ([]*model.ShareCode, int64, error) {
	dateStr := time.Now().Format("2006-01-02")
	keyData := fmt.Sprintf("market:data:%s", dateStr)

	// 决定查询哪个索引
	var targetIdx string
	if sortType == "time" {
		targetIdx = fmt.Sprintf("market:idx:time:%s", dateStr)
	} else {
		targetIdx = fmt.Sprintf("market:idx:price:%s", dateStr)
	}

	// 1. 获取总数
	total, err := r.rdb.ZCard(ctx, targetIdx).Result()
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.ShareCode{}, 0, nil
	}

	// 2. 计算分页范围
	start := int64((page - 1) * pageSize)
	stop := start + int64(pageSize) - 1

	// 3. 从 ZSet 倒序取指定范围的 Code
	codes, err := r.rdb.ZRevRange(ctx, targetIdx, start, stop).Result()
	if err != nil || len(codes) == 0 {
		return []*model.ShareCode{}, total, nil
	}

	// 4. 拿着 Code 去 Hash 取详情
	jsonStrings, err := r.rdb.HMGet(ctx, keyData, codes...).Result()
	if err != nil {
		return nil, 0, err
	}

	// 5. 反序列化
	var result []*model.ShareCode
	for _, s := range jsonStrings {
		if s == nil {
			continue
		}
		if str, ok := s.(string); ok {
			var item model.ShareCode
			_ = json.Unmarshal([]byte(str), &item)
			result = append(result, &item)
		}
	}
	return result, total, nil
}
