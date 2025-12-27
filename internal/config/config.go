package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read(filepaths ...string) (Config, error) {
	var config Config
	filePath, err := getConfigFilePath(filepaths...)
	if err != nil {
		return Config{}, fmt.Errorf("error config filepath: %w", err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error read config file: %w", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("error decode config file: %w", err)
	}
	return config, nil
}

func (c Config) SetUser(userName string, filepaths ...string) error {
	c.CurrentUserName = userName
	filePath, err := getConfigFilePath(filepaths...)
	if err != nil {
		return fmt.Errorf("error config filepath: %w", err)
	}
	if err = write(c, filePath); err != nil {
		return fmt.Errorf("error write config file: %w", err)
	}
	return nil
}

func getConfigFilePath(filepaths ...string) (string, error) {
	// Explicit override always wins
	if len(filepaths) > 0 {
		return filepaths[0], nil
	}

	// 1. Try executable-relative (works for built binaries)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidate := filepath.Join(exeDir, configFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// 2. Fallback: walk upward from CWD (works for go run .)
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, configFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("%s not found", configFileName)
}

func write(c Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error overwrite config file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(&c); err != nil {
		return fmt.Errorf("error encode config file: %w", err)
	}
	return nil

}
