package store

import (
	"github.com/aayushsaxena100/testproject/models"
	"github.com/gin-gonic/gin"
)

type Client interface {
	Create(c *gin.Context, client *models.Client) (*models.Client, error)
	GetByUsername(c *gin.Context, username string) (*models.Client, error)
}
