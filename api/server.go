package api

import (
	"fmt"

	db "github.com/GiorgiMakharadze/bank-API-golang/db/sqlc"
	"github.com/GiorgiMakharadze/bank-API-golang/token"
	"github.com/GiorgiMakharadze/bank-API-golang/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/api/v1/users", server.createUser)

	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.GET("/api/v1/accounts", server.listAccount)

	router.POST("/api/v1/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
