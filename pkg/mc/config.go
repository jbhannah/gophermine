package mc

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the vanilla Minecraft server configuration, defined by default in
// a server.properties file in the current directory.
type Config struct {
	EnableRCON bool   `mapstructure:"enable-rcon"`
	ServerIP   string `mapstructure:"server-ip"`
	ServerPort int    `mapstructure:"server-port"`
	RCON       struct {
		Password string
		Port     int
	}
}

// NewServerConfig reads the vanilla Minecraft server configuration file and
// returns it unamrshaled into a Config.
func NewServerConfig() (*Config, error) {
	viper.SetConfigName("server")

	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetDefault("enable-rcon", false)
	viper.SetDefault("server-ip", "")
	viper.SetDefault("server-port", 25565)
	viper.SetDefault("rcon.password", "")
	viper.SetDefault("rcon.port", 25575)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}

		if werr := viper.SafeWriteConfigAs("server.properties"); werr != nil {
			return nil, werr
		}
	}

	if werr := viper.WriteConfig(); werr != nil {
		return nil, werr
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// RCONAddr returns the address and port to which the RCON listener is bound.
func (config *Config) RCONAddr() string {
	if config.EnableRCON && config.RCON.Password != "" {
		return fmt.Sprintf("%s:%d", config.ServerIP, config.RCON.Port)
	}

	return ""
}

// ServerAddr returns the address and port to which the Minecraft listener is
// bound.
func (config *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", config.ServerIP, config.ServerPort)
}
