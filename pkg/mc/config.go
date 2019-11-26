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

var conf *config

type config struct {
	*viper.Viper
	EnableRCON bool   `mapstructure:"enable-rcon"`
	ServerIP   string `mapstructure:"server-ip"`
	ServerPort int    `mapstructure:"server-port"`
	RCON       struct {
		Password string
		Port     int
	}
}

func init() {
	conf = &config{}
	conf.Viper = viper.New()

	conf.SetConfigName("server")
	conf.AddConfigPath(".")

	conf.SetDefault("enable-rcon", EnableRCON)
	conf.SetDefault("server-ip", ServerIP)
	conf.SetDefault("server-port", ServerPort)
	conf.SetDefault("rcon.password", "")
	conf.SetDefault("rcon.port", RCONPort)
}

// LoadConfig loads the server.properties file.
func LoadConfig() error {
	if err := conf.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}

		if werr := conf.SafeWriteConfigAs("server.properties"); werr != nil {
			return werr
		}
	}

	if err := conf.WriteConfig(); err != nil {
		return err
	}

	return conf.Unmarshal(conf)
}

// Config returns the current server configuration.
func Config() *config {
	return conf
}

// RCONAddr returns the address and port to which the RCON listener is bound.
func (conf *config) RCONAddr() string {
	if conf.EnableRCON && conf.RCON.Password != "" {
		return fmt.Sprintf("%s:%d", conf.ServerIP, conf.RCON.Port)
	}

	return ""
}

// ServerAddr returns the address and port to which the Minecraft listener is
// bound.
func (conf *config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", conf.ServerIP, conf.ServerPort)
}
