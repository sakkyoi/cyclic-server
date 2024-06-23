package colonel

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
)

var Writ = &Config{}

type Config struct {
	Server   *ServerConfig   `koanf:"server"`
	Signup   *SignupConfig   `koanf:"signup"`
	SMTP     *SMTPConfig     `koanf:"smtp"`
	Database *DatabaseConfig `koanf:"database"`
	JWT      *JWTConfig      `koanf:"jwt"`
	Logger   *LoggerConfig   `koanf:"logger"`
}

type ServerConfig struct {
	Mode           string   `koanf:"mode"`
	Host           string   `koanf:"host"`
	Port           int      `koanf:"port"`
	TrustedProxies []string `koanf:"trustedProxies"`
}

type SignupConfig struct {
	Enabled      bool `koanf:"enabled"`
	Verification bool `koanf:"verification"`
}

type SMTPConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
}

type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
	Driver   string `koanf:"driver"`
}

type JWTConfig struct {
	PrivateKey string `koanf:"privateKey"`
	PublicKey  string `koanf:"publicKey"`
	Expiration int    `koanf:"expiration"`
}

type LoggerConfig struct {
	Level      string `koanf:"level"`
	File       string `koanf:"file"`
	MaxSize    int    `koanf:"maxSize"`
	MaxBackups int    `koanf:"maxBackups"`
	MaxAge     int    `koanf:"maxAge"`
	Compress   bool   `koanf:"compress"`
}

func Init() {
	k := koanf.New(".")

	// Load from yaml file
	f := file.Provider("config.yml")
	if err := k.Load(f, yaml.Parser()); err != nil {
		panic(err)
	}

	// Load from environment variables and merge into the loaded config
	if err := k.Load(env.Provider("CYCLIC_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CYCLIC_")), "_", ".", -1)
	}), nil); err != nil {
		panic(err)
	}

	// Unmarshal the whole config into the Writ variable
	if err := k.Unmarshal("", &Writ); err != nil {
		panic(err)
	}
}
