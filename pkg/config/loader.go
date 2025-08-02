package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration with the following precedence (highest to lowest):
// 1. Environment variables (with envPrefix)
// 2. External configuration file (if provided)
// 3. Built-in default configuration
//
// The config parameter must be a pointer to a configuration struct to populate.
//
// Example usage:
//
//	type AppConfig struct {
//	    Database DatabaseConfig `mapstructure:"database"`
//	    Redis    RedisConfig    `mapstructure:"redis"`
//	}
//
//	var config AppConfig
//	err := LoadConfig("MYAPP", "config.yaml", defaultYAML, &config)
//	if err != nil {
//	    return err
//	}
//	// config is now populated with merged configuration
func LoadConfig(
	envPrefix, configFilename string,
	builtinConfig []byte,
	config any,
	decodeHookFuncs ...mapstructure.DecodeHookFunc,
) error {
	// Validate config parameter
	if err := validateConfigParameter(config); err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// 1. Load built-in default configuration first (lowest priority)
	if err := v.ReadConfig(bytes.NewReader(builtinConfig)); err != nil {
		return fmt.Errorf("failed to load default configuration: %w", err)
	}

	// 2. Merge external config file if provided (medium priority)
	if configFilename != "" {
		v.SetConfigFile(configFilename)
		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("failed to merge config file %q: %w", configFilename, err)
		}
	}

	// 3. Environment variables are automatically handled by viper.AutomaticEnv() (highest priority)

	// Prepare decode hooks
	hooks := []mapstructure.DecodeHookFunc{
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		TimeLocationDecodeHook(),
		CMQTypeDecodeHook(),
	}
	hooks = append(hooks, decodeHookFuncs...)

	// Unmarshal final merged configuration
	if err := v.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(hooks...))); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return nil
}

// validateConfigParameter validates that the config parameter is a pointer to a struct
func validateConfigParameter(config any) error {
	if config == nil {
		return fmt.Errorf("config parameter cannot be nil")
	}

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("config parameter must be a pointer, got %T", config)
	}

	if rv.IsNil() {
		return fmt.Errorf("config parameter cannot be a nil pointer")
	}

	elem := rv.Elem()
	if !elem.CanSet() {
		return fmt.Errorf("config parameter must be a settable pointer")
	}

	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("config parameter must be a pointer to a struct, got pointer to %s", elem.Kind())
	}

	return nil
}
