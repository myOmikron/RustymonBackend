package configs

import (
	"errors"
	"fmt"
	"github.com/myOmikron/RustymonBackend/utils"
	"github.com/myOmikron/echotools/color"
	"strings"
)

var Config *RustymonConfig

var allowedDrivers = []string{"sqlite", "mysql", "postgresql"}

type Server struct {
	ListenAddress string
	ListenPort    uint16
}

type Database struct {
	Driver   string
	Name     string
	Host     string
	Port     uint16
	User     string
	Password string
}

type Rustymon struct {
	RegistrationDisabled bool
}

type RustymonConfig struct {
	Server   Server
	Database Database
	Rustymon Rustymon
}

type ConfigError struct {
	Section   string
	Parameter string
	Err       error
}

func (c *ConfigError) Error() string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s: %s\n",
		color.Colorize(color.RED, "Configuration failure:"),
		color.Colorize(color.BLUE, c.Section),
		c.Parameter,
		c.Err.Error(),
	)
}

func (conf *RustymonConfig) CheckConfig() *ConfigError {
	if conf.Server.ListenPort < 1 || conf.Server.ListenPort > 1<<15-1 {
		return &ConfigError{
			Err:       errors.New("invalid server port"),
			Section:   "[Server]",
			Parameter: "ListenPort",
		}
	}

	// Check database part
	if !utils.Contains(conf.Database.Driver, allowedDrivers) {
		return &ConfigError{
			Err:       errors.New(fmt.Sprintf("driver must be in range: %s", strings.Join(allowedDrivers, ","))),
			Section:   "[Database]",
			Parameter: "Driver",
		}
	}
	if conf.Database.Name == "" {
		return &ConfigError{
			Err:       errors.New("name must not be empty"),
			Section:   "[Database]",
			Parameter: "Name",
		}
	}
	switch conf.Database.Driver {
	case "mysql", "postgresql":
		if conf.Database.Port <= 0 || conf.Database.Port > 1<<15-1 {
			return &ConfigError{Err: errors.New("not a valid port"), Section: "[Database]", Parameter: "Port"}
		}
		if conf.Database.Host == "" {
			return &ConfigError{Err: errors.New("invalid host"), Section: "[Database]", Parameter: "Host"}
		}
		if conf.Database.User == "" {
			return &ConfigError{Err: errors.New("must not be empty"), Section: "[Database]", Parameter: "User"}
		}
		if conf.Database.Password == "" {
			return &ConfigError{Err: errors.New("must not be empty"), Section: "[Database]", Parameter: "Password"}
		}
	}

	return nil
}

func (conf *RustymonConfig) GetListenString() string {
	return conf.Server.ListenAddress + ":" + fmt.Sprint(conf.Server.ListenPort)
}
