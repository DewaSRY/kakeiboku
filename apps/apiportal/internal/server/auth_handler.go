package server

import (
	"net/http"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) SignUpHandler(ctx *gin.Context) {
	var req SignupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// create user
	hash_password, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user_created, err := server.store.CreateUser(ctx, services.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash_password,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create jwt token
	access_token, access_payload, err := server.token.CreateToken(
		user_created.ID,
		user_created.Email,
		server.config.AccessTokenDuration,
		token.TokenTypeAccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refresh_token, refresh_payload, err := server.token.CreateToken(
		user_created.ID,
		user_created.Email,
		server.config.RefreshTokenDuration,
		token.TokenTypeRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, services.CreateSessionParams{
		ID:           access_payload.ID,
		Email:        user_created.Email,
		RefreshToken: refresh_token,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_payload.ExpiredAt,
	})
}


