package lock

import (
	"fmt"
	redis "github.com/go-redis/redis/v7"
	"math/rand"
	"sync/atomic"
	"time"
)

type RedisLock struct {
	store   *redis.Client
	seconds uint32
	count   int32
	key     string
	id      string
}

const (
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`

	SetnxExCommand = `if redis.call("GET", KEYS[1]) == "" then
		redis.call("SET", KEYS[1],ARGV[1])
		redis.call("EXPRIE", KEYS[1],ARGV[2])
    return 1
else
    return 0
end`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRedisLock(store *redis.Client, key string, id uint64) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    fmt.Sprint(rand.Uint64()),
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() (bool, error) {

	newCount := atomic.AddInt32(&rl.count, 1)
	if newCount > 1 {
		return true, nil
	}
	seconds := atomic.LoadUint32(&rl.seconds)
	ok, err := rl.store.Eval(SetnxExCommand, []string{rl.key}, []string{rl.id, fmt.Sprint(seconds + 1)}).Result()
	if err == redis.Nil {
		atomic.AddInt32(&rl.count, -1)
		return false, nil
	} else if err != nil {
		atomic.AddInt32(&rl.count, -1)
		return false, err
	} else if ok.(int) == 0 {
		atomic.AddInt32(&rl.count, -1)
		return false, nil
	}
	return true, nil
}

// Release releases the lock.
func (rl *RedisLock) Release() (bool, error) {
	newCount := atomic.AddInt32(&rl.count, -1)
	if newCount > 0 {
		return true, nil
	}
	resp, err := rl.store.Eval(delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
