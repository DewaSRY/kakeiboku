package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	router.GET("/", server.HelloWorldHandler)
	router.GET("/health", server.healthHandler)

	account_routes := router.Group("/accounts") 

	account_routes.POST("/", server.CreateAccountHandler)
	return router
}

func (server *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (server *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, server.db.Health())
}
