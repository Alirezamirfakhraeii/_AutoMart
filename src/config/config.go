package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// Config Structures
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	InternalPort string
	Port         string
	ExternalPort string
	RunMode      string
	Domain       string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}

// GetConfig 1. Main Execution Flow
// GetConfig: The main function that orchestrates fetching the directory,
// filename,  loading the configuration file, and parsing it into the Config struct.

func GetConfig() *Config {
	cfgDir := getConfigDir()
	cfgName := getConfigFileName(os.Getenv("APP_ENV"))

	v, err := LoadConfig(cfgName, "yml", cfgDir)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := ParsConfig(v)
	if err != nil {
		log.Fatalf("Erro in parse %v", err)
	}
	return cfg
}

// 2. Configuration Directory Determination
// getConfigDir: Finds and returns the absolute path of the directory
// where the configuration files are located (relative to this Go file).

func getConfigDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	return currentDir
}

// 3. Configuration File Naming
// getConfigFileName: Returns the base name of the configuration file (e.g., "config-development")
// based on the current APP_ENV environment variable.

func getConfigFileName(env string) string {
	if env == "docker" {
		return "config-docker"
	} else if env == "production" {
		return "config-production"
	} else {
		return "config-development"
	}
}

// LoadConfig 4. Loading the Configuration File (I/O)
// LoadConfig: Uses the Viper library to read the configuration file from the specified path
// and environment variables, returning a Viper object.

func LoadConfig(filename string, fileType string, configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType(fileType)
	v.AddConfigPath(configPath)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New(fmt.Sprintf("file Not Found in %s", configPath))
		}
		return nil, err
	}
	return v, nil
}

// ParsConfig 5. Parsing the Loaded Data
// ParsConfig: Unmarshals (converts) the data from the Viper object into the
// Go-defined 'Config' struct.

func ParsConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
