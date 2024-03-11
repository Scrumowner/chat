package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	authservice "server/internal/modules/auth/service"
	"server/internal/modules/chat/models"
	"server/internal/modules/chat/service"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ChatController struct {
	Hub  *service.Hub
	Auth *authservice.AuthService
}

func NewChatController(h *service.Hub, db *sqlx.DB, token []byte) *ChatController {
	return &ChatController{
		Hub:  h,
		Auth: authservice.NewAuthService(db, token),
	}
}
func (c *ChatController) App(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "appv2.html", gin.H{
		"title": "Main website",
	})

}
func (c *ChatController) CreateRoom(ctx *gin.Context) {
	var req models.CreateRoomRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Hub.Rooms[req.ID] = &service.Room{
		Name:    req.Name,
		ID:      req.ID,
		Clients: make(map[string]*service.Client),
	}
	ctx.JSON(http.StatusOK, req)
}
func (c *ChatController) GetRooms(ctx *gin.Context) {
	rooms := make([]models.GetRoomsResponse, 0)
	for _, room := range c.Hub.Rooms {
		rooms = append(rooms, models.GetRoomsResponse{
			ID:   room.ID,
			Name: room.Name,
		})
	}
	ctx.JSON(http.StatusOK, rooms)
}
func (c *ChatController) GetClients(ctx *gin.Context) {
	clients := make([]models.GetClientsResponse, 0)
	roomid := ctx.Param("roomid")
	_, ok := c.Hub.Rooms[roomid]
	if !ok {
		ctx.JSON(http.StatusOK, clients)
	}
	for _, client := range c.Hub.Rooms[roomid].Clients {
		clients = append(clients, models.GetClientsResponse{
			ID:       client.ID,
			Username: client.Username,
		})
	}
	ctx.JSON(http.StatusOK, clients)
}

func (c *ChatController) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	roomid := ctx.Param("roomid")
	clientid := ctx.Query("clientid")
	username := ctx.Query("username")
	msg := &service.Message{
		Content:  fmt.Sprintf("%s are entry", username),
		RoomID:   roomid,
		Username: username,
	}
	client := &service.Client{
		ID:       clientid,
		Username: username,
		RoomID:   roomid,
		Conn:     conn,
		Message:  make(chan *service.Message, 10),
	}

	c.Hub.Register <- client
	c.Hub.Broadcast <- msg
	go client.Read(c.Hub)
	client.Write()

}
