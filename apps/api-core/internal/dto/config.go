package dto

// ConfigResponse represents the server configuration response
type ConfigResponse struct {
	API      APIConfig      `json:"api"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
}

// APIConfig represents API server configuration
type APIConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Profile string `json:"profile"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// EnvironmentResponse represents the server environment information
type EnvironmentResponse struct {
	Profile  string `json:"profile"`
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Uptime   string `json:"uptime"`
}
