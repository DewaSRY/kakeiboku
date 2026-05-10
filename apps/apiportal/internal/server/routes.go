package server

import (
	"net/http"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *Server) http.Handler {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	v1_routes := router.Group("/v1")
	v1_routes.Use(middleware.WebMiddleware)
	v1_routes.GET("/", server.HelloWorldHandler)
	// v1_routes.GET("/health", server.healthHandler)

	auth_routes := v1_routes.Group("/auth")
	auth_routes.POST("/login", server.LoginHandler)
	auth_routes.POST("/signup", server.SignUpHandler)

	user_routes := v1_routes.Group("/users")
	user_routes.Use(middleware.AuthMiddleware(server.Token))
	user_routes.GET("/me", server.signInHandler)

	account_routes := v1_routes.Group("/accounts")
	account_routes.Use(middleware.AuthMiddleware(server.Token))

	account_routes.POST("/", server.CreateAccountHandler)

	return router
}

func (server *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// func (server *Server) healthHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, server.store.)
// }
