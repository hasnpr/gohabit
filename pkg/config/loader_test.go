package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Test configuration structs
type TestConfig struct {
	Name     string        `mapstructure:"name"`
	Port     int           `mapstructure:"port"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Features []string      `mapstructure:"features"`
	Database TestDatabase  `mapstructure:"database"`
}

type TestDatabase struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var defaultConfigYAML = `
name: "default-app"
port: 8080
timeout: "30s"
features:
  - "feature1"
  - "feature2"
database:
  host: "localhost"
  username: "default_user"
  password: "default_pass"
`

var overrideConfigYAML = `
name: "override-app"
port: 9090
database:
  host: "override-host"
  username: "override_user"
`

func TestLoadConfig_Success(t *testing.T) {
	var config TestConfig
	err := LoadConfig("TEST", "", []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify config was loaded correctly
	if config.Name != "default-app" {
		t.Errorf("Config.Name = %v, want default-app", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Config.Port = %v, want 8080", config.Port)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("Config.Timeout = %v, want 30s", config.Timeout)
	}
	if len(config.Features) != 2 {
		t.Errorf("Config.Features length = %v, want 2", len(config.Features))
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Config.Database.Host = %v, want localhost", config.Database.Host)
	}
}

func TestLoadConfig_WithExternalFile(t *testing.T) {
	// Create temporary config file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(overrideConfigYAML)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	var config TestConfig
	err = LoadConfig("TEST", tmpFile.Name(), []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify config was merged correctly (external file overrides defaults)
	if config.Name != "override-app" {
		t.Errorf("Config.Name = %v, want override-app", config.Name)
	}
	if config.Port != 9090 {
		t.Errorf("Config.Port = %v, want 9090", config.Port)
	}
	if config.Database.Host != "override-host" {
		t.Errorf("Config.Database.Host = %v, want override-host", config.Database.Host)
	}
	// Values not in override file should come from defaults
	if config.Timeout != 30*time.Second {
		t.Errorf("Config.Timeout = %v, want 30s", config.Timeout)
	}
	if config.Database.Password != "default_pass" {
		t.Errorf("Config.Database.Password = %v, want default_pass", config.Database.Password)
	}
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("TEST_NAME", "env-app")
	os.Setenv("TEST_PORT", "7070")
	os.Setenv("TEST_DATABASE_HOST", "env-host")
	defer func() {
		os.Unsetenv("TEST_NAME")
		os.Unsetenv("TEST_PORT")
		os.Unsetenv("TEST_DATABASE_HOST")
	}()

	var config TestConfig
	err := LoadConfig("TEST", "", []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Environment variables should override defaults
	if config.Name != "env-app" {
		t.Errorf("Config.Name = %v, want env-app", config.Name)
	}
	if config.Port != 7070 {
		t.Errorf("Config.Port = %v, want 7070", config.Port)
	}
	if config.Database.Host != "env-host" {
		t.Errorf("Config.Database.Host = %v, want env-host", config.Database.Host)
	}
	// Values without env vars should come from defaults
	if config.Timeout != 30*time.Second {
		t.Errorf("Config.Timeout = %v, want 30s", config.Timeout)
	}
}

func TestLoadConfig_PrecedenceOrder(t *testing.T) {
	// Create temporary config file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(overrideConfigYAML)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Set environment variable (highest priority)
	os.Setenv("TEST_NAME", "env-priority")
	defer os.Unsetenv("TEST_NAME")

	var config TestConfig
	err = LoadConfig("TEST", tmpFile.Name(), []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Environment should win over file and defaults
	if config.Name != "env-priority" {
		t.Errorf("Config.Name = %v, want env-priority (env should override file and defaults)", config.Name)
	}
	// File should win over defaults
	if config.Port != 9090 {
		t.Errorf("Config.Port = %v, want 9090 (file should override defaults)", config.Port)
	}
	// Default values should be used when not overridden
	if config.Timeout != 30*time.Second {
		t.Errorf("Config.Timeout = %v, want 30s (should use default)", config.Timeout)
	}
}

func TestLoadConfig_InvalidBuiltinConfig(t *testing.T) {
	var config TestConfig
	err := LoadConfig("TEST", "", []byte("invalid yaml content"), &config)

	if err == nil {
		t.Error("LoadConfig() should fail with invalid builtin config")
	}
	if !strings.Contains(err.Error(), "failed to load default configuration") {
		t.Errorf("Error should mention default configuration: %v", err)
	}
}

func TestLoadConfig_InvalidExternalFile(t *testing.T) {
	var config TestConfig
	err := LoadConfig("TEST", "nonexistent-file.yaml", []byte(defaultConfigYAML), &config)

	if err == nil {
		t.Error("LoadConfig() should fail with nonexistent external file")
	}
	if !strings.Contains(err.Error(), "failed to merge config file") {
		t.Errorf("Error should mention config file merging: %v", err)
	}
}

func TestLoadConfig_UnmarshalError(t *testing.T) {
	invalidConfigYAML := `
name: "test"
port: "invalid-port"  # This should be an integer
`
	var config TestConfig
	err := LoadConfig("TEST", "", []byte(invalidConfigYAML), &config)

	if err == nil {
		t.Error("LoadConfig() should fail with unmarshal error")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal configuration") {
		t.Errorf("Error should mention unmarshal failure: %v", err)
	}
}

func TestLoadConfig_WithCustomDecodeHooks(t *testing.T) {
	// Create a custom decode hook for testing
	customHook := func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t.Kind() == reflect.String {
			if str, ok := data.(string); ok && str == "CUSTOM" {
				return "TRANSFORMED", nil
			}
		}
		return data, nil
	}

	configYAML := `
name: "CUSTOM"
port: 8080
`

	var config TestConfig
	err := LoadConfig("TEST", "", []byte(configYAML), &config, customHook)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Name != "TRANSFORMED" {
		t.Errorf("Config.Name = %v, want TRANSFORMED (custom hook should transform)", config.Name)
	}
}

func TestValidateConfigParameter(t *testing.T) {
	tests := []struct {
		name          string
		config        interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid pointer to struct",
			config:      &TestConfig{},
			expectError: false,
		},
		{
			name:          "Nil config",
			config:        nil,
			expectError:   true,
			errorContains: "cannot be nil",
		},
		{
			name:          "Non-pointer config",
			config:        TestConfig{},
			expectError:   true,
			errorContains: "must be a pointer",
		},
		{
			name:          "Nil pointer",
			config:        (*TestConfig)(nil),
			expectError:   true,
			errorContains: "cannot be a nil pointer",
		},
		{
			name:          "Pointer to non-struct",
			config:        new(string),
			expectError:   true,
			errorContains: "must be a pointer to a struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfigParameter(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("validateConfigParameter() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error should contain %q, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("validateConfigParameter() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadConfig_ValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		config interface{}
	}{
		{"Nil config", nil},
		{"Non-pointer config", TestConfig{}},
		{"Nil pointer", (*TestConfig)(nil)},
		{"Pointer to non-struct", new(string)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadConfig("TEST", "", []byte(defaultConfigYAML), tt.config)
			if err == nil {
				t.Errorf("LoadConfig() should fail with invalid config parameter")
			}
		})
	}
}

func TestLoadConfig_EmptyEnvPrefix(t *testing.T) {
	var config TestConfig
	err := LoadConfig("", "", []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() should work with empty env prefix: %v", err)
	}

	if config.Name != "default-app" {
		t.Errorf("Config should be loaded correctly even with empty env prefix")
	}
}

func TestLoadConfig_EnvironmentVariableReplacer(t *testing.T) {
	// Test that dots and dashes in config keys are replaced with underscores in env vars
	os.Setenv("TEST_DATABASE_HOST", "env-replaced-host")
	defer os.Unsetenv("TEST_DATABASE_HOST")

	var config TestConfig
	err := LoadConfig("TEST", "", []byte(defaultConfigYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Database.Host != "env-replaced-host" {
		t.Errorf("Config.Database.Host = %v, want env-replaced-host", config.Database.Host)
	}
}

// Test that the built-in decode hooks work correctly
func TestLoadConfig_BuiltinDecodeHooks(t *testing.T) {
	type ConfigWithHooks struct {
		Duration time.Duration `mapstructure:"duration"`
		List     []string      `mapstructure:"list"`
	}

	configYAML := `
duration: "5m"
list: "item1,item2,item3"
`

	var config ConfigWithHooks
	err := LoadConfig("TEST", "", []byte(configYAML), &config)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Duration != 5*time.Minute {
		t.Errorf("Config.Duration = %v, want 5m", config.Duration)
	}

	expectedList := []string{"item1", "item2", "item3"}
	if !reflect.DeepEqual(config.List, expectedList) {
		t.Errorf("Config.List = %v, want %v", config.List, expectedList)
	}
}
