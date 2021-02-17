package routes

import (
	"testing"
	"github.com/gin-gonic/gin"
)


func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupRouter()
}