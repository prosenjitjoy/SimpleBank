package api

import (
	"fmt"
	"main/database/db"
	"main/token"
	"main/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for our banking service.
type Server struct {
	config     *util.ConfigDatabase
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store, cfg *util.ConfigDatabase) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.POST("/tokens/renew", s.renewAccessToken)

	authRoutes := router.Group("/")
	authRoutes.Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAcount)
	authRoutes.GET("/accounts", s.listAcount)
	authRoutes.PATCH("/accounts/:id", s.updateAccount)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
