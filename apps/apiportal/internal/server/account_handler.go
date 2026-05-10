package server

import (
	"net/http"

	db "github.com/dewasurya/kakeiboku/apps/apiportal/internal/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/utils"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Balance  float64 `json:"balance" binding:"required,min=0"`
	Currency string  `json:"currency" binding:"required,oneof=USD EUR JPY"`
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

	_, err := s.Store.CreateAccounts(ctx, db.CreateAccountsParams{
		UserID:   1,
		Balance:  utils.IntToPgTypeNumeric(0),
		Currency: "",
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.CommonResponse("success create account"))

}
