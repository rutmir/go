package common

import (
	"fmt"
	"os"
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
		panic(fmt.Sprintf("Environment parameter required %s", propName))
	}
	return result
}

// UniqueStrSlice returns a unique subset of the string slice provided.
func UniqueStrSlice(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
