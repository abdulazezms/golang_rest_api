package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alpha"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		FullName:  user.FullName,
	}
}

func (s Server) createUser(ctx *gin.Context) {
	var requestBody createUserRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		handleErrorBinding(ctx, err)
		return
	}

	hashedPassword, err := util.HashPassword(requestBody.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username:       requestBody.Username,
		FullName:       requestBody.FullName,
		Email:          requestBody.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, newUserResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (s Server) loginUser(ctx *gin.Context) {
	var requestBody loginUserRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		handleErrorBinding(ctx, err)
		return
	}

	user, err := s.store.GetUser(ctx, requestBody.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	err = util.CheckPassword(requestBody.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	t, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenValidation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserResponse{
		Token: t,
		User:  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, res)
}
