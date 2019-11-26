package mc

import (
	"fmt"

	"github.com/spf13/viper"
)

// Properties default values
const (
	EnableRCON = false
	ServerIP   = ""
	ServerPort = 25565
	RCONPort   = 25575
)

var props *properties

type properties struct {
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
	props = &properties{}
	props.Viper = viper.New()

	props.SetConfigName("server")
	props.AddConfigPath(".")

	props.SetDefault("enable-rcon", EnableRCON)
	props.SetDefault("server-ip", ServerIP)
	props.SetDefault("server-port", ServerPort)
	props.SetDefault("rcon.password", "")
	props.SetDefault("rcon.port", RCONPort)
}

// LoadProperties loads the server.properties file.
func LoadProperties() error {
	if err := props.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}

		if werr := props.SafeWriteConfigAs("server.properties"); werr != nil {
			return werr
		}
	}

	if err := props.WriteConfig(); err != nil {
		return err
	}

	return props.Unmarshal(props)
}

// Properties returns the current server configuration.
func Properties() *properties {
	return props
}

// RCONAddr returns the address and port to which the RCON listener is bound.
func (p *properties) RCONAddr() string {
	if p.EnableRCON && p.RCON.Password != "" {
		return fmt.Sprintf("%s:%d", p.ServerIP, p.RCON.Port)
	}

	return ""
}

// ServerAddr returns the address and port to which the Minecraft listener is
// bound.
func (p *properties) ServerAddr() string {
	return fmt.Sprintf("%s:%d", p.ServerIP, p.ServerPort)
}
