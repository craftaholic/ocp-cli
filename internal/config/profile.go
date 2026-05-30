package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Profile struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

func GetProfileDir(name string) (string, error) {
	profilesDir, err := GetProfilesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(profilesDir, name), nil
}

func LoadProfile(name string) (*Profile, error) {
	profileDir, err := GetProfileDir(name)
	if err != nil {
		return nil, err
	}

	profilePath := filepath.Join(profileDir, "profile.json")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile '%s' does not exist", name)
		}
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}

	return &profile, nil
}

func SaveProfile(profile *Profile) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	profileDir, err := GetProfileDir(profile.Name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	profilePath := filepath.Join(profileDir, "profile.json")
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	tmpPath := profilePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	if err := os.Rename(tmpPath, profilePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save profile: %w", err)
	}

	return nil
}

func DeleteProfile(name string) error {
	profileDir, err := GetProfileDir(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	if err := os.RemoveAll(profileDir); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

func ListProfiles() ([]string, error) {
	profilesDir, err := GetProfilesDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	var profiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			profilePath := filepath.Join(profilesDir, entry.Name(), "profile.json")
			if _, err := os.Stat(profilePath); err == nil {
				profiles = append(profiles, entry.Name())
			}
		}
	}

	sort.Strings(profiles)
	return profiles, nil
}

func ProfileExists(name string) (bool, error) {
	profileDir, err := GetProfileDir(name)
	if err != nil {
		return false, err
	}

	profilePath := filepath.Join(profileDir, "profile.json")
	_, err = os.Stat(profilePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check profile: %w", err)
	}

	return true, nil
}

func UpdateSymlink(profileName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	opencodeLink := filepath.Join(homeDir, ".config", "opencode")
	
	info, err := os.Lstat(opencodeLink)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(opencodeLink); err != nil {
				return fmt.Errorf("failed to remove old symlink: %w", err)
			}
		} else {
			fmt.Printf("⚠️  Warning: %s exists but is not a symlink. Skipping symlink creation.\n", opencodeLink)
			fmt.Printf("   Profile switching will work, but opencode will not use the profile directory.\n")
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check symlink: %w", err)
	}

	profileDir, err := GetProfileDir(profileName)
	if err != nil {
		return err
	}

	if err := os.Symlink(profileDir, opencodeLink); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
