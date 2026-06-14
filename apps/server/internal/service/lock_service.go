package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/pkg/redislock"

	"github.com/redis/go-redis/v9"
)

type LockService struct {
	rdb *redis.Client
}

func NewLockService(rdb *redis.Client) *LockService {
	return &LockService{rdb: rdb}
}

func (s *LockService) AcquireLock(req *model.LockAcquireRequest) (*model.LockAcquireResponse, error) {
	if s.rdb == nil {
		return nil, errors.New("Redis 不可用")
	}

	ttl := time.Duration(req.TTL) * time.Second
	lock := redislock.New(s.rdb, req.Resource, ttl)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := lock.Acquire(ctx); err != nil {
		if errors.Is(err, redislock.ErrLockNotAcquired) {
			// 获取冲突时返回当前锁持有信息
			owner, lockTTL, _ := redislock.GetLockInfo(ctx, s.rdb, req.Resource)
			return nil, fmt.Errorf("资源 %q 已被 %s 锁定（TTL: %dms）", req.Resource, owner[:8], lockTTL.Milliseconds())
		}
		return nil, err
	}

	return &model.LockAcquireResponse{
		Resource: req.Resource,
		Owner:    lock.Owner(),
		TTL:      req.TTL,
	}, nil
}

func (s *LockService) ReleaseLock(req *model.LockReleaseRequest) error {
	if s.rdb == nil {
		return errors.New("Redis 不可用")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 先检查锁是否存在
	_, _, err := redislock.GetLockInfo(ctx, s.rdb, req.Resource)
	if err != nil {
		return err
	}

	// 演示用途：获取当前持有者并强制释放
	// 生产环境中每个客户端应持有自己的 Lock 实例
	owner, _, err := redislock.GetLockInfo(ctx, s.rdb, req.Resource)
	if err != nil {
		return err
	}
	if owner == "" {
		return errors.New("锁不存在")
	}

	// 演示场景直接 DEL 强制释放
	key := "lock:" + req.Resource
	return s.rdb.Del(ctx, key).Err()
}

func (s *LockService) GetStatus(req *model.LockStatusRequest) (*model.LockStatusResponse, error) {
	if s.rdb == nil {
		return nil, errors.New("Redis 不可用")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	owner, ttl, err := redislock.GetLockInfo(ctx, s.rdb, req.Resource)
	if err != nil {
		return nil, err
	}

	return &model.LockStatusResponse{
		Resource: req.Resource,
		Locked:   owner != "",
		Owner:    owner,
		TTL:      ttl.Milliseconds(),
	}, nil
}

func (s *LockService) ContentionDemo(req *model.ContentionDemoRequest) (*model.ContentionDemoResponse, error) {
	if s.rdb == nil {
		return nil, errors.New("Redis 不可用")
	}

	// 先清除资源上已有的锁
	ctx := context.Background()
	s.rdb.Del(ctx, "lock:"+req.Resource)

	ttl := time.Duration(req.TTL) * time.Second
	holdDuration := time.Duration(req.HoldMs) * time.Millisecond

	var (
		mu      sync.Mutex
		results = make([]model.ContentionResult, req.Goroutines)
		wg      sync.WaitGroup
	)

	for i := 0; i < req.Goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			lock := redislock.New(s.rdb, req.Resource, ttl)
			start := time.Now()

			err := lock.Acquire(ctx)
			waitMs := time.Since(start).Milliseconds()

			if err != nil {
				mu.Lock()
				results[id] = model.ContentionResult{
					GoroutineID: id + 1,
					Acquired:    false,
					WaitMs:      waitMs,
					Message:     "获取锁失败：资源已被占用",
				}
				mu.Unlock()
				return
			}

			// 模拟持有锁期间的业务处理
			time.Sleep(holdDuration)

			// 释放锁
			_ = lock.Release(ctx)

			mu.Lock()
			results[id] = model.ContentionResult{
				GoroutineID: id + 1,
				Acquired:    true,
				WaitMs:      waitMs,
				Message:     fmt.Sprintf("获取锁成功，持有 %dms 后释放", req.HoldMs),
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	succeeded := 0
	for _, r := range results {
		if r.Acquired {
			succeeded++
		}
	}

	return &model.ContentionDemoResponse{
		Resource: req.Resource,
		Results:  results,
		Summary: model.ContentionSummary{
			Total:     req.Goroutines,
			Succeeded: succeeded,
			Failed:    req.Goroutines - succeeded,
		},
	}, nil
}
