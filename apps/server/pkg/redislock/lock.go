package redislock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotAcquired = errors.New("未获取到锁")
	ErrLockNotHeld     = errors.New("当前持有者不拥有该锁")
)

const lockPrefix = "lock:"

// Lock 表示一个 Redis 分布式锁实例
type Lock struct {
	rdb    *redis.Client
	key    string
	value  string
	ttl    time.Duration
	locked bool
}

// New 创建一个新的分布式锁实例
// value 使用随机 Token，确保只有锁持有者能释放
func New(rdb *redis.Client, resource string, ttl time.Duration) *Lock {
	return &Lock{
		rdb:   rdb,
		key:   lockPrefix + resource,
		value: randomToken(),
		ttl:   ttl,
	}
}

func randomToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	)
}

// Acquire 尝试通过 SET NX EX 获取锁，如果锁已被占用则返回 ErrLockNotAcquired
func (l *Lock) Acquire(ctx context.Context) error {
	ok, err := l.rdb.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrLockNotAcquired
	}
	l.locked = true
	return nil
}

// Release 通过 Lua 脚本原子释放锁，先比较 value 再删除，防止误删其他客户端的锁
func (l *Lock) Release(ctx context.Context) error {
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, l.rdb, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotHeld
	}
	l.locked = false
	return nil
}

// IsLocked 返回当前实例是否持有锁
func (l *Lock) IsLocked() bool {
	return l.locked
}

// Owner 返回锁持有者的唯一标识
func (l *Lock) Owner() string {
	return l.value
}

// TTL 返回锁的过期时间
func (l *Lock) TTL() time.Duration {
	return l.ttl
}

// Key 返回锁在 Redis 中的 key
func (l *Lock) Key() string {
	return l.key
}

// GetLockInfo 查询指定资源的锁在 Redis 中的状态
func GetLockInfo(ctx context.Context, rdb *redis.Client, resource string) (owner string, ttl time.Duration, err error) {
	key := lockPrefix + resource

	owner, err = rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, err
	}

	ttl, err = rdb.TTL(ctx, key).Result()
	if err != nil {
		return owner, 0, err
	}

	return owner, ttl, nil
}

// StatusResponse 锁状态查询的响应结构
type StatusResponse struct {
	Resource string `json:"resource"`
	Locked   bool   `json:"locked"`
	Owner    string `json:"owner,omitempty"`
	TTL      int64  `json:"ttl_ms"`
}
