package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //an implementation of a driver
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/util"
)

func NewTestServer(t *testing.T, store db.Store) (*Server, error) {
	config := util.Config{
		TokenSymmetricKey:     util.RandomString(32),
		AccessTokenValidation: time.Minute,
	}
	server, err := NewServer(config, store)
	return server, err
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
