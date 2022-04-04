package app

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	log2 "github.com/labstack/gommon/log"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/execution"
	"github.com/myOmikron/echotools/logging"
	"github.com/myOmikron/echotools/worker"
	"github.com/pelletier/go-toml"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
)

var log logging.Logger

var asciiArt = `
 ______
|  ___ \            _                           
| |___) |_   _  ___| |_ _   _ ____   ___  ____  
|  __  /| | | |/___)  _) | | |    \ / _ \|  _ \ 
| |  \ \\ |_| |___ | |_| |_| | | | | |_| | | | |
|_|   \_|\____(___/ \___)__  |_|_|_|\___/|_| |_|
                        (____/   & a bunch of other languages`

func StartServer(configPath string, isReloading bool) {
	config := &configs.RustymonConfig{}

	if configBytes, err := ioutil.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(config)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, config); err != nil {
			panic(err)
		}
	}

	// Check for valid config values
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Logging
	logging.Initialize(&logging.Config{
		LogQueue:      config.Logging.LogQueueSize,
		LogFile:       config.Logging.LogFile,
		LogMaxDays:    config.Logging.LogMaxDays,
		LogMaxBackups: config.Logging.LogMaxBackups,
		LogMaxSize:    config.Logging.LogMaxCapacity,
	})

	log = logging.GetLogger("app")

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	// Disable default logger as it just sucks
	e.Logger.SetLevel(log2.OFF)

	color.Println(color.BLUE, asciiArt)
	fmt.Println()

	// Initialize DB
	db := InitializeDatabase(config)

	// WorkerPool
	poolConf := &worker.PoolConfig{
		NumWorker: 8,
		QueueSize: 80,
	}
	wp := worker.NewPool(poolConf)
	wp.Start()

	// Middleware
	InitializeMiddleware(e, db, config)

	// Template rendering
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gohtml")),
	}
	e.Renderer = renderer

	// Routes
	defineRoutes(e, config, db, wp)

	// Start RPC listener on unix socket
	var cliSock net.Listener
	go InitializeRPC(&cliSock, config.Server.CLIUnixSocket, db, isReloading)

	log.Infof("Listening on %s", color.Colorize(color.PURPLE, config.GetListenString()))
	execution.SignalStart(e, config.GetListenString(), &execution.Config{
		ReloadFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
			StartServer(configPath, true)
		},
		StopFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
			logging.Stop()
		},
		TerminateFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
			logging.Stop()
		},
	})
}
