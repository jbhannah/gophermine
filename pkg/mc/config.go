package mc

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration default values
const (
	EnableRCON = false
	ServerIP   = ""
	ServerPort = 25565
	RCONPort   = 25575
)

// Config is the vanilla Minecraft server configuration, defined by default in
// a server.properties file in the current directory.
type Config struct {
	*viper.Viper
	EnableRCON bool   `mapstructure:"enable-rcon"`
	ServerIP   string `mapstructure:"server-ip"`
	ServerPort int    `mapstructure:"server-port"`
	RCON       struct {
		Password string
		Port     int
	}
}

// NewConfig reads the vanilla Minecraft server configuration file and
// returns it unamrshaled into a Config.
func NewConfig() (*Config, error) {
	config := &Config{}
	config.Viper = viper.New()

	config.SetConfigName("server")
	config.AddConfigPath(".")

	config.SetDefault("enable-rcon", EnableRCON)
	config.SetDefault("server-ip", ServerIP)
	config.SetDefault("server-port", ServerPort)
	config.SetDefault("rcon.password", "")
	config.SetDefault("rcon.port", RCONPort)

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}

		if werr := config.SafeWriteConfigAs("server.properties"); werr != nil {
			return nil, werr
		}
	}

	if werr := config.WriteConfig(); werr != nil {
		return nil, werr
	}

	if err := config.Unmarshal(config); err != nil {
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
