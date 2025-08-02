package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
)

func TestTimeLocationDecodeHook(t *testing.T) {
	hook := TimeLocationDecodeHook()

	tests := []struct {
		name        string
		from        reflect.Type
		to          reflect.Type
		data        interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name:        "Valid timezone string to Location",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf((*time.Location)(nil)),
			data:        "America/New_York",
			expected:    "America/New_York", // We'll check the location name
			expectError: false,
		},
		{
			name:        "UTC timezone",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf((*time.Location)(nil)),
			data:        "UTC",
			expected:    "UTC",
			expectError: false,
		},
		{
			name:        "Invalid timezone",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf((*time.Location)(nil)),
			data:        "Invalid/Timezone",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Non-string input",
			from:        reflect.TypeOf(123),
			to:          reflect.TypeOf((*time.Location)(nil)),
			data:        123,
			expected:    123,
			expectError: false,
		},
		{
			name:        "Wrong target type",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(""),
			data:        "America/New_York",
			expected:    "America/New_York",
			expectError: false,
		},
		{
			name:        "Non-string data for string type",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf((*time.Location)(nil)),
			data:        123,
			expected:    123,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookFunc := hook.(func(reflect.Type, reflect.Type, interface{}) (interface{}, error))
			result, err := hookFunc(tt.from, tt.to, tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("TimeLocationDecodeHook() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("TimeLocationDecodeHook() unexpected error: %v", err)
				return
			}

			// For successful Location conversions, check the location name
			if tt.to == reflect.TypeOf((*time.Location)(nil)) && tt.from == reflect.TypeOf("") && tt.expected != 123 {
				if loc, ok := result.(*time.Location); ok {
					if loc.String() != tt.expected {
						t.Errorf("TimeLocationDecodeHook() location = %v, want %v", loc.String(), tt.expected)
					}
				} else {
					t.Errorf("TimeLocationDecodeHook() result is not *time.Location")
				}
			} else {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("TimeLocationDecodeHook() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestCMQTypeDecodeHook(t *testing.T) {
	hook := CMQTypeDecodeHook()

	tests := []struct {
		name        string
		from        reflect.Type
		to          reflect.Type
		data        interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name:        "Valid nats_streaming",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(CMQType(0)),
			data:        "nats_streaming",
			expected:    CMQNatsStreaming,
			expectError: false,
		},
		{
			name:        "Valid jet_stream",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(CMQType(0)),
			data:        "jet_stream",
			expected:    CMQJetStream,
			expectError: false,
		},
		{
			name:        "Valid nats",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(CMQType(0)),
			data:        "nats",
			expected:    CMQNats,
			expectError: false,
		},
		{
			name:        "Invalid CMQ type",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(CMQType(0)),
			data:        "invalid_type",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Non-string input",
			from:        reflect.TypeOf(123),
			to:          reflect.TypeOf(CMQType(0)),
			data:        123,
			expected:    123,
			expectError: false,
		},
		{
			name:        "Wrong target type",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(""),
			data:        "nats_streaming",
			expected:    "nats_streaming",
			expectError: false,
		},
		{
			name:        "Non-string data for string input type",
			from:        reflect.TypeOf(""),
			to:          reflect.TypeOf(CMQType(0)),
			data:        123,
			expected:    123,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookFunc := hook.(func(reflect.Type, reflect.Type, interface{}) (interface{}, error))
			result, err := hookFunc(tt.from, tt.to, tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("CMQTypeDecodeHook() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CMQTypeDecodeHook() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CMQTypeDecodeHook() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHooksIntegration(t *testing.T) {
	// Test that hooks work with mapstructure
	type TestConfig struct {
		Location *time.Location `mapstructure:"location"`
		CMQType  CMQType        `mapstructure:"cmq_type"`
	}

	input := map[string]interface{}{
		"location": "America/Los_Angeles",
		"cmq_type": "jet_stream",
	}

	var config TestConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			TimeLocationDecodeHook(),
			CMQTypeDecodeHook(),
		),
		Result: &config,
	})
	if err != nil {
		t.Fatalf("Failed to create decoder: %v", err)
	}

	err = decoder.Decode(input)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	// Check location
	if config.Location == nil {
		t.Error("Location should not be nil")
	} else if config.Location.String() != "America/Los_Angeles" {
		t.Errorf("Location = %v, want America/Los_Angeles", config.Location.String())
	}

	// Check CMQ type
	if config.CMQType != CMQJetStream {
		t.Errorf("CMQType = %v, want %v", config.CMQType, CMQJetStream)
	}
}

func TestHooksWithInvalidData(t *testing.T) {
	type TestConfig struct {
		Location *time.Location `mapstructure:"location"`
		CMQType  CMQType        `mapstructure:"cmq_type"`
	}

	input := map[string]interface{}{
		"location": "Invalid/Location",
		"cmq_type": "invalid_cmq",
	}

	var config TestConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			TimeLocationDecodeHook(),
			CMQTypeDecodeHook(),
		),
		Result: &config,
	})
	if err != nil {
		t.Fatalf("Failed to create decoder: %v", err)
	}

	err = decoder.Decode(input)
	if err == nil {
		t.Error("Expected decode to fail with invalid data")
	}
}

func TestHookTypes(t *testing.T) {
	// Test that hooks return the correct function type
	timeHook := TimeLocationDecodeHook()
	cmqHook := CMQTypeDecodeHook()

	// These should not panic - if they do, the function signatures are wrong
	_ = mapstructure.ComposeDecodeHookFunc(timeHook, cmqHook)
}
