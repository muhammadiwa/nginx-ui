package policy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

const PolicyPath = "/etc/cp/conf/local_policy.yaml"

// GetPolicy returns the content of policy file
func GetPolicy(c *gin.Context) {
	// Check if file exists, if not create empty
	if _, err := os.Stat(PolicyPath); os.IsNotExist(err) {
		err = ioutil.WriteFile(PolicyPath, []byte(""), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create policy file",
			})
			return
		}
	}

	// Read policy file
	content, err := ioutil.ReadFile(PolicyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read policy file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": string(content),
	})
}

// SavePolicy saves the policy content to file
func SavePolicy(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Write to policy file
	err := ioutil.WriteFile(PolicyPath, []byte(req.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write policy file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Policy saved successfully",
	})
}

// ApplyPolicy applies the policy configuration using open-appsec-ctl command
func ApplyPolicy(c *gin.Context) {
	// Jalankan command dengan sudo
	cmd := exec.Command("sudo", "open-appsec-ctl", "-ap")

	// Ambil output command
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to apply policy: %v", err),
			"details": outputStr,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Policy applied successfully",
		"details": outputStr,
	})
}
