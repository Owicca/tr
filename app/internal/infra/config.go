package infra

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type SessionsConfig struct {
	AuthenticationKey string
	EncryptionKey     string
	Key               string
}

type Config struct {
	CfgPath    string
	HttpHost   string
	HttpPort   string
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
	Logger     zap.Config
	Sessions   SessionsConfig
}

func (c Config) String() string {
	cfgStr, err := json.Marshal(c.CfgPath)
	if err != nil {
		// fmt.Errorf("Could not serialize config struct %s (%s)", c.CfgPath, err)
		return ""
	}

	return string(cfgStr)
}

func LoadConfig(path string) (Config, error) {
	cfg := NewConfig(path)

	file, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("Could not open config file %s (%s)", path, err)
	}
	dec := json.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("Could not decode contents of config file %s (%s)", path, err)
	}

	return cfg, nil
}

func NewConfig(path string) Config {
	var cfg = Config{
		CfgPath: path,
	}

	return cfg
}
