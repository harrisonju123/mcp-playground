// internal/config/load.go
package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

func Load() Config {
	k := koanf.New(".")
	// 1) defaults
	err := k.Load(confmap.Provider(map[string]any{
		"port":              defaultConfig.Port,
		"enable_reflection": defaultConfig.EnableReflection,
		"registry_url":      defaultConfig.RegistryURL,
		"otel_endpoint":     defaultConfig.OTelEndpoint,
	}, "."), nil)
	if err != nil {
		panic(err)
	}

	// 2) MCP_ env-vars override
	err = k.Load(env.Provider("MCP_", ".", nil), nil)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
