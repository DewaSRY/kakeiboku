package server

import "github.com/gin-gonic/gin"

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}



func commonResponse(message string) gin.H {
	return gin.H{"message": message, "success": true}
}	
