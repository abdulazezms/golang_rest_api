package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	//FieldLevel is an interface that contains all info and helper functions to validate a field.
	if currency, ok := fieldLevel.Field().Interface().(string); ok { //get the value of a field.
		return util.IsSupportedCurrency(currency)
	}
	//the field isn't a string.
	return false
}

func (s Server) isValidAccount(ctx *gin.Context, accountID int64, currency string) (bool, db.Account) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return false, account
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false, account
	}
	return true, account
}
