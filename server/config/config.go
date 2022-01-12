package config

import "github.com/spf13/viper"

func LoadConfig(path string, filename string) (*FacadeConfig, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("server.port", 3126)
	viper.SetDefault("server.max_workers", 32)
	viper.SetDefault("server.spin_down", true)
	viper.SetDefault("server.name", "FacadeServer")

	viper.SetDefault("cache.ttl_ms", 1200000)
	viper.SetDefault("cache.clean_ms", 3600000)

	viper.SetDefault("response.raw", false)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	} else {
		config := new(FacadeConfig)
		err := viper.Unmarshal(config)

		return config, err
	}
}

// FacadeConfig is the configuration file marshalled.
type FacadeConfig struct {
	Connection ServerProperties `mapstructure:"server"`
	Cache      CacheProperties  `mapstructure:"cache"`
	Response   ResponseOptions  `mapstructure:"response"`
}

// ---------- Configuration structures ----------

// ServerProperties contains generic properties of the server.
type ServerProperties struct {
	Port       int    `mapstructure:"port"`        // Port for the server
	MaxWorkers int    `mapstructure:"max_workers"` // Max workers
	SpinDown   bool   `mapstructure:"spin_down"`   // Spin down idle workers?
	Name       string `mapstructure:"name"`        // Name for the service
}

// CacheProperties relates specifically to the cache.
type CacheProperties struct {
	TtlMs   int64 `mapstructure:"ttl_ms"`   // Time to live in ms
	CleanMs int64 `mapstructure:"clean_ms"` // How often to clean cache in ms
}

// ResponseOptions contains options regarding how Facade responds to requests.
type ResponseOptions struct {
	Raw bool `mapstructure:"raw"` // Return raw response rather than marshalled
}
