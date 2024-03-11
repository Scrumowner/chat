package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	controller "server/internal/modules/chat/contorller"
	"strings"
)

type Router struct {
	c *controller.ChatController
}

func NewRouter(controller *controller.ChatController) *Router {
	return &Router{
		c: controller,
	}
}
func (r *Router) Route(router *gin.Engine) {
	///home/host/go/src/telegram-chat/server/static/appv2.html
	router.LoadHTMLFiles("./appv2.html")
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()

	})
	router.Use(func(ctx *gin.Context) {
		if strings.HasPrefix("/ws/create", ctx.Request.URL.Path) || strings.HasPrefix("/ws/clients", ctx.Request.URL.Path) || strings.HasPrefix("/ws/rooms", ctx.Request.URL.Path) {
			raw := ctx.Request.Header.Get("Authorization")
			token := strings.TrimSpace(raw)
			isValid := r.c.Auth.Verify(token)
			if !isValid {
				log.Println("JWT TOKEN IS INVALID", token)
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.Next()
		}
	})
	auth := router.Group("/auth")
	//
	//
	//	http://localhost:4321/auth/register
	//	{"username":"Peter","password":"123123"}
	//
	auth.POST("/register", r.c.Register)
	//
	//
	//http://localhost:4321/auth/register
	// {"username":"Peter","password":"123123"}
	//
	//
	auth.POST("/login", r.c.Login)
	//
	//
	//
	router.GET("/app", r.c.App)
	//
	// http://localhost:4321/ws/create
	//	{  "id":1,"name":"osu" }
	//	{  "id":2, "name":"heroes of might and magic" }
	//	{  "id":3, "name":"apex legends" }
	//
	router.POST("/ws/create", r.c.CreateRoom)
	//
	// ws://localhost:4321/ws/join/2?clientid=2&username=peter
	//
	router.GET("/ws/join/:roomid", r.c.JoinRoom)
	//
	//http://168.34.32.22:4321/ws/rooms
	//
	router.GET("/ws/rooms", r.c.GetRooms)
	//
	// http://localhost:4321/ws/clietns/1
	//
	router.GET("/ws/clients/:roomid", r.c.GetClients)

}
