package model

import (
	"encoding/json"
	"time"
)

type AppSecConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppSecPolicy struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Rules       []string  `json:"rules" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppSecRule struct {
	ID          string           `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Action      string           `json:"action"`
	Conditions  json.RawMessage  `json:"conditions" gorm:"type:json"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type AppSecStats struct {
	TotalRequests   int64       `json:"totalRequests"`
	BlockedRequests int64       `json:"blockedRequests"`
	AlertedRequests int64       `json:"alertedRequests"`
	TopAttacks      []TopAttack `json:"topAttacks"`
}

type TopAttack struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type AppSecEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RequestID   string    `json:"requestId"`
	PolicyID    string    `json:"policyId"`
	RuleID      string    `json:"ruleId"`
	Action      string    `json:"action"`
	RequestData string    `json:"requestData" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}
