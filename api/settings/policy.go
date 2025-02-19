package settings

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const PolicyPath = "/etc/cp/conf/local_policy.yaml"

// InitPolicyRouter initializes the policy routes
func InitPolicyRouter(r *gin.RouterGroup) {
	r.GET("/policy", GetPolicy)
	r.POST("/policy", SavePolicy)
}

// GetPolicy handles retrieving the policy file content
func GetPolicy(c *gin.Context) {
	content, err := ioutil.ReadFile(PolicyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, create it with default content
			defaultContent := "# Default policy configuration\n"
			err = ioutil.WriteFile(PolicyPath, []byte(defaultContent), 0644)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to create policy file: %v", err),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"content": defaultContent,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to read policy file: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": string(content),
	})
}

// SavePolicy handles saving new content to the policy file
func SavePolicy(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err := ioutil.WriteFile(PolicyPath, []byte(req.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to save policy file: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Policy saved successfully",
	})
}
