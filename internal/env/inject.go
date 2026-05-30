package env

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/craftaholic/ocp/internal/config"
)

func InjectProfileVars(profile *config.Profile) []string {
	env := os.Environ()
	envMap := make(map[string]string)

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	for key, value := range profile.Vars {
		envMap[key] = ExpandPath(value)
	}

	var result []string
	for key, value := range envMap {
		result = append(result, key+"="+value)
	}

	return result
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

func IsSensitive(key string) bool {
	lowerKey := strings.ToLower(key)
	sensitiveWords := []string{"key", "secret", "token", "password"}
	for _, word := range sensitiveWords {
		if strings.Contains(lowerKey, word) {
			return true
		}
	}
	return false
}

func MaskValue(value string) string {
	if len(value) <= 8 {
		return "***"
	}
	return value[:8] + "..."
}
