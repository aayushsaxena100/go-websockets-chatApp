package server

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/aayushsaxena100/testproject/models"
	"github.com/aayushsaxena100/testproject/service"
	"github.com/aayushsaxena100/testproject/websocket"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type server struct {
	service         service.Client
	key             *rsa.PrivateKey
	webSocketServer websocket.Server
}

// Claims is a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type response struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type user struct {
	Name string `json:"username"`
}

func New(svc service.Client, websocketServer websocket.Server, key *rsa.PrivateKey) *server {
	return &server{service: svc, webSocketServer: websocketServer, key: key}
}

func (s *server) Register(ctx *gin.Context) {
	user := &models.Client{}
	err := ctx.Bind(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling user model. Error: %v", err))
		ctx.JSON(http.StatusBadRequest, response{
			Message: "Invalid request body",
		})
		return
	}

	err = s.service.Register(ctx, user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error registering a user. Error: %v", err))
		ctx.JSON(http.StatusInternalServerError, response{
			Message: "Some error",
		})
		return
	}
}

func (s *server) Signin(ctx *gin.Context) {
	user := &models.Client{}
	err := ctx.Bind(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling user model. Error: %v", err))
		ctx.JSON(http.StatusBadRequest, response{
			Message: "Invalid request body",
		})
		return
	}

	err = s.service.Signin(ctx, user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in signin. Error: %v", err))
		ctx.JSON(http.StatusInternalServerError, response{
			Message: "Some error",
		})
		return
	}

	//Generate jwt-token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	})

	signedToken, err := token.SignedString(s.key)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in signing jwt token. Error: %v", err))
		ctx.JSON(http.StatusInternalServerError, response{
			Message: "Some error",
		})
		return
	}

	//ctx.JSON(http.StatusOK, response{
	//	Message: "Signin Successful",
	//	Token:   signedToken,
	//})

	ctx.HTML(http.StatusOK, "home.html", gin.H{"token": signedToken})
}

func (s *server) WebSocketHandler(ctx *gin.Context) {
	tokenString := ctx.Query("Authtoken")
	jwtData := strings.Split(tokenString, ".")

	payload, _ := jwt.DecodeSegment(jwtData[1])
	u := &user{}
	_ = json.Unmarshal(payload, u)

	go func() {
		s.webSocketServer.RegisterConnectionAndStartListening(ctx, u.Name)
	}()
}
