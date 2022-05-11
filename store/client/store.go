package client

import (
	"github.com/aayushsaxena100/testproject/models"
	store2 "github.com/aayushsaxena100/testproject/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const clientsHashSet = "clientsHashSet"

type store struct {
	rds *redis.Client
}

func New(rds *redis.Client) store2.Client {
	return &store{rds: rds}
}

func (s *store) Create(c *gin.Context, client *models.Client) (*models.Client, error) {
	_, err := s.rds.HSet(c, clientsHashSet, map[string]string{client.Username: client.Password}).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *store) GetByUsername(c *gin.Context, username string) (*models.Client, error) {
	bytes, err := s.rds.HGet(c, clientsHashSet, username).Bytes()
	if err != nil {
		if err.Error() != "redis: nil" {
			return nil, err
		}
		return nil, nil
	}

	var client = &models.Client{Username: username, Password: string(bytes)}

	return client, nil
}
