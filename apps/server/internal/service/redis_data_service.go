package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fullstack-engineering-lab/server/internal/model"

	"github.com/redis/go-redis/v9"
)

// RedisDataService 提供 Hash / Set / Sorted Set 三类常用数据结构的操作
type RedisDataService struct {
	rdb *redis.Client
}

func NewRedisDataService(rdb *redis.Client) *RedisDataService {
	return &RedisDataService{rdb: rdb}
}

// ===== Hash: 用户画像缓存 =====

// SetHashField 设置单个字段
func (s *RedisDataService) SetHashField(req *model.HashFieldUpdateRequest) (*model.HashProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.rdb.HSet(ctx, req.Key, req.Field, req.Value).Err(); err != nil {
		return nil, fmt.Errorf("HSet 失败: %w", err)
	}

	return s.getHashProfile(ctx, req.Key)
}

// GetHashProfile 查询完整画像
func (s *RedisDataService) GetHashProfile(key string) (*model.HashProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.getHashProfile(ctx, key)
}

// MultiSetHash 批量设置字段
func (s *RedisDataService) MultiSetHash(req *model.HashMultiSetRequest) (*model.HashProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.rdb.HSet(ctx, req.Key, req.Fields).Err(); err != nil {
		return nil, fmt.Errorf("HMSet 失败: %w", err)
	}

	return s.getHashProfile(ctx, req.Key)
}

// DeleteHashField 删除字段
func (s *RedisDataService) DeleteHashField(key string, field string) (*model.HashProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	n, err := s.rdb.HDel(ctx, key, field).Result()
	if err != nil {
		return nil, fmt.Errorf("HDel 失败: %w", err)
	}
	if n == 0 {
		return nil, fmt.Errorf("字段 %q 不存在", field)
	}

	return s.getHashProfile(ctx, key)
}

// getHashProfile 内部方法：读取完整 Hash
func (s *RedisDataService) getHashProfile(ctx context.Context, key string) (*model.HashProfileResponse, error) {
	fields, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("HGetAll 失败: %w", err)
	}
	if len(fields) == 0 {
		return nil, errors.New("key 不存在或已过期")
	}

	return &model.HashProfileResponse{
		Key:      key,
		Fields:   fields,
		NumOfFld: len(fields),
	}, nil
}

// ===== Set: 标签/收藏夹管理 =====

// SetAddMembers 添加成员
func (s *RedisDataService) SetAddMembers(req *model.SetAddMembersRequest) (*model.SetListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := make([]interface{}, len(req.Members))
	for i, m := range req.Members {
		args[i] = m
	}

	if err := s.rdb.SAdd(ctx, req.Key, args...).Err(); err != nil {
		return nil, fmt.Errorf("SAdd 失败: %w", err)
	}

	return s.getSetMembers(ctx, req.Key)
}

// SetRemoveMembers 移除成员
func (s *RedisDataService) SetRemoveMembers(req *model.SetRemoveMembersRequest) (*model.SetListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := make([]interface{}, len(req.Members))
	for i, m := range req.Members {
		args[i] = m
	}

	s.rdb.SRem(ctx, req.Key, args...)
	return s.getSetMembers(ctx, req.Key)
}

// GetSetMembers 查询所有成员
func (s *RedisDataService) GetSetMembers(req *model.SetListRequest) (*model.SetListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.getSetMembers(ctx, req.Key)
}

// SetIntersect 计算交集
func (s *RedisDataService) SetIntersect(req *model.SetOperationRequest) (*model.SetOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	members, err := s.rdb.SInter(ctx, req.Keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("SInter 失败: %w", err)
	}

	return &model.SetOperationResponse{
		Keys:    req.Keys,
		Op:      "intersect",
		Members: members,
		Count:   len(members),
	}, nil
}

// SetUnion 计算并集
func (s *RedisDataService) SetUnion(req *model.SetOperationRequest) (*model.SetOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	members, err := s.rdb.SUnion(ctx, req.Keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("SUnion 失败: %w", err)
	}

	return &model.SetOperationResponse{
		Keys:    req.Keys,
		Op:      "union",
		Members: members,
		Count:   len(members),
	}, nil
}

// SetDiff 计算差集（以第一个集合为基准）
func (s *RedisDataService) SetDiff(req *model.SetOperationRequest) (*model.SetOperationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	members, err := s.rdb.SDiff(ctx, req.Keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("SDiff 失败: %w", err)
	}

	return &model.SetOperationResponse{
		Keys:    req.Keys,
		Op:      "diff",
		Members: members,
		Count:   len(members),
	}, nil
}

// getSetMembers 内部方法：读取 Set 所有成员
func (s *RedisDataService) getSetMembers(ctx context.Context, key string) (*model.SetListResponse, error) {
	members, err := s.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("SMembers 失败: %w", err)
	}

	return &model.SetListResponse{
		Key:     key,
		Members: members,
		Count:   len(members),
	}, nil
}

// ===== ZSet: 实时排行榜 =====

// ZSetAddScore 添加成员/增加分数（ZINCRBY）
func (s *RedisDataService) ZSetAddScore(req *model.ZSetAddScoreRequest) (*model.ZSetMemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newScore, err := s.rdb.ZIncrBy(ctx, req.Key, req.Score, req.Member).Result()
	if err != nil {
		return nil, fmt.Errorf("ZIncrBy 失败: %w", err)
	}

	// 查询最新排名
	rank, _ := s.rdb.ZRevRank(ctx, req.Key, req.Member).Result()

	return &model.ZSetMemberResponse{
		Key:    req.Key,
		Member: req.Member,
		Score:  newScore,
		Rank:   int(rank) + 1,
	}, nil
}

// ZSetGetTopN 获取前 N 名
func (s *RedisDataService) ZSetGetTopN(req *model.ZSetTopNRequest) (*model.ZSetTopNResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	total, err := s.rdb.ZCard(ctx, req.Key).Result()
	if err != nil {
		return nil, fmt.Errorf("ZCard 失败: %w", err)
	}

	results, err := s.rdb.ZRevRangeWithScores(ctx, req.Key, 0, int64(req.N-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("ZRevRangeWithScores 失败: %w", err)
	}

	members := make([]model.ZSetRankEntry, len(results))
	for i, z := range results {
		members[i] = model.ZSetRankEntry{
			Rank:   i + 1,
			Member: z.Member.(string),
			Score:  z.Score,
		}
	}

	return &model.ZSetTopNResponse{
		Key:     req.Key,
		Total:   int(total),
		Members: members,
	}, nil
}

// ZSetGetRank 查询某成员的排名和分数
func (s *RedisDataService) ZSetGetRank(req *model.ZSetRankRequest) (*model.ZSetMemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rank, err := s.rdb.ZRevRank(ctx, req.Key, req.Member).Result()
	if err != nil {
		return nil, errors.New("成员不存在")
	}

	score, err := s.rdb.ZScore(ctx, req.Key, req.Member).Result()
	if err != nil {
		return nil, errors.New("成员不存在")
	}

	return &model.ZSetMemberResponse{
		Key:    req.Key,
		Member: req.Member,
		Score:  score,
		Rank:   int(rank) + 1,
	}, nil
}

// ===== List: 最新活动流 / 简易消息队列 =====

// ListPush 推入消息到列表头部或尾部
func (s *RedisDataService) ListPush(req *model.ListPushRequest) (*model.ListRangeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var n int64
	var err error
	if req.Pos == "left" {
		n, err = s.rdb.LPush(ctx, req.Key, req.Value).Result()
	} else {
		n, err = s.rdb.RPush(ctx, req.Key, req.Value).Result()
	}
	if err != nil {
		return nil, fmt.Errorf("LPush/RPush 失败: %w", err)
	}

	// 返回最近 20 条记录
	return s.getListRange(ctx, req.Key, 0, 19, n)
}

// ListPop 从列表头部或尾部弹出消息
func (s *RedisDataService) ListPop(req *model.ListPopRequest) (*model.ListPopResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var val string
	var err error
	if req.Pos == "left" {
		val, err = s.rdb.LPop(ctx, req.Key).Result()
	} else {
		val, err = s.rdb.RPop(ctx, req.Key).Result()
	}
	if err != nil {
		return nil, fmt.Errorf("LPop/RPop 失败（列表可能为空）: %w", err)
	}

	return &model.ListPopResponse{
		Key:   req.Key,
		Value: val,
		Pos:   req.Pos,
	}, nil
}

// ListRange 查询列表指定范围的元素
func (s *RedisDataService) ListRange(req *model.ListRangeRequest) (*model.ListRangeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	length, err := s.rdb.LLen(ctx, req.Key).Result()
	if err != nil {
		return nil, fmt.Errorf("LLen 失败: %w", err)
	}

	return s.getListRange(ctx, req.Key, req.Start, req.Stop, length)
}

// getListRange 内部方法：获取列表范围
func (s *RedisDataService) getListRange(ctx context.Context, key string, start, stop int64, length int64) (*model.ListRangeResponse, error) {
	values, err := s.rdb.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("LRange 失败: %w", err)
	}

	return &model.ListRangeResponse{
		Key:    key,
		Values: values,
		Length: length,
	}, nil
}

// ===== String: 验证码存储 / 计数器 =====

// StringSet 设置键值（可选 TTL）
func (s *RedisDataService) StringSet(req *model.StringSetRequest) (*model.StringSetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if req.TTL > 0 {
		if err := s.rdb.Set(ctx, req.Key, req.Value, time.Duration(req.TTL)*time.Second).Err(); err != nil {
			return nil, fmt.Errorf("Set 失败: %w", err)
		}
	} else {
		if err := s.rdb.Set(ctx, req.Key, req.Value, 0).Err(); err != nil {
			return nil, fmt.Errorf("Set 失败: %w", err)
		}
	}

	return &model.StringSetResponse{
		Key:   req.Key,
		Value: req.Value,
		TTL:   req.TTL,
	}, nil
}

// StringGet 获取键值
func (s *RedisDataService) StringGet(key string) (*model.StringGetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return &model.StringGetResponse{
				Key:    key,
				Value:  "",
				TTL:    -2,
				Exists: false,
			}, nil
		}
		return nil, fmt.Errorf("Get 失败: %w", err)
	}

	ttl, _ := s.rdb.TTL(ctx, key).Result()
	ttlSec := int(ttl.Seconds())
	if ttlSec < 0 {
		ttlSec = -1 // -1 表示无过期
	}

	return &model.StringGetResponse{
		Key:    key,
		Value:  val,
		TTL:    ttlSec,
		Exists: true,
	}, nil
}

// StringIncr 自增计数器
func (s *RedisDataService) StringIncr(req *model.StringIncrRequest) (*model.StringIncrResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 先获取当前值作为 before
	beforeVal, err := s.rdb.Get(ctx, req.Key).Int64()
	if err != nil {
		if err == redis.Nil {
			beforeVal = 0 // key 不存在时，INCR 从 0 开始
		} else {
			return nil, fmt.Errorf("获取当前值失败: %w", err)
		}
	}

	delta := req.Delta
	if delta == 0 {
		delta = 1
	}

	var afterVal int64
	if delta > 0 {
		afterVal, err = s.rdb.IncrBy(ctx, req.Key, delta).Result()
	} else {
		afterVal, err = s.rdb.DecrBy(ctx, req.Key, -delta).Result()
	}
	if err != nil {
		return nil, fmt.Errorf("InrBy/DecrBy 失败: %w", err)
	}

	return &model.StringIncrResponse{
		Key:    req.Key,
		Before: beforeVal,
		After:  afterVal,
		Delta:  delta,
	}, nil
}
