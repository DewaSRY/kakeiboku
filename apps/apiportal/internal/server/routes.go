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
	v1_routes.GET("/health", server.HealthHandler)

	auth_routes := v1_routes.Group("/auth")
	auth_routes.POST("/login", server.LoginHandler)
	auth_routes.POST("/signup", server.SignUpHandler)
	auth_routes.POST("/refresh-token", server.RefreshTokenHandler)

	user_routes := v1_routes.Group("/user")
	user_routes.Use(middleware.AuthMiddleware(server.Token))
	user_routes.GET("/me", server.signInHandler)

	account_routes := v1_routes.Group("/account")
	account_routes.Use(middleware.AuthMiddleware(server.Token))
	account_routes.POST("/", server.CreateAccountHandler)
	account_routes.GET("/", server.GetAccountHandler)

	transaction_routes := v1_routes.Group("/transaction")
	transaction_routes.Use(middleware.AuthMiddleware(server.Token))
	transaction_routes.POST("/", server.TransactionHandler)
	
	return router
}
