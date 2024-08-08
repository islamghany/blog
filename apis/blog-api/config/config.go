package config

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// WEB
	// adding default values for the configuration
	ReadTimeout     time.Duration `mapstructure:"READ_TIMEOUT" default:"5s"`
	WriteTimeout    time.Duration `mapstructure:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout     time.Duration `mapstructure:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT" omit:"true" default:"20s"`
	APIHost         string        `mapstructure:"API_HOST" default:"0.0.0.0:8000"`
	DebugHost       string        `mapstructure:"DEBUG_HOST" mask:"true" default:"0.0.0.0:8080"`
	WHITELIST       []string      `mapstructure:"WHITELIST" default:"http://localhost:3000,http://localhost:5173"`
	// DB
	DBUser       string `mapstructure:"DB_USER" default:"blog"`
	DBPassword   string `mapstructure:"DB_PASSWORD" omit:"true" default:"islamghany"`
	DBHost       string `mapstructure:"DB_HOST" default:"localhost:5432"`
	DBName       string `mapstructure:"DB_NAME" default:"blog_db"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS" default:"25"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS" default:"25"`
	DisabelTLS   bool   `mapstructure:"DISABLE_TLS" default:"true"`
	// Security
	JWTSecret string `mapstructure:"JWT_SECRET" omit:"true" default:"secret key"`
}

func LoadConfig(isProduction bool) (config Config, err error) {
	// if production read from environment variables directly
	if isProduction {
		viper.SetDefault("READ_TIMEOUT", "5s")
		viper.SetDefault("WRITE_TIMEOUT", "10s")
		viper.SetDefault("IDLE_TIMEOUT", "120s")
		viper.SetDefault("SHUTDOWN_TIMEOUT", "20s")
		viper.SetDefault("API_HOST", "0.0.0.0:8000")
		viper.SetDefault("DEBUG_HOST", "0.0.0.0:8001")
		viper.SetDefault("WHITELIST", "http://localhost:3000,http://localhost:5173")
		viper.SetDefault("DB_USER", "blog")
		viper.SetDefault("DB_PASSWORD", "islamghany")
		viper.SetDefault("DB_HOST", "localhost:5432")
		viper.SetDefault("DB_NAME", "blog_db")
		viper.SetDefault("MAX_IDLE_CONNS", 25)
		viper.SetDefault("MAX_OPEN_CONNS", 25)
		viper.SetDefault("DISABLE_TLS", true)
		viper.SetDefault("JWT_SECRET", "secret key")
		viper.AutomaticEnv()
		err = viper.Unmarshal(&config)
		return
	}
	viper.AddConfigPath(".")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// Parse will omit or mask the keys that has omit or mask tags from the config
func (c Config) Parse() map[string]any {
	result := make(map[string]interface{})
	// Get the type of the Config struct
	configType := reflect.TypeOf(c)

	// Get the value of the Config struct
	configValue := reflect.ValueOf(c)

	// Iterate through the fields of the struct
	for i := 0; i < configType.NumField(); i++ {
		// Get the field
		field := configType.Field(i)
		// Get the field

		// Get the value of the field
		fieldValue := configValue.Field(i).Interface()

		// Check if the field has a omit tag
		if omitTag, ok := field.Tag.Lookup("omit"); ok && omitTag == "true" {
			continue // Skip this field
		}
		// Check if the field has a mask tag
		if maskTag, ok := field.Tag.Lookup("mask"); ok && maskTag == "true" {
			result[field.Name] = "************"
			continue
		}
		// Add the field value to the result
		// Add the field to the result
		result[field.Name] = fieldValue
	}

	return result
}
