package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	main()
}
