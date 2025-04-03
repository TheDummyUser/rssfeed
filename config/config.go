package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	DBPath string
}

var (
	loadEnvOnce sync.Once
	Envs        DbConfig
)

func init() {
	loadEnvOnce.Do(func() {
		loadEnv()
		Envs = DbConfig{
			DBPath: "rssdb.sqlite",
		}
	})
}

// Load environment variables from .env file
func loadEnv() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Warning: Could not determine current directory:", err)
		return
	}

	// Try loading .env from multiple possible locations
	paths := []string{
		filepath.Join(currentDir, ".env"),
		filepath.Join(currentDir, "..", ".env"),
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			fmt.Println("Loaded .env file from:", path)
			return
		}
	}

	fmt.Println("Warning: Could not load .env file. Checked:", paths)
}

// GetDSN returns the database connection string
func (c DbConfig) GetDSN() string {
	return c.DBPath
}

// Config fetches environment variables
func Config(key string) string {
	return os.Getenv(key)
}
