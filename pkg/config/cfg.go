package config

import "strings"

type Config struct {
	VaultAddr string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) UpdateVaultAddr(value string) {
	v := strings.TrimSuffix(value, "/")
	cfg.VaultAddr = v
}
