package jwtauth

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type handler struct {
	publicKey *rsa.PublicKey
}

func New(publicKey *rsa.PublicKey) *handler {
	return &handler{publicKey: publicKey}
}

func (h handler) ValidateJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path != "/ws" {
			return
		}

		tokenString, exists := extractAuthToken(context)
		if !exists {
			return
		}

		if !h.isTokenValid(context, tokenString) {
			return
		}
	}
}

func (h *handler) isTokenValid(context *gin.Context, tokenString string) bool {
	if context.Request.URL.Path == "/register" || context.Request.URL.Path == "/signin" {
		return false
	}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.publicKey, nil
	})

	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return false
	}

	return true
}

func extractAuthToken(context *gin.Context) (string, bool) {
	headerString := context.Query("Authtoken")
	if headerString == "" {
		context.AbortWithStatus(http.StatusUnauthorized)
		return "", false
	}

	return headerString, true
}
