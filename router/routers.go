package router

import (
	"net/http"

	"github.com/0xJacky/Nginx-UI/api/analytic"
	"github.com/0xJacky/Nginx-UI/api/certificate"
	"github.com/0xJacky/Nginx-UI/api/cluster"
	"github.com/0xJacky/Nginx-UI/api/config"
	"github.com/0xJacky/Nginx-UI/api/nginx"
	nginxLog "github.com/0xJacky/Nginx-UI/api/nginx_log"
	"github.com/0xJacky/Nginx-UI/api/notification"
	"github.com/0xJacky/Nginx-UI/api/openai"
	"github.com/0xJacky/Nginx-UI/api/public"
	"github.com/0xJacky/Nginx-UI/api/settings"
	"github.com/0xJacky/Nginx-UI/api/sites"
	"github.com/0xJacky/Nginx-UI/api/streams"
	"github.com/0xJacky/Nginx-UI/api/system"
	"github.com/0xJacky/Nginx-UI/api/template"
	"github.com/0xJacky/Nginx-UI/api/terminal"
	"github.com/0xJacky/Nginx-UI/api/upstream"
	"github.com/0xJacky/Nginx-UI/api/user"
	"github.com/0xJacky/Nginx-UI/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

// InitRouter initializes the router
func InitRouter() {
	r := cosy.GetEngine()

	initEmbedRoute(r)
	initUiRoute(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	})

	root := r.Group("/api")
	{
		public.InitRouter(root)
		system.InitPublicRouter(root)
		user.InitAuthRouter(root)

		// Authorization required and not websocket request
		g := root.Group("/", middleware.AuthRequired(), middleware.Proxy())
		{
			user.InitUserRouter(g)
			analytic.InitRouter(g)
			user.InitManageUserRouter(g)
			nginx.InitRouter(g)
			sites.InitCategoryRouter(g)
			sites.InitRouter(g)
			streams.InitRouter(g)
			config.InitRouter(g)
			template.InitRouter(g)
			certificate.InitCertificateRouter(g)
			certificate.InitDNSCredentialRouter(g)
			certificate.InitAcmeUserRouter(g)
			system.InitPrivateRouter(g)
			settings.InitRouter(g)
			openai.InitRouter(g)
			cluster.InitRouter(g)
			notification.InitRouter(g)
		}

		// Authorization required and websocket request
		w := root.Group("/", middleware.AuthRequired(), middleware.ProxyWs())
		{
			analytic.InitWebSocketRouter(w)
			certificate.InitCertificateWebSocketRouter(w)
			o := w.Group("", middleware.RequireSecureSession())
			{
				terminal.InitRouter(o)
			}
			nginxLog.InitRouter(w)
			upstream.InitRouter(w)
			system.InitWebSocketRouter(w)
		}

		// AppSec routes
		appsec := root.Group("/v1/appsec")
		{
			appsec.GET("/config", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec config",
				})
			})
			appsec.PUT("/config", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec config updated",
				})
			})
			appsec.GET("/policies", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec policies",
				})
			})
			appsec.POST("/policies", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec policy created",
				})
			})
			appsec.PUT("/policies/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec policy updated",
				})
			})
			appsec.DELETE("/policies/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec policy deleted",
				})
			})
			appsec.GET("/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "AppSec stats",
				})
			})
		}
	}
}
