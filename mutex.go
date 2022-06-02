package ez3

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname)

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}

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
