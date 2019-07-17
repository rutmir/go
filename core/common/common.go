package common

import (
	"os"

	"github.com/rutmir/go/core/logger"
)

// GetPropertyOrDefault ...
func GetPropertyOrDefault(propName, defaultValue string) string {
	var result string
	if result = os.Getenv(propName); len(result) > 0 {
		return result
	}
	return defaultValue
}

// GetRequiredProperty ...
func GetRequiredProperty(propName string) string {
	var result string
	if result = os.Getenv(propName); len(result) == 0 {
		logger.Fatal("Environment parameter required", propName)
	}
	return result
}
