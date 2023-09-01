package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "tutorial.sqlc.dev/app/db/sqlc"
)

// Server serves HTTP requests coming
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//returns the underlying validator engine which powers the StructValidator implementation.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	grp := router.Group("/v1")
	grp.POST("/accounts", server.createAccount)
	grp.GET("/accounts/:id", server.getAccount)
	grp.GET("/accounts", server.listAccounts)
	grp.POST("/transfers", server.createTransfer)
	grp.POST("/users", server.createUser)
	server.router = router
	return server
}

// Start runs the server and makes it listen to a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
