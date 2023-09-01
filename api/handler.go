package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func handleErrorBinding(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	}
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "currency":
		return "Invalid currency!"
	case "alpha":
		return "must contain only alpha characters"
	case "email":
		return "must be valid email"
	}
	return ""
}
