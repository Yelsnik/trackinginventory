package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Yelsnik/trackinginventory/db/sqlc"
	"github.com/Yelsnik/trackinginventory/util"
	"github.com/gin-gonic/gin"
	// "github.com/lib/pq"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID                int64     `json:"id"`
	Name              string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.db.CreateUser(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the user by email from the db
	user, err := server.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if the user password matches the req password
	err = util.ComparePassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// create a token
	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// send it to the user
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	fmt.Println(user.ID)

	ctx.JSON(http.StatusOK, rsp)
}

//
