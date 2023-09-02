package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/token"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}
type ApiError struct {
	Field string
	Msg   string
}

func (s Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleErrorBinding(ctx, err)
		return
	}
	//retrieve the account with the FromID given
	valid, accountFrom := s.isValidAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	valid, _ = s.isValidAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	//make sure that the account actually belongs to the user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if accountFrom.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("From Account ID cannot be found in your list of accounts")))
		return
	}

	result, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, result)
}
