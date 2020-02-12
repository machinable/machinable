package hooks

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	code := m.Run()
	os.Exit(code)
}

func TestUpdateHook(t *testing.T) {

}

func TestListHooks(t *testing.T) {

}

func TestAddHook(t *testing.T) {

}

func TestGetHook(t *testing.T) {

}

func TestDeleteHook(t *testing.T) {

}
