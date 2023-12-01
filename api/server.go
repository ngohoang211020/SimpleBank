package api

import (
	"github.com/gin-gonic/gin"
	db "simplebank/db/sqlc"
)

// Server servers HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

func (server *Server) Start(port string) error {
	return server.router.Run(":" + port)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
