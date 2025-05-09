// pkg/config/config.go
package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `env:"SERVER"`
	Database DatabaseConfig `env:"DB"`
	JWT      JWTConfig      `env:"JWT"`
	VideoCfg VideoConfig    `env:"VIDEO"`
}

type ServerConfig struct {
	Port    string `env:"PORT" default:"8080"`
	Mode    string `env:"MODE" default:"debug"`
	Address string `env:"ADDRESS" default:"localhost"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST" default:"localhost"`
	Port     string `env:"PORT" default:"5432"`
	User     string `env:"USER" default:"postgres"`
	Password string `env:"PASSWORD" default:""`
	Name     string `env:"NAME" default:"fluxio"`
	SSLMode  string `env:"SSL_MODE" default:"disable"`
}

type JWTConfig struct {
	Secret string `env:"SECRET" default:""`
}

type VideoConfig struct {
	S3BucketName           string `env:"BUCKET_NAME" default:""`
	S3Region               string `env:"BUCKET_REGION" default:""`
	S3AccessKey            string `env:"BUCKET_ACCESS_KEY" default:""`
	S3SecretKey            string `env:"BUCKET_SECRET_KEY" default:""`
	S3UploadCallbackSecret string `env:"BUCKET_UPLOAD_CALLBACK_SECRET" default:""`
	S3Endpoint             string `env:"BUCKET_ENDPOINT" default:""`
}

const envPrefix = "FLUXIO"

// LoadConfig reads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{}
	populateConfigFromEnv(config, envPrefix)

	return config, nil
}

// GetDatabaseURL constructs a database URL from config
func (c *DatabaseConfig) GetDatabaseURL() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" sslmode=" + c.SSLMode
}

// populateConfigFromEnv uses reflection to populate a struct from environment variables
func populateConfigFromEnv(config interface{}, prefix string) {
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)

		// Get the env tag for this field
		envTag := field.Tag.Get("env")

		// If this field is a nested struct
		if field.Type.Kind() == reflect.Struct {
			// Determine the new prefix for nested fields
			newPrefix := prefix
			if envTag != "" {
				newPrefix = prefix + "_" + envTag
			}

			// Recursively populate nested struct
			populateConfigFromEnv(fieldValue.Addr().Interface(), newPrefix)
			continue
		}

		// For regular fields, get value from environment
		if envTag == "" {
			continue // Skip fields without env tag
		}

		// Construct the full environment variable name
		envName := prefix + "_" + envTag

		// Get default value
		defaultValue := field.Tag.Get("default")

		// Get the value from environment or use default
		value := getEnvOrDefault(envName, defaultValue)

		// Set the field value directly based on its type
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int, reflect.Int64:
			if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				fieldValue.SetInt(intValue)
			}
		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				fieldValue.SetBool(boolValue)
			}
		case reflect.Float64:
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				fieldValue.SetFloat(floatValue)
			}
		}
	}
}

// Helper function to get an environment variable or return a default value
func getEnvOrDefault(key, defaultValue string) string {
	// Make the key uppercase
	key = strings.ToUpper(key)

	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
