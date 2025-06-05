// internal/config/config.go
package config

type Config struct {
	Port             int    `koanf:"port"`              // MCP_PORT
	EnableReflection bool   `koanf:"enable_reflection"` // MCP_ENABLE_REFLECTION
	RegistryURL      string `koanf:"registry_url"`      // MCP_REGISTRY_URL
	OTelEndpoint     string `koanf:"otel_endpoint"`     // MCP_OTEL_ENDPOINT
}

// sane built-ins for unit tests / go run
var defaultConfig = Config{
	Port:             50051,
	EnableReflection: true,
	RegistryURL:      "./servers.yaml",
	OTelEndpoint:     "",
}
