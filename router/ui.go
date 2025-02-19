package router

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func initUiRoute(r *gin.Engine) {
	// Add MIME type for .vue files
	r.Use(func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, ".vue") {
			c.Header("Content-Type", "text/javascript")
		}
	})

	// Serve the frontend static files
	r.Static("/assets", "./app/dist/assets")
	r.StaticFile("/favicon.ico", "./app/dist/favicon.ico")

	// Handle the main HTML file
	r.GET("/", func(c *gin.Context) {
		c.File("./app/dist/index.html")
	})

	// Handle all other routes for SPA
	r.NoRoute(func(c *gin.Context) {
		// Skip API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		// Serve index.html for all other routes
		ext := path.Ext(c.Request.URL.Path)
		if ext == "" {
			c.File("./app/dist/index.html")
		} else {
			c.Status(http.StatusNotFound)
		}
	})
}
