package routes

import (
	"github.com/devops-kung-fu/domi/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRouter - Set up gin router
func SetupRouter() *gin.Engine { 
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	misc := r.Group("/")
	{
		misc.GET("alive", controllers.CanYouHearMeNow)
	}
	github := r.Group("/github/v1")
	{
		github.POST("webhook", controllers.ReceiveGitHubWebHook)
	}
	return r
}