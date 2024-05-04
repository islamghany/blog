package config

import (
	"reflect"

	"github.com/spf13/viper"
)

func LoadConfig(dest any, path, name, ext string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(ext)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(dest)
	return
}

// Parse will omit or mask the keys that has omit or mask tags from the config
func Parse(c any) map[string]any {
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
