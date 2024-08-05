package handlers

import (
	"chat_app/internal/services"
	"chat_app/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	service   *services.ChatService
	JwtSecret string
	clients   map[uint]*websocket.Conn
	mutex     sync.Mutex
	upgrader  websocket.Upgrader
}
type Message struct {
	Type        string `json:"type"`
	UserId      uint   `json:"user_id"`
	RecipientID uint   `json:"recipient_id"`
	Content     string `json:"content"`
	AudioURL    string `json:"audio_url,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

func NewChathandler(s *services.ChatService, jwtSecret string) *ChatHandler {
	return &ChatHandler{service: s,
		clients:   make(map[uint]*websocket.Conn),
		JwtSecret: jwtSecret,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Minute * 2,
			CheckOrigin: func(r *http.Request) bool {

				return true
			},
		},
	}
}

func (c *ChatHandler) ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	userId, err := c.authenticateUser(r)

	if err != nil {
		log.Println(err.Error())
		return
	}

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.registerClinet(uint(userId), conn)

	go c.handleMessage(uint(userId), conn)
}

func (c *ChatHandler) registerClinet(userId uint, conn *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.clients[userId] = conn
}

func (c *ChatHandler) unRegisterClinet(userId uint) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if conn, ok := c.clients[userId]; !ok {
		conn.Close()
		delete(c.clients, userId)
	}
}

func (c *ChatHandler) handleMessage(userId uint, conn *websocket.Conn) {
	defer c.unRegisterClinet(userId)

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("read error:", err)
			break
		}
		msg.UserId = userId
		c.handleChatMessage(userId, &msg)
		// switch msg.Type {
		// case "chat":
		// 	c.handleChatMessage(userId, &msg)
		// case "read":
		// 	c.handleReadMessage(userId, &msg)
		// }
	}
}

func (c *ChatHandler) handleChatMessage(senderID uint, msg *Message) {
	savedMsg, err := c.service.SendMessage(msg.Content, msg.AudioURL, msg.ImageURL, senderID, msg.RecipientID)
	if err != nil {
		log.Println("Error saving message:", err)
		return
	}
	c.mutex.Lock()
	if conn, ok := c.clients[msg.RecipientID]; ok {
		err = conn.WriteJSON(savedMsg)
		if err != nil {
			log.Println("Error sending message to recipient:", err)
		}
	}
	c.mutex.Unlock()

}
func (c *ChatHandler) authenticateUser(r *http.Request) (uint, error) {
	authToken := r.Header.Get("token")
	if authToken == "" {
		return 0, fmt.Errorf("authorization header is required")
	}
	authToken, ok := strings.CutPrefix(authToken, "Bearer ")
	if !ok {
		return 0, fmt.Errorf("invalid Token Format")
	}

	// Parse the token
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.JwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	// Extract user ID from claims
	log.Println(claims["ID"])
	userIDFloat, ok := claims["ID"].(float64)
	if !ok {
		return 0, fmt.Errorf("user ID type assertion")
	}
	userID := uint(userIDFloat)

	return userID, nil
}
func (c *ChatHandler) GetMessagedUsers(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(utils.ID).(float64)
	if !ok {
		utils.ErrorJSON(w, "Can not Convert userId", http.StatusInternalServerError)

		return
	}
	users, err := c.service.GetConversation(uint(userId))
	if err != nil {
		utils.ErrorJSON(w, "Conversation Error", http.StatusBadRequest)
		return
	}

	utils.ResponseHandler(w, users)
}

func (c *ChatHandler) GetChatBetweenTwoUsers(w http.ResponseWriter, r *http.Request) {
	userId1, ok := r.Context().Value(utils.ID).(float64)
	if !ok {
		utils.ErrorJSON(w, "Can not Convert userId", http.StatusInternalServerError)

		return
	}
	userID2, err := strconv.ParseUint(r.URL.Query().Get("with"), 10, 32)
	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	messages, err := c.service.GetChatPartners(uint(userId1), uint(userID2), int(limit), 0)
	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseHandler(w, messages)
}
