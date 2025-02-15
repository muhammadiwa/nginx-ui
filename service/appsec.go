package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/muhammadiwa/nginx-ui/model"
)

type AppSecService struct {
	configDir string
	modsecDir string
	mu        sync.Mutex
}

func NewAppSecService(nginxConfigDir string) *AppSecService {
	return &AppSecService{
		configDir: filepath.Join(nginxConfigDir, "conf.d", "appsec"),
		modsecDir: filepath.Join(nginxConfigDir, "modsec"),
	}
}

func (s *AppSecService) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create necessary directories
	dirs := []string{s.configDir, s.modsecDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Initialize ModSecurity main configuration
	mainConfPath := filepath.Join(s.modsecDir, "main.conf")
	if _, err := os.Stat(mainConfPath); os.IsNotExist(err) {
		defaultConfig := `
# Basic ModSecurity configuration
SecRuleEngine On
SecRequestBodyAccess On
SecResponseBodyAccess On
SecResponseBodyMimeType text/plain text/html text/xml application/json
SecDataDir /tmp
`
		if err := os.WriteFile(mainConfPath, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("failed to create ModSecurity main config: %v", err)
		}
	}

	// Initialize rules directory
	rulesDir := filepath.Join(s.modsecDir, "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create rules directory: %v", err)
	}

	return nil
}

func (s *AppSecService) UpdatePolicy(policy model.AppSecPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules := []string{}
	for _, rule := range policy.Rules {
		rules = append(rules, s.convertToModSecurityRule(rule))
	}

	policyPath := filepath.Join(s.modsecDir, "rules", fmt.Sprintf("policy_%s.conf", policy.ID))
	content := strings.Join(rules, "\n\n")

	if err := os.WriteFile(policyPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write policy file: %v", err)
	}

	return nil
}

func (s *AppSecService) convertToModSecurityRule(ruleID string) string {
	// This is a simplified example. You should implement proper rule conversion logic
	return fmt.Sprintf(`
SecRule REQUEST_URI "@rx .*" \
    "id:%s,\
    phase:2,\
    deny,\
    status:403,\
    log,\
    msg:'Custom Rule Violation'"
`, ruleID)
}

func (s *AppSecService) DeletePolicy(policyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	policyPath := filepath.Join(s.modsecDir, "rules", fmt.Sprintf("policy_%s.conf", policyID))
	if err := os.Remove(policyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete policy file: %v", err)
	}

	return nil
}

func (s *AppSecService) GetStats() (*model.AppSecStats, error) {
	// This is a placeholder. You should implement proper log parsing
	// from ModSecurity audit logs to get real statistics
	return &model.AppSecStats{
		TotalRequests:   1000,
		BlockedRequests: 50,
		AlertedRequests: 20,
		TopAttacks: []model.TopAttack{
			{Type: "XSS", Count: 30},
			{Type: "SQL Injection", Count: 15},
			{Type: "File Upload", Count: 5},
		},
	}, nil
}
