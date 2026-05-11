package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) HealthHandler(c *gin.Context) {
	stats := server.Store.Health(c.Request.Context())

	statusCode := http.StatusOK
	
	if stats.Status == "down" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, stats)
}
