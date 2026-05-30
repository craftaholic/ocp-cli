package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Active string `json:"active"`
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "ocp"), nil
}

func GetProfilesDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "profiles"), nil
}

func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	profilesDir := filepath.Join(configDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

func IsFirstRun() (bool, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return false, err
	}
	
	configPath := filepath.Join(configDir, "config.json")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func MigrateExistingConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	opencodeDir := filepath.Join(homeDir, ".config", "opencode")
	
	info, err := os.Lstat(opencodeDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check opencode directory: %w", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return nil
	}

	if !info.IsDir() {
		return nil
	}

	fmt.Println("🔍 Detected existing opencode configuration at ~/.config/opencode")
	fmt.Println("📦 Migrating to ocp default profile...")

	if err := EnsureConfigDir(); err != nil {
		return err
	}

	profilesDir, err := GetProfilesDir()
	if err != nil {
		return err
	}

	defaultProfileDir := filepath.Join(profilesDir, "default")
	if _, err := os.Stat(defaultProfileDir); err == nil {
		return fmt.Errorf("default profile already exists")
	}

	backupDir := opencodeDir + ".backup"
	if err := os.Rename(opencodeDir, backupDir); err != nil {
		fmt.Printf("⚠️  Cannot migrate ~/.config/opencode (permission denied or busy)\n")
		fmt.Printf("   Skipping migration. You can set up profiles manually.\n\n")
		return nil
	}

	if err := os.Rename(backupDir, defaultProfileDir); err != nil {
		os.Rename(backupDir, opencodeDir)
		return fmt.Errorf("failed to move config to default profile: %w", err)
	}

	profile := &Profile{
		Name: "default",
		Vars: make(map[string]string),
	}

	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		profile.Vars["ANTHROPIC_API_KEY"] = apiKey
	}

	if err := SaveProfile(profile); err != nil {
		return fmt.Errorf("failed to create default profile: %w", err)
	}

	cfg := &Config{
		Active: "default",
	}
	if err := SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := UpdateSymlink("default"); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Println("✅ Migration complete!")
	fmt.Println("   - Moved ~/.config/opencode to ~/.config/ocp/profiles/default/")
	fmt.Println("   - Created symlink: ~/.config/opencode -> ~/.config/ocp/profiles/default/")
	fmt.Println("   - Set 'default' as active profile")
	fmt.Println()
	fmt.Println("Your existing configuration is now in the 'default' profile.")
	fmt.Println("You can create additional profiles with: ocp add <profile>")
	fmt.Println()

	return nil
}

func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	
	configPath := filepath.Join(configDir, "config.json")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		if err := MigrateExistingConfig(); err != nil {
			return nil, err
		}
		data, err = os.ReadFile(configPath)
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	tmpPath := configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
