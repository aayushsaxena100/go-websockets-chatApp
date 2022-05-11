package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	usernameToConnection map[string]*websocket.Conn
	connectionToUsername map[*websocket.Conn]string
)

type response struct {
	Packet string `json:"packet"`
}

type packet struct {
	To      string `json:"to,omitempty"`
	From    string `json:"from,omitempty"`
	Content string `json:"content"`
}

type server struct {
}

func New() *server {
	usernameToConnection = make(map[string]*websocket.Conn, 0)
	connectionToUsername = make(map[*websocket.Conn]string, 0)
	return &server{}
}

func (s *server) RegisterConnectionAndStartListening(ctx *gin.Context, username string) {
	connection, isSuccessful := upgradeAndGetWebhookConnection(ctx)
	if !isSuccessful {
		return
	}

	usernameToConnection[username] = connection
	connectionToUsername[connection] = username

	for {
		_, msgBytes, err := connection.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				// Need to remove from active connection
			}
			log.Println(err)
			return
		}

		msg := &packet{}
		_ = json.Unmarshal(msgBytes, msg)

		log.Println(fmt.Sprintf("Message received by server: %v", string(msgBytes)))

		sendToConnection := &websocket.Conn{}

		sendToConnection, ok := usernameToConnection[msg.To]
		if !ok {
			err = connection.WriteMessage(websocket.TextMessage, []byte("User not online or does not exist"))
			if err != nil {
				log.Println(err)
			}
			continue
		}

		from := connectionToUsername[connection]

		err = sendToConnection.WriteJSON(&packet{From: from, Content: msg.Content})
		if err != nil {
			log.Println(err)
		}
	}
}

func upgradeAndGetWebhookConnection(ctx *gin.Context) (con *websocket.Conn, isSuccessful bool) {
	isSuccessful = true

	upgrader := websocket.Upgrader{ReadBufferSize: 1042, WriteBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	con, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(fmt.Sprintf("Error in upgrading connection to websocket. Error: %v", err))

		isSuccessful = false
	}

	return con, isSuccessful
}
