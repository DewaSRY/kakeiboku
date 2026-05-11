package server

import (
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/middleware"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/gin-gonic/gin"
)



func (server *Server) signInHandler	(ctx *gin.Context) {		

	auth_payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	user, err := server.Store.GetUserByID(ctx, auth_payload.UserID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, user)
}
