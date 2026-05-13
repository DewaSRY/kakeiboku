package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) SignUpHandler(ctx *gin.Context) {
	var req SignupRequest

	if err := utils.BindJSON(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// create user
	hash_password, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	user_created, err := server.Store.CreateUserTx(ctx, services.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash_password,
	})

	if err != nil {
		if errors.Is(err, services.ErrCreatingUserWIthDuplicateEmail) {
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(fmt.Errorf("email %s already exists", req.Email)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	// create jwt token
	access_token, access_payload, err := server.Token.CreateToken(
		user_created.ID,
		user_created.Email,
		server.Config.AccessTokenDuration,
		token.TokenTypeAccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	refresh_token, refresh_payload, err := server.Token.CreateToken(
		user_created.ID,
		user_created.Email,
		server.Config.RefreshTokenDuration,
		token.TokenTypeRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	_, err = server.Store.SetSession(ctx, services.CreateSessionParams{
		ID:           access_payload.ID,
		Email:        user_created.Email,
		RefreshToken: refresh_token,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	server.setCookies(ctx, utils.KeyAccessToken, access_token, int(server.Config.AccessTokenDuration.Seconds()))
	server.setCookies(ctx, utils.KeyRefreshToken, refresh_token, int(server.Config.RefreshTokenDuration.Seconds()))

	ctx.JSON(http.StatusOK, AuthResponse{
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_payload.ExpiredAt,
	})
}

func (server *Server) LoginHandler(ctx *gin.Context) {
	var req LoginRequest

	fmt.Printf("hit this end point")

	if err := utils.BindJSON(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := server.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	// create jwt token
	access_token, access_payload, err := server.Token.CreateToken(
		user.ID,
		user.Email,
		server.Config.AccessTokenDuration,
		token.TokenTypeAccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	refresh_token, refresh_payload, err := server.Token.CreateToken(
		user.ID,
		user.Email,
		server.Config.RefreshTokenDuration,
		token.TokenTypeRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	_, err = server.Store.SetSession(ctx, services.CreateSessionParams{
		ID:           access_payload.ID,
		Email:        user.Email,
		RefreshToken: refresh_token,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	server.setCookies(ctx, utils.KeyAccessToken, access_token, int(server.Config.AccessTokenDuration.Seconds()))
	server.setCookies(ctx, utils.KeyRefreshToken, refresh_token, int(server.Config.RefreshTokenDuration.Seconds()))

	ctx.JSON(http.StatusOK, AuthResponse{
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_payload.ExpiredAt,
	})
}

func (server *Server) RefreshTokenHandler(ctx *gin.Context) {
	var req RefreshTokenRequest

	if err := utils.BindJSON(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	refresh_token := req.RefreshToken
	if(req.RefreshToken == "" ) {
		cookie, err := ctx.Cookie(utils.KeyRefreshToken)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("refresh token is required")))
			return
		}
		refresh_token = cookie
	}

	 if refresh_token == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("refresh token is required")))
		return
	}

	refresh_token_payload, err := server.Token.VerifyToken(refresh_token, token.TokenTypeRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	user, err := server.Store.GetUserByEmail(ctx, refresh_token_payload.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	access_token, access_payload, err := server.Token.CreateToken(
		user.ID,
		user.Email,
		server.Config.AccessTokenDuration,
		token.TokenTypeAccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	server.setCookies(ctx, utils.KeyAccessToken, access_token, int(server.Config.AccessTokenDuration.Seconds()))

	ctx.JSON(http.StatusOK, AuthResponse{
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_token_payload.ExpiredAt,
	})
}


