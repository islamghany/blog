package config

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// WEB
	ReadTimeout     time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT" omit:"true"`
	APIHost         string        `mapstructure:"API_HOST"`
	DebugHost       string        `mapstructure:"DEBUG_HOST" mask:"true"`
	WHITELIST       []string      `mapstructure:"WHITELIST"`
	// DB
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD" omit:"true"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBName       string `mapstructure:"DB_NAME"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	DisabelTLS   bool   `mapstructure:"DISABLE_TLS"`
	// Security
	JWTSecret string `mapstructure:"JWT_SECRET" omit:"true"`
}

func LoadConfig() (config Config, err error) {
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

		// Add the field to the result
		result[field.Name] = fieldValue
	}

	return result
}
