package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadiwa/nginx-ui/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AppSecHandler struct {
	db *gorm.DB
}

func NewAppSecHandler(db *gorm.DB) *AppSecHandler {
	return &AppSecHandler{db: db}
}

// GetConfig returns the current AppSec configuration
func (h *AppSecHandler) GetConfig(c *gin.Context) {
	var config model.AppSecConfig
	result := h.db.First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			config = model.AppSecConfig{
				Enabled: false,
			}
			h.db.Create(&config)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, config)
}

// UpdateConfig updates the AppSec configuration
func (h *AppSecHandler) UpdateConfig(c *gin.Context) {
	var config model.AppSecConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&model.AppSecConfig{}).Where("id = ?", config.ID).Updates(config)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetPolicies returns all AppSec policies
func (h *AppSecHandler) GetPolicies(c *gin.Context) {
	var policies []model.AppSecPolicy
	result := h.db.Find(&policies)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, policies)
}

// CreatePolicy creates a new AppSec policy
func (h *AppSecHandler) CreatePolicy(c *gin.Context) {
	var policy model.AppSecPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	result := h.db.Create(&policy)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, policy)
}

// UpdatePolicy updates an existing AppSec policy
func (h *AppSecHandler) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	var policy model.AppSecPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy.UpdatedAt = time.Now()

	result := h.db.Model(&model.AppSecPolicy{}).Where("id = ?", id).Updates(policy)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// DeletePolicy deletes an AppSec policy
func (h *AppSecHandler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	result := h.db.Delete(&model.AppSecPolicy{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Policy deleted successfully"})
}

// GetStats returns AppSec statistics
func (h *AppSecHandler) GetStats(c *gin.Context) {
	var stats model.AppSecStats

	// Get total requests
	h.db.Model(&model.AppSecEvent{}).Count(&stats.TotalRequests)

	// Get blocked requests
	h.db.Model(&model.AppSecEvent{}).Where("action = ?", "block").Count(&stats.BlockedRequests)

	// Get alerted requests
	h.db.Model(&model.AppSecEvent{}).Where("action = ?", "alert").Count(&stats.AlertedRequests)

	// Get top attacks
	var topAttacks []model.TopAttack
	h.db.Raw(`
		SELECT rule_id as type, COUNT(*) as count 
		FROM app_sec_events 
		WHERE created_at >= ? 
		GROUP BY rule_id 
		ORDER BY count DESC 
		LIMIT 10
	`, time.Now().AddDate(0, 0, -7)).Scan(&topAttacks)

	stats.TopAttacks = topAttacks
	c.JSON(http.StatusOK, stats)
}
