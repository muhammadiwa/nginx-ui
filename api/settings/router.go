package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/0xJacky/Nginx-UI/internal/middleware"
)

func InitRouter(r *gin.RouterGroup) {
	// Settings routes
	r.GET("server/name", GetServerName)
	r.GET("", GetSettings)
	r.POST("", middleware.RequireSecureSession(), SaveSettings)

	r.GET("auth/banned_ips", GetBanLoginIP)
	r.DELETE("auth/banned_ip", RemoveBannedIP)
}
