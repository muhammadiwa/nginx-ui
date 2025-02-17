package preference

import (
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetPolicy(c *gin.Context) {
	content, err := os.ReadFile("/etc/cp/conf/local_policy.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

func SavePolicy(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := exec.Command("sudo", "tee", "/etc/cp/conf/local_policy.yaml")
	cmd.Stdin = strings.NewReader(req.Content)
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Policy saved successfully"})
}

func InitRouter(r *gin.RouterGroup) {
	r.GET("/policy", GetPolicy)
	r.POST("/policy", SavePolicy)
}
