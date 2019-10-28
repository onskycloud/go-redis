package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/onskycloud/structs"
	"github.com/vmihailenco/msgpack"
)

// NewRedis instance
func NewRedis(options interface{}) *Redis {
	opts := new(redis.Options)
	structs.Merge(opts, options)
	client := redis.NewClient(opts)
	return &Redis{
		db: client,
	}
}

// Redis struct
type Redis struct {
	db       *redis.Client
	pipeline *redis.Pipeline
}

// DB constructor
func (r *Redis) DB() *redis.Client {
	return r.db
}

// Close close redis
func (r *Redis) Close() error {
	return r.db.Close()
}

// Ping ping
func (r *Redis) Ping() *redis.StatusCmd {
	return r.db.Ping()
}

// Del del object by keys
func (r *Redis) Del(keys ...string) error {
	_, err := r.db.Del(keys...).Result()
	return err
}

// HDel del object by key and field
func (r *Redis) HDel(key string, fields ...string) error {
	_, err := r.db.HDel(key, fields...).Result()
	return err
}

// Get get object by key
func (r *Redis) Get(key string) (string, error) {
	return r.db.Get(key).Result()
}

// Set set object by key
func (r *Redis) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return r.db.Set(key, value, expiration).Result()
}

// SetObject set object by key and field
func (r *Redis) SetObject(objectKey string, field string, value interface{}) (bool, error) {
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.db.HSet(objectKey, field, bytes).Result()
}

// GetObject get object by key and field
func (r *Redis) GetObject(objectKey string, field string, result interface{}) error {
	temp, err := r.db.HGet(objectKey, field).Bytes()
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(temp, result)
	return err
}

// CheckExistedObject will return true if the object is existed.
func (r *Redis) CheckExistedObject(objectKey string, field string) (bool, error) {
	existed, err := r.db.HExists(objectKey, field).Result()
	if err != nil {
		return false, err
	}
	if !existed {
		return false, nil
	}
	return true, nil
}
