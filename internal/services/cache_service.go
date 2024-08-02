package services

import "chat_app/internal/repositories"

type CacheService struct {
	repo *repositories.CacheRepo
}

func NewCacheService(repo *repositories.CacheRepo) *CacheService {
	return &CacheService{repo: repo}
}

func (c *CacheService) SetAsBlock(ip string) error {
	return c.repo.SetAsBlock(ip)
}
func (c *CacheService) IsIPBlocked(ip string) (bool, error) {
	return c.repo.IsIPBlocked(ip)
}
