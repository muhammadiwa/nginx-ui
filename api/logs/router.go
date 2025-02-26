package logs

import "github.com/gin-gonic/gin"

func InitRouter(g *gin.RouterGroup) {
	// Remove the /logs prefix since it's handled by the group
	g.GET("", GetLogs)
}
