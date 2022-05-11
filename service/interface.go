package service

import (
	"github.com/aayushsaxena100/testproject/models"
	"github.com/gin-gonic/gin"
)

type Client interface {
	Register(ctx *gin.Context, client *models.Client) error
	Signin(ctx *gin.Context, client *models.Client) error
}
