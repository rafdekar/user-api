package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/rafdekar/user-api/db/sqlc"
	"net/http"
)

// Server serves all HTTP requests for banking service
type Server struct {
	queries db.Querier
	router  *gin.Engine
}

// NewServer starts a new server
func NewServer(db db.Querier) *Server {
	server := &Server{
		queries: db,
	}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users", server.listUsers)
	router.PUT("/users", server.updateUser)
	router.DELETE("/users", server.deleteUser)

	router.HEAD("/_health", server.health)

	server.router = router
	return server
}

// Start method starts gin on specified address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// health is a health check endpoint
func (s *Server) health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "PONG")
}

// errorResponse is function for formatting error responses to be returned by gin handler
func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
