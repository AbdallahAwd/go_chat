package repositories

import (
	"time"

	"github.com/go-redis/redis"
)

type CacheRepo struct {
	client *redis.Client
}

func NewCacheRepo(client *redis.Client) *CacheRepo {
	return &CacheRepo{client: client}
}
func (c *CacheRepo) SetAsBlock(ip string) error {
	key := "Block" + ip
	return c.client.Set(key, ip, time.Minute*5).Err()
}
func (c *CacheRepo) IsIPBlocked(ip string) (bool, error) {
	key := "Block" + ip
	err := c.client.Get(key).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
