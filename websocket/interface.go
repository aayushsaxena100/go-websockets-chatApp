package websocket

import "github.com/gin-gonic/gin"

type Server interface {
	RegisterConnectionAndStartListening(ctx *gin.Context, username string)
}
