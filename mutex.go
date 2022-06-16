package ez3

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

// Mutex guards against concurrent write access to an EZ3 client using Redis as a source of truth.
type Mutex struct {
	ez3     *EZ3
	redsync *redsync.Redsync
}

// NewMutex creates a new Mutex instance, wrapping around an EZ3 instance and a Redis connection.
func NewMutex(ez3 *EZ3, redisClient *redis.Client) Mutex {
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)
	return Mutex{ez3: ez3, redsync: rs}
}

func (m Mutex) Get(key string, dst Serializable) error {
	return (*m.ez3).Get(key, dst)
}

func (m Mutex) Set(key string, val Serializable) error {
	mx := m.redsync.NewMutex(key)
	if err := mx.Lock(); err != nil {
		return err
	}
	defer mx.Unlock()
	return (*m.ez3).Set(key, val)
}

func (m Mutex) Del(key string) error {
	mx := m.redsync.NewMutex(key)
	if err := mx.Lock(); err != nil {
		return err
	}
	defer mx.Unlock()
	return (*m.ez3).Del(key)
}

func (m Mutex) List(prefix string) (keys []string, err error) {
	return (*m.ez3).List(prefix)
}
