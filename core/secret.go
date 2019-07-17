package secret

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rutmir/go/core/logger"
)

var (
	secretPath = "/etc/secrets"
)

// Loads a value from the specified secret file
func getValue(name string, panic bool) string {
	path := filepath.Join(secretPath, name)
	secret, err := ioutil.ReadFile(path)
	if err != nil {
		if panic {
			logger.Fatalf("Unable to read secret - %s", path)
		} else {
			logger.Errf("Unable to read secret -  %s", path)
		}
	}
	return string(secret)
}

// GetValue return secret
func GetValue(name string) string {
	return getValue(name, false)
}

// GetValueOrPanic return secret or run panic
func GetValueOrPanic(name string) string {
	return getValue(name, true)
}

// GetSecretFilePath return path to file secret
func GetSecretFilePath(name string) (string, bool) {
	path := filepath.Join(secretPath, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", false
	}
	return path, true
}

// GetSecretFilePathOrPanic ...
func GetSecretFilePathOrPanic(name string) string {
	path, ok := GetSecretFilePath(name)
	if !ok {
		logger.Fatalf("Unable to get secret file path - %v, name - %v", path, name)
	}
	return path
}

// Initialize ...
func Initialize(path string) {
	if len(path) > 0 {
		secretPath = path
	}
}
