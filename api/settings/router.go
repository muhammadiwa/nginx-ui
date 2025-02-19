package settings

import (
	"github.com/0xJacky/Nginx-UI/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	// Settings routes
	r.GET("server/name", GetServerName)
	r.GET("", GetSettings)
	r.POST("", middleware.RequireSecureSession(), SaveSettings)

	r.GET("auth/banned_ips", GetBanLoginIP)
	r.DELETE("auth/banned_ip", RemoveBannedIP)

	// Policy routes - tanpa auth middleware untuk testing
	r.GET("policy", GetPolicy)
	r.POST("policy", SavePolicy)
}
