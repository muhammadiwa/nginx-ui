package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func formatEventTime(eventTime string) string {
	t, _ := time.Parse("2006-01-02T15:04:05.999", eventTime)
	jakartaLoc, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(jakartaLoc).Format("January 2 2006 15:04")
}

func GetLogs(c *gin.Context) {
	// Get filter parameters
	eventSeverity := c.Query("severity")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	cmd := exec.Command("sudo", "open-appsec-ctl", "-vl")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Failed to execute command: %s", err.Error()),
		})
		return
	}

	lines := strings.Split(out.String(), "\n")
	var events []map[string]interface{}

	// Parse start and end dates if provided
	var startTime, endTime time.Time
	var hasDateFilter bool
	if startDate != "" && endDate != "" {
		startTime, err = time.Parse("2006-01-02", startDate)
		if err == nil {
			endTime, err = time.Parse("2006-01-02", endDate)
			if err == nil {
				endTime = endTime.Add(24 * time.Hour) // Include the entire end date
				hasDateFilter = true
			}
		}
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		// Apply severity filter if specified
		if eventSeverity != "" {
			if severity, ok := event["eventSeverity"].(string); ok {
				if severity != eventSeverity {
					continue
				}
			}
		}

		// Format and filter by date range if specified
		if eventTimeStr, ok := event["eventTime"].(string); ok {
			eventTime, err := time.Parse("2006-01-02T15:04:05.999", eventTimeStr)
			if err == nil && hasDateFilter {
				if eventTime.Before(startTime) || eventTime.After(endTime) {
					continue
				}
			}
			formattedTime := formatEventTime(eventTimeStr)
			event["eventTime"] = formattedTime
		}

		events = append(events, event)
	}

	// Sort events by time descending (newest first)
	sort.Slice(events, func(i, j int) bool {
		timeI, _ := time.Parse("January 2 2006 15:04", events[i]["eventTime"].(string))
		timeJ, _ := time.Parse("January 2 2006 15:04", events[j]["eventTime"].(string))
		return timeI.After(timeJ)
	})

	c.JSON(200, gin.H{
		"data":    events,
		"message": "Logs retrieved successfully",
	})
}
