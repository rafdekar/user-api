package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/rafdekar/user-api/db/sqlc"
	"net/http"
)

type createUserRequest struct {
	FirstName string `json:"first_name" binding:"alpha"`
	LastName  string `json:"last_name" binding:"alpha"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email" binding:"email"`
	Country   string `json:"country" binding:"len=2,alpha"`
}

// createUser method defines endpoint for createing a user
func (s *Server) createUser(ctx *gin.Context) {
	request := &createUserRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateUserParams{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Nickname:  request.Nickname,
		Password:  request.Password,
		Email:     request.Email,
		Country:   request.Country,
	}

	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type updateUserRequest struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	FirstName string    `json:"first_name" binding:"alpha"`
	LastName  string    `json:"last_name" binding:"alpha"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email" binding:"email"`
	Country   string    `json:"country" binding:"len=2,alpha"`
}

// updateUser method defines endpoint for updating selected user data
func (s *Server) updateUser(ctx *gin.Context) {
	request := &updateUserRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.UpdateUserParams{
		ID:        request.ID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Nickname:  request.Nickname,
		Password:  request.Password,
		Email:     request.Email,
		Country:   request.Country,
	}

	user, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUsersRequest struct {
	PageSize   int32 `json:"page_size" binding:"required,min=1"`
	PageNumber int32 `json:"page_number" binding:"required,min=1"`
}

// listUsers method defines endpoint for listing users from page X of size Y
func (s *Server) listUsers(ctx *gin.Context) {
	request := &listUsersRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListUsersParams{
		Offset: (request.PageNumber - 1) * request.PageSize,
		Limit:  request.PageSize,
	}

	users, err := s.queries.ListUsers(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type deleteUserRequest struct {
	ID uuid.UUID `json:"ID" binding:"required"`
}

// deleteUser defines endpoint for completely deleting user data
func (s *Server) deleteUser(ctx *gin.Context) {
	request := &deleteUserRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.queries.DeleteUser(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
