package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// ReceiveGitHubWebHook - Receives and processes GitHub WebHook Events
func ReceiveGitHubWebHook(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}