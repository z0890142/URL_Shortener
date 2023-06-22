package config

import (
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

var globalConfig *Config
var configOnce sync.Once

// ResetConfig set config to Nil, used for tests
func ResetConfig() {
	globalConfig = nil
}

// GetConfig 獲取該服務相關配置
func GetConfig() *Config {
	configOnce.Do(func() {
		globalConfig = &Config{}
	})
	return globalConfig
}

// Config 該服務相關配置
type Config struct {
	Env              string         `mapstructure:"ENV"`
	Service          Service        `mapstructure:"SERVICE"`
	LogLevel         string         `mapstructure:"LOG_LEVEL"`
	LogFile          string         `mapstructure:"LOG_FILE"`
	EndPoints        EndPoints      `mapstructure:"ENDPOINTS"`
	Databases        DatabaseOption `mapstructure:"DATABASES"`
	MaxRetry         int            `mapstructure:"MAX_RETRY"`
	Redis            RedisOption    `mapstructure:"REDIS"`
	HashPoolSize     int            `mapstructure:"HASH_POOL_SIZE"`
	EnableKeyService bool           `mapstructure:"ENABLE_KEY_SERVICE"`
}

// Service defines service configuration struct.
type Service struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

type EndPoints struct {
	KeyServer Endpoint `mapstructure:"KEY_SERVER"`
}
type Endpoint struct {
	EnableGRCP    bool              `mapstructure:"ENABLE_GRCP"`
	Grpc          GrpcEndpoint      `mapstructure:"GRPC"`
	Http          HttpEndpoint      `mapstructure:"HTTP"`
	ExtensionInfo map[string]string `mapstructure:"EXTENSION_INFO"`
}

type HttpEndpoint struct {
	Host      string `mapstructure:"HOST"`
	Port      string `mapstructure:"PORT"`
	EnableTls bool   `mapstructure:"ENABLE_TLS"`
}

type GrpcEndpoint struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	Insecure bool   `mapstructure:"INSECURE"`
}

type DatabaseOption struct {
	Driver       string        `mapstructure:"DRIVER"`
	Host         string        `mapstructure:"HOST"`
	Port         uint16        `mapstructure:"PORT"`
	Username     string        `mapstructure:"USERNAME"`
	Password     string        `mapstructure:"PASSWORD"`
	DBName       string        `mapstructure:"DBNAME"`
	Timezone     string        `mapstructure:"TIMEZONE"`
	Charset      string        `mapstructure:"CHARSET"`
	PoolSize     int           `mapstructure:"POOL_SIZE"`
	Timeout      time.Duration `mapstructure:"TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
}

type RedisOption struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Password string `mapstructure:"PASSWORD"`
}

// MarshalLogObject use for logger *zap.Logger
func (o DatabaseOption) MarshalLogObject(en zapcore.ObjectEncoder) error {
	en.AddString("driver", o.Driver)
	en.AddString("host", o.Host)
	en.AddUint16("port", o.Port)
	en.AddString("username", o.Username)
	en.AddString("password", "********")
	en.AddString("dbname", o.DBName)

	if len(o.Timezone) > 0 {
		en.AddString("timezone", o.Timezone)
	}

	if len(o.Charset) > 0 {
		en.AddString("charset", o.Charset)
	}

	if o.PoolSize > 0 {
		en.AddInt("pool_size", o.PoolSize)
	}

	if o.Timeout > 0 {
		en.AddString("timeout", o.Timeout.String())
	}

	if o.ReadTimeout > 0 {
		en.AddString("read_timeout", o.ReadTimeout.String())
	}

	if o.WriteTimeout > 0 {
		en.AddString("write_timeout", o.WriteTimeout.String())
	}

	return nil
}
