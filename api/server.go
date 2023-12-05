package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
)

// Server servers HTTP requests for our banking service
type Server struct {
	store       db.Store
	tokenMarker token.Maker
	router      *gin.Engine
	config      *util.Configuration
}

func NewServer(config *util.Configuration, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:       store,
		tokenMarker: tokenMaker,
		config:      config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(port string) error {
	return server.router.Run(":" + port)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
