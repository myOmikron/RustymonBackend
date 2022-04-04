package configs

import (
	"errors"
	"fmt"
	"github.com/myOmikron/RustymonBackend/utils"
	"github.com/myOmikron/echotools/color"
	"net/url"
	"os"
	"strings"
)

var allowedDrivers = []string{"sqlite", "mysql", "postgresql"}

type Server struct {
	ListenAddress           string
	ListenPort              uint16
	PublicURI               string
	AllowedHosts            []string
	UseForwardedProtoHeader bool
	TemplateDir             string
	CLIUnixSocket           string
}

type Logging struct {
	LogFile        string
	LogQueueSize   int
	LogMaxCapacity int
	LogMaxDays     int
	LogMaxBackups  int
}

type Database struct {
	Driver   string
	Name     string
	Host     string
	Port     uint16
	User     string
	Password string
}

type Mail struct {
	Host     string
	Port     uint16
	From     string
	User     string
	Password string
}

type Rustymon struct {
	RegistrationDisabled bool
}

type RustymonConfig struct {
	Server   Server
	Logging  Logging
	Database Database
	Mail     Mail
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
			Err:       errors.New("invalid app port"),
			Section:   "[Server]",
			Parameter: "ListenPort",
		}
	}

	if !strings.HasPrefix(conf.Server.PublicURI, "https://") &&
		!strings.HasPrefix(conf.Server.PublicURI, "http://") {
		return &ConfigError{
			Err:       errors.New("invalid public uri prefix. Only https:// or http:// are valid prefixes"),
			Section:   "[Server]",
			Parameter: "PublicURI",
		}
	}

	if _, err := url.Parse(conf.Server.PublicURI); err != nil {
		return &ConfigError{
			Err:       errors.New("invalid public uri"),
			Section:   "[Server]",
			Parameter: "PublicURI",
		}
	}

	if _, err := os.Stat(conf.Server.TemplateDir); err != nil {
		return &ConfigError{
			Err:       err,
			Section:   "[Server]",
			Parameter: "TemplateDir",
		}
	}

	if len(conf.Server.AllowedHosts) == 0 {
		return &ConfigError{
			Err:       errors.New("empty value is forbidden, as the server will not respond to anything"),
			Section:   "[Server]",
			Parameter: "AllowedHosts",
		}
	}
	for _, allowedHost := range conf.Server.AllowedHosts {
		if !strings.HasPrefix(allowedHost, "https://") && !strings.HasPrefix(allowedHost, "http://") {
			return &ConfigError{
				Err:       errors.New("must be starting with either http:// or https://"),
				Section:   "[Section]",
				Parameter: "AllowedHosts",
			}
		}
		if _, err := url.Parse(allowedHost); err != nil {
			return &ConfigError{
				Err:       err,
				Section:   "[Section]",
				Parameter: "AllowedHosts",
			}
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

	if _, serr := os.Stat(conf.Logging.LogFile); errors.Is(serr, os.ErrNotExist) {
		if testlog, ferr := os.Create(conf.Logging.LogFile); ferr != nil {
			return &ConfigError{
				Err:       ferr,
				Section:   "[Logging]",
				Parameter: "LogFile",
			}
		} else {
			// We can leave the file empty, as logging will open it anyway
			testlog.Close()
		}
	} else if errors.Is(serr, os.ErrPermission) {
		c := ConfigError{
			Err:       serr,
			Section:   "[Logging]",
			Parameter: "LogFile",
		}
		fmt.Println(c.Error())
		os.Exit(1)
	}

	if conf.Logging.LogMaxDays <= 0 {
		return &ConfigError{
			Err:       errors.New("must be > 0"),
			Section:   "[Logging]",
			Parameter: "LogMaxDays",
		}
	}

	if conf.Logging.LogMaxCapacity <= 0 {
		return &ConfigError{
			Err:       errors.New("must be > 0"),
			Section:   "[Logging]",
			Parameter: "LogMaxCapacity",
		}
	}

	if conf.Logging.LogMaxBackups < 0 {
		return &ConfigError{
			Err:       errors.New("must be >= 0"),
			Section:   "[Logging]",
			Parameter: "LogMaxBackups",
		}
	}

	if conf.Logging.LogQueueSize <= 0 {
		return &ConfigError{
			Err:       errors.New("must be > 0"),
			Section:   "[Logging]",
			Parameter: "LogQueueSize",
		}
	}

	return nil
}

func (conf *RustymonConfig) GetListenString() string {
	return conf.Server.ListenAddress + ":" + fmt.Sprint(conf.Server.ListenPort)
}
