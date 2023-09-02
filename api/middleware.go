package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"tutorial.sqlc.dev/app/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

// higher order functions that will return the authentication middleware function
func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeaderValue := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeaderValue) == 0 {
			err := errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Split(authorizationHeaderValue, " ")
		if len(fields) != 2 {
			err := errors.New("Authorization header is invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if fields[0] != authorizationTypeBearer {
			err := fmt.Errorf("Unsuppported Authorization type %s", fields[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
