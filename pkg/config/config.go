package config

import (
	"fmt"
	"time"
)

type ApplicationType string

const (
	GoHabitApplicationType ApplicationType = "gohabit"
	NoHabitApplicationType ApplicationType = "nohabit"
)

// Tracing config struct
type Tracing struct {
	Enabled      bool    `mapstructure:"enabled"`
	AgentHost    string  `mapstructure:"agent_host"`
	AgentPort    string  `mapstructure:"agent_port"`
	SamplerRatio float64 `mapstructure:"sampler_ratio"`
}

// Logger config values
type Logger struct {
	Level string `mapstructure:"level"`
}

// HTTPServer config values
type HTTPServer struct {
	Listen            string        `mapstructure:"listen"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	SkipURLs          []string      `mapstructure:"skip_urls"`
}

// String represents most usable values for the HTTPServer config
func (i HTTPServer) String() string {
	return fmt.Sprintf("Listen: %s", i.Listen)
}

// Locale config about localization
type Locale struct {
	Default string `mapstructure:"default"`
	Path    string `mapstructure:"path"`
}

// NatsStreaming config values
type NatsStreaming struct {
	Address            string        `mapstructure:"address" validate:"required"`
	ConnectWait        time.Duration `mapstructure:"connect_wait" validate:"required"`
	PubAckWait         time.Duration `mapstructure:"pub_ack_wait" validate:"required"`
	MaxPubAcksInflight int           `mapstructure:"max_pub_acks_in_flight" validate:"required"`
	PingInterval       int           `mapstructure:"ping_interval" validate:"required"`
	PingMaxOut         int           `mapstructure:"ping_max_out" validate:"required"`
	ClusterID          string        `mapstructure:"cluster_id" validate:"required"`
	ClientID           string        `mapstructure:"client_id" validate:"required"`
}

// Nats is a struct to specify config for connecting to nats
type Nats struct {
	Type               CMQType       `mapstructure:"type"`
	Address            string        `mapstructure:"address" validate:"required"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	ConnectWait        time.Duration `mapstructure:"connect_wait" validate:"required"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout" validate:"required"`
	FlushTimeout       time.Duration `mapstructure:"flush_timeout" validate:"required"`
	FlusherTimeout     time.Duration `mapstructure:"flusher_timeout" validate:"required"`
	PingInterval       time.Duration `mapstructure:"ping_interval" validate:"required"`
	ConnectBufSize     int           `mapstructure:"connect_buf_size" validate:"required"`
	MaxChanLen         int           `mapstructure:"max_chan_len" validate:"required"`
	MaxPingOut         int           `mapstructure:"max_ping_out" validate:"required"`
	MaxReconnect       int           `mapstructure:"max_reconnect" validate:"required"`
	ClusterID          string        `mapstructure:"cluster_id" validate:"required"`
	ClientName         string        `mapstructure:"client_name" validate:"required"`
	ClientID           string        `mapstructure:"client_id" validate:"required"`
	PubAckWait         time.Duration `mapstructure:"pub_ack_wait" validate:"required"`
	MaxPubAcksInflight int           `mapstructure:"max_pub_acks_in_flight" validate:"required"`
	PingMaxOut         int           `mapstructure:"ping_max_out" validate:"required"`
}

// CMQType specifies type of CMQ
type CMQType int

const (
	CMQNatsStreaming CMQType = iota
	CMQJetStream
	CMQNats
)

// String is a method to convert the CMQType values to string.
func (c CMQType) String() string {
	return [...]string{
		"nats_streaming",
		"jet_stream",
		"nats",
	}[c]
}

// String represents most usable values for the NatsStreaming config
func (i NatsStreaming) String() string {
	return fmt.Sprintf("Address: %s, ClientID: %s, ClusterID: %v", i.Address, i.ClientID, i.ClusterID)
}

// Redis config struct
type Redis struct {
	Address            string        `mapstructure:"address"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	DB                 int           `mapstructure:"db"`
	MaxRetries         int           `mapstructure:"max_retries"`
	MinRetryBackoff    time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `mapstructure:"max_retry_backoff"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	PoolSize           int           `mapstructure:"pool_size"`
	MinIdleConns       int           `mapstructure:"min_idle_conns"`
	MaxConnAge         time.Duration `mapstructure:"max_conn_age"`
	PoolTimeout        time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idle_check_frequency"`
	DialRetry          int           `mapstructure:"dial_retry"`
	Sentinel           Sentinel      `mapstructure:"sentinel"`
}

// Sentinel specific configs for connecting to sentinel
type Sentinel struct {
	Enabled               bool     `mapstructure:"enabled"`
	MasterName            string   `mapstructure:"master_name"`
	Addresses             []string `mapstructure:"addresses"`
	SentinelPassword      string   `mapstructure:"sentinel_password"`
	SlaveOnly             bool     `mapstructure:"slave_only"`
	UseDisconnectedSlaves bool     `mapstructure:"use_disconnected_slaves"`
	QuerySentinelRandomly bool     `mapstructure:"query_sentinel_randomly"`
}

// String represents most usable values for the Redis config
func (i Redis) String() string {
	if i.Sentinel.Enabled {
		return fmt.Sprintf("Sentinel mode addresses: %v, DB: %d, MasterName: %s",
			i.Sentinel.Addresses,
			i.DB,
			i.Sentinel.MasterName)
	}

	return fmt.Sprintf("Single instance mode with address: %s, DB: %d", i.Address, i.DB)
}

// SQLDatabase is a configuration structure
type SQLDatabase struct {
	Driver                   string        `mapstructure:"driver"`
	Host                     string        `mapstructure:"host"`
	Port                     int           `mapstructure:"port"`
	DB                       string        `mapstructure:"db"`
	User                     string        `mapstructure:"user"`
	Password                 string        `mapstructure:"password"`
	MaxConn                  int           `mapstructure:"max_conn"`
	IdleConn                 int           `mapstructure:"idle_conn"`
	Timeout                  time.Duration `mapstructure:"timeout"`
	DialRetry                int           `mapstructure:"dial_retry"`
	StatementCacheCapacity   int           `mapstructure:"statement_cache_capacity"`
	DescriptionCacheCapacity int           `mapstructure:"description_cache_capacity"`
	DialTimeout              time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout              time.Duration `mapstructure:"read_timeout"`
	WriteTimeout             time.Duration `mapstructure:"write_timeout"`
	UpdateTimeout            time.Duration `mapstructure:"update_timeout"`
	DeleteTimeout            time.Duration `mapstructure:"delete_timeout"`
	QueryTimeout             time.Duration `mapstructure:"query_timeout"`
}

// DSN returns the Data Source Name for the SQLDatabase config
func (d SQLDatabase) DSN() string {
	if d.Driver == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&"+
			"collation=utf8mb4_general_ci", d.User, d.Password, d.Host, d.Port, d.DB)
	}
	if d.Driver == "postgresql" {
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.DB)
	}

	panic("SQLDatabase driver is not supported")
}

// String represents most usable values for the SQLDatabase config
func (d SQLDatabase) String() string {
	if d.Driver == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&"+
			"collation=utf8mb4_general_ci", d.User, "***", d.Host, d.Port, d.DB)
	}

	panic("SQLDatabase driver is not supported")
}

// RabbitMQ is the RabbitMQ config structure
type RabbitMQ struct {
	User        string        `mapstructure:"user"`
	Password    string        `mapstructure:"password"`
	Host        string        `mapstructure:"host"`
	Port        uint          `mapstructure:"port"`
	Vhost       string        `mapstructure:"vhost"`
	CtxTimeOut  time.Duration `mapstructure:"ctx_timeout"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
	DialRetry   uint          `mapstructure:"dial_retry"`
}

// DSN returns the Data Source Name for the RabbitMQ config proper for this RabbitMQ package "github.com/rabbitmq/amqp091-go"
func (r *RabbitMQ) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.User, r.Password, r.Host, r.Port, r.Vhost)
}

// BasicService is a struct to specify the config for a general service
type BasicService struct {
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// AuthTokenService is a struct to specify the config for a service with token
type AuthTokenService struct {
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
	Token   string        `mapstructure:"token"`
}

// AuthPasswordService is a struct to specify the config for a service with username and password
type AuthPasswordService struct {
	BaseURL  string        `mapstructure:"base_url"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	ClientID string        `mapstructure:"client_id"`
}

// Grpc is a struct to specify the config for application's gRPC server
type Grpc struct {
	Address string `mapstructure:"address"`
}
