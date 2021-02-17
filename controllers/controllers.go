package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ReceiveGitHubWebHook - Receives and processes GitHub WebHook Events
func ReceiveGitHubWebHook(c *gin.Context) {
	fmt.Println(c.Request.Body)
}