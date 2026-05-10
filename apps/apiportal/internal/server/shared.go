package server

import "github.com/gin-gonic/gin"


func (server *Server) setCookies(ctx *gin.Context, key string, value string, maxAge int) {
	ctx.SetCookie(
		key, 
		value, 
		maxAge, 
		"/", 
		server.Config.AppDomain, server.Config.AppEnv == "production", 
		true,
	)
}
