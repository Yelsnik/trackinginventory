package api

import (
	"fmt"

	db "github.com/Yelsnik/trackinginventory/db/sqlc"
	"github.com/Yelsnik/trackinginventory/token"
	"github.com/Yelsnik/trackinginventory/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db         db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/inventory", server.createInventory)
	authRoutes.GET("/inventory/:id", server.getInventory)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

func NewServer(config util.Config, db db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		db:         db,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
