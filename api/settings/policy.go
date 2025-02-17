package settings

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const policyPath = "/etc/cp/conf/local_policy.yaml"

func GetLocalPolicy(c *gin.Context) {
	// Baca file local_policy.yaml
	data, err := ioutil.ReadFile(policyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read policy file",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": string(data),
		"path": policyPath,
	})
}

func SaveLocalPolicy(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Hapus spasi di awal baris jika ada
	content := strings.TrimSpace(req.Content)

	// Tulis ke file
	err := ioutil.WriteFile(policyPath, []byte(content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save policy file",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Policy saved successfully",
	})
}
