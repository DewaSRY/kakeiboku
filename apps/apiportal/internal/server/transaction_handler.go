package server

import (
	"net/http"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CreateTransactionRequest struct {
	FromAccount int64 `json:"from_account" binding:"required"`
	ToAccount   int64 `json:"to_account" binding:"required"`
	Amount      int   `json:"amount" binding:"required,gt=0"`
}

func (server *Server) TransactionHandler(ctx *gin.Context) {
	var req CreateTransactionRequest 

	if err := utils.BindJSON(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// auth_payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	_, err := server.Store.CreateTransferTx(ctx, services.CreateTransactionParams{
		FromAccountID: req.FromAccount,
		ToAccountID:   req.ToAccount,
		Amount:        utils.IntToPgTypeNumeric(req.Amount),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, utils.CommonResponse("success create transaction"))
}
