package main

import (
	"crypto/rsa"
	"fmt"
	httpClient "github.com/aayushsaxena100/testproject/http/server"
	"github.com/aayushsaxena100/testproject/middleware/jwtauth"
	clientSvc "github.com/aayushsaxena100/testproject/service/client"
	clientStore "github.com/aayushsaxena100/testproject/store/client"
	websocketServer "github.com/aayushsaxena100/testproject/websocket/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	g := gin.Default()

	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	privateKey, err := loadRsaPrivateKey()
	if err != nil {
		return
	}

	publicKey, err := loadRsaPublicKey()
	if err != nil {
		return
	}

	webSocketServer := websocketServer.New()
	cltStore := clientStore.New(rds)
	clientService := clientSvc.New(cltStore)
	httpClt := httpClient.New(clientService, webSocketServer, privateKey)

	addJWTAuth(g, publicKey)

	g.Static("public", "./public")
	g.LoadHTMLGlob("public/html/*")

	g.GET("/", serveRegisterPage)
	g.GET("/home", serveHomePage)
	g.GET("/signin", serveSigninPage)
	g.POST("/signin", httpClt.Signin)
	g.POST("/register", httpClt.Register)
	g.GET("/ws", httpClt.WebSocketHandler)

	log.Println("Listening on port 8080")

	err = g.Run("localhost:8080")
	if err != nil {
		log.Println(fmt.Sprintf("Error starting server. Error: %v", err))
		return
	}
}

func addJWTAuth(g *gin.Engine, publicKey *rsa.PublicKey) {
	g.Use(jwtauth.New(publicKey).ValidateJWT())
}

func serveHomePage(context *gin.Context) {
	context.HTML(200, "home.html", nil)
}

func serveSigninPage(context *gin.Context) {
	context.HTML(200, "signin.html", nil)
}

func serveRegisterPage(context *gin.Context) {
	context.HTML(200, "register.html", nil)
}

func loadRsaPrivateKey() (*rsa.PrivateKey, error) {
	cwd, _ := os.Getwd()
	bytes, err := ioutil.ReadFile(cwd + "/jwtRS256.key")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in fetching PRIVATE KEY. Error: %v", err))
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(bytes)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in parsing PRIVATE KEY. Error: %v", err))
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}

func loadRsaPublicKey() (*rsa.PublicKey, error) {
	cwd, _ := os.Getwd()
	bytes, err := ioutil.ReadFile(cwd + "/jwtRS256.key.pub")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in fetching PUBLIC KEY. Error: %v", err))
		return nil, err
	}

	key, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in parsing PUBLIC KEY. Error: %v", err))
		return nil, err
	}

	parsedCryptoKey := key.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	return pubCrypto.(*rsa.PublicKey), nil
}
