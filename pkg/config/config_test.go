package config

import (
	"testing"
	"time"
)

func TestApplicationType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		appType  ApplicationType
		expected string
	}{
		{"GoHabit type", GoHabitApplicationType, "gohabit"},
		{"NoHabit type", NoHabitApplicationType, "nohabit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.appType) != tt.expected {
				t.Errorf("ApplicationType = %v, want %v", string(tt.appType), tt.expected)
			}
		})
	}
}

func TestDriver_Constants(t *testing.T) {
	if MySQLDriver != "mysql" {
		t.Errorf("MySQLDriver = %v, want mysql", MySQLDriver)
	}
	if PostgreSQLDriver != "postgresql" {
		t.Errorf("PostgreSQLDriver = %v, want postgresql", PostgreSQLDriver)
	}
}

func TestHTTPServer_String(t *testing.T) {
	server := HTTPServer{
		Listen:            ":8080",
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		SkipURLs:          []string{"/health", "/metrics"},
	}

	expected := "Listen: :8080"
	result := server.String()

	if result != expected {
		t.Errorf("HTTPServer.String() = %v, want %v", result, expected)
	}
}

func TestCMQType_String(t *testing.T) {
	tests := []struct {
		name     string
		cmqType  CMQType
		expected string
	}{
		{"NATS Streaming", CMQNatsStreaming, "nats_streaming"},
		{"JetStream", CMQJetStream, "jet_stream"},
		{"NATS", CMQNats, "nats"},
		{"Invalid type", CMQType(999), "invalid CMQType(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmqType.String()
			if result != tt.expected {
				t.Errorf("CMQType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNatsStreaming_String(t *testing.T) {
	nats := NatsStreaming{
		Address:   "nats://localhost:4222",
		ClientID:  "test-client",
		ClusterID: "test-cluster",
	}

	expected := "Address: nats://localhost:4222, ClientID: test-client, ClusterID: test-cluster"
	result := nats.String()

	if result != expected {
		t.Errorf("NatsStreaming.String() = %v, want %v", result, expected)
	}
}

func TestRedis_String(t *testing.T) {
	t.Run("Single instance mode", func(t *testing.T) {
		redis := Redis{
			Address: "localhost:6379",
			DB:      0,
			Sentinel: Sentinel{
				Enabled: false,
			},
		}

		expected := "Single instance mode with address: localhost:6379, DB: 0"
		result := redis.String()

		if result != expected {
			t.Errorf("Redis.String() = %v, want %v", result, expected)
		}
	})

	t.Run("Sentinel mode", func(t *testing.T) {
		redis := Redis{
			DB: 1,
			Sentinel: Sentinel{
				Enabled:    true,
				MasterName: "mymaster",
				Addresses:  []string{"sentinel1:26379", "sentinel2:26379"},
			},
		}

		expected := "Sentinel mode addresses: [sentinel1:26379 sentinel2:26379], DB: 1, MasterName: mymaster"
		result := redis.String()

		if result != expected {
			t.Errorf("Redis.String() = %v, want %v", result, expected)
		}
	})
}

func TestSQLDatabase_DSN(t *testing.T) {
	tests := []struct {
		name        string
		db          SQLDatabase
		expectedDSN string
		expectError bool
	}{
		{
			name: "MySQL DSN",
			db: SQLDatabase{
				Driver:   MySQLDriver,
				Host:     "localhost",
				Port:     3306,
				DB:       "testdb",
				User:     "testuser",
				Password: "testpass",
			},
			expectedDSN: "testuser:testpass@tcp(localhost:3306)/testdb?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_general_ci",
			expectError: false,
		},
		{
			name: "PostgreSQL DSN",
			db: SQLDatabase{
				Driver:   PostgreSQLDriver,
				Host:     "localhost",
				Port:     5432,
				DB:       "testdb",
				User:     "testuser",
				Password: "testpass",
			},
			expectedDSN: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
			expectError: false,
		},
		{
			name: "Unsupported driver",
			db: SQLDatabase{
				Driver: "sqlite",
			},
			expectedDSN: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := tt.db.DSN()

			if tt.expectError {
				if err == nil {
					t.Errorf("SQLDatabase.DSN() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("SQLDatabase.DSN() unexpected error: %v", err)
				return
			}

			if dsn != tt.expectedDSN {
				t.Errorf("SQLDatabase.DSN() = %v, want %v", dsn, tt.expectedDSN)
			}
		})
	}
}

func TestSQLDatabase_String(t *testing.T) {
	tests := []struct {
		name     string
		db       SQLDatabase
		expected string
	}{
		{
			name: "MySQL string representation",
			db: SQLDatabase{
				Driver:   MySQLDriver,
				Host:     "localhost",
				Port:     3306,
				DB:       "testdb",
				User:     "testuser",
				Password: "testpass",
			},
			expected: "testuser:***@tcp(localhost:3306)/testdb?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_general_ci",
		},
		{
			name: "PostgreSQL string representation",
			db: SQLDatabase{
				Driver:   PostgreSQLDriver,
				Host:     "localhost",
				Port:     5432,
				DB:       "testdb",
				User:     "testuser",
				Password: "testpass",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "Unsupported driver",
			db: SQLDatabase{
				Driver: "sqlite",
			},
			expected: "SQLDatabase driver is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.db.String()
			if result != tt.expected {
				t.Errorf("SQLDatabase.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRabbitMQ_DSN(t *testing.T) {
	rabbit := &RabbitMQ{
		User:     "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     5672,
		Vhost:    "/",
	}

	expected := "amqp://guest:guest@localhost:5672//"
	result := rabbit.DSN()

	if result != expected {
		t.Errorf("RabbitMQ.DSN() = %v, want %v", result, expected)
	}
}

func TestRabbitMQ_DSN_CustomVhost(t *testing.T) {
	rabbit := &RabbitMQ{
		User:     "admin",
		Password: "secret",
		Host:     "rabbitmq.example.com",
		Port:     5672,
		Vhost:    "production",
	}

	expected := "amqp://admin:secret@rabbitmq.example.com:5672/production"
	result := rabbit.DSN()

	if result != expected {
		t.Errorf("RabbitMQ.DSN() = %v, want %v", result, expected)
	}
}

func TestTracing_Struct(t *testing.T) {
	tracing := Tracing{
		Enabled:      true,
		AgentHost:    "jaeger",
		AgentPort:    "6831",
		SamplerRatio: 0.1,
	}

	if !tracing.Enabled {
		t.Error("Tracing.Enabled should be true")
	}
	if tracing.AgentHost != "jaeger" {
		t.Errorf("Tracing.AgentHost = %v, want jaeger", tracing.AgentHost)
	}
	if tracing.SamplerRatio != 0.1 {
		t.Errorf("Tracing.SamplerRatio = %v, want 0.1", tracing.SamplerRatio)
	}
}

func TestLogger_Struct(t *testing.T) {
	logger := Logger{
		Level: "debug",
	}

	if logger.Level != "debug" {
		t.Errorf("Logger.Level = %v, want debug", logger.Level)
	}
}

func TestBasicService_Struct(t *testing.T) {
	service := BasicService{
		BaseURL: "https://api.example.com",
		Timeout: 30 * time.Second,
	}

	if service.BaseURL != "https://api.example.com" {
		t.Errorf("BasicService.BaseURL = %v, want https://api.example.com", service.BaseURL)
	}
	if service.Timeout != 30*time.Second {
		t.Errorf("BasicService.Timeout = %v, want 30s", service.Timeout)
	}
}

func TestAuthTokenService_Struct(t *testing.T) {
	service := AuthTokenService{
		BaseURL: "https://api.example.com",
		Timeout: 30 * time.Second,
		Token:   "secret-token",
	}

	if service.Token != "secret-token" {
		t.Errorf("AuthTokenService.Token = %v, want secret-token", service.Token)
	}
}

func TestAuthPasswordService_Struct(t *testing.T) {
	service := AuthPasswordService{
		BaseURL:  "https://api.example.com",
		Username: "admin",
		Password: "secret",
		ClientID: "client-123",
	}

	if service.Username != "admin" {
		t.Errorf("AuthPasswordService.Username = %v, want admin", service.Username)
	}
	if service.ClientID != "client-123" {
		t.Errorf("AuthPasswordService.ClientID = %v, want client-123", service.ClientID)
	}
}

func TestGrpc_Struct(t *testing.T) {
	grpc := Grpc{
		Address: ":9090",
	}

	if grpc.Address != ":9090" {
		t.Errorf("Grpc.Address = %v, want :9090", grpc.Address)
	}
}
