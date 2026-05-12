package server

import (
	"net/http"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/middleware"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	db "github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Balance  float64 `json:"balance" binding:"required,min=0"`
	Currency string  `json:"currency" binding:"required,currency"`
}

type CreateAccountResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (s *Server) CreateAccountHandler(ctx *gin.Context) {
	var req CreateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	auth_payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	_, err := s.Store.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   auth_payload.UserID,
		Balance:  utils.IntToPgTypeNumeric(0),
		Currency: req.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.CommonResponse("success create account"))
}

type GetAccountQuery struct {
	utils.PaginationQueryParam
}

type GetAccountResponse struct {
	Data     []services.Account       `json:"data"`
	Metadata utils.PaginationMetaData `json:"metadata"`
}

func (server *Server) GetAccountHandler(ctx *gin.Context) {
	var req GetAccountQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	auth_payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	accounts, err := server.Store.ListUserAccounts(ctx, services.ListUserAccountsParams{
		UserID: auth_payload.UserID,
		Limit:  req.Limit,
		Offset: req.GetOffset(),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	account_amount, err := server.Store.GetUserAccountCount(ctx, auth_payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	response := GetAccountResponse{
		Data:     accounts,
		Metadata: req.GetMetadata(int32(account_amount)),
	}

	ctx.JSON(http.StatusOK, response)
}
