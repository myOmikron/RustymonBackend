package app

import (
	"errors"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/handler"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/execution"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"github.com/pelletier/go-toml"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

var asciiArt = `
 ______
|  ___ \            _                           
| |___) |_   _  ___| |_ _   _ ____   ___  ____  
|  __  /| | | |/___)  _) | | |    \ / _ \|  _ \ 
| |  \ \\ |_| |___ | |_| |_| | | | | |_| | | | |
|_|   \_|\____(___/ \___)__  |_|_|_|\___/|_| |_|
                        (____/   & a bunch of other languages`

func StartServer(configPath string) {
	color.Println(color.BLUE, asciiArt)
	fmt.Println()

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

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Set debug level
	e.Logger.SetLevel(log.DEBUG)

	// Initialize DB
	var driver gorm.Dialector

	switch config.Database.Driver {
	case "sqlite":
		driver = sqlite.Open(config.Database.Name)
	case "mysql":
		mysqlConf := mysqlDriver.NewConfig()
		mysqlConf.Net = fmt.Sprintf("tcp(%s)", net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))))
		mysqlConf.DBName = config.Database.Name
		mysqlConf.User = config.Database.User
		mysqlConf.Passwd = config.Database.Password
		mysqlConf.ParseTime = true
		mysqlConf.Params = map[string]string{
			"charset": "utf8mb4",
		}
		driver = mysql.Open(mysqlConf.FormatDSN())
	case "postgresql":
		dsn := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(config.Database.User, config.Database.Password),
			Host:   net.JoinHostPort(config.Database.Host, strconv.Itoa(int(config.Database.Port))),
			Path:   config.Database.Name,
		}
		driver = postgres.Open(dsn.String())
	}
	db := database.Initialize(
		driver,
		&utilitymodels.Session{},

		// Static
		&models.Modifier{},

		&models.Move{},
		&models.Item{},
		&models.Pokemon{},
		&models.SpawnArea{},

		&models.WeatherType{},
		&models.MoonType{},
		&models.TimeType{},

		&models.HeldItemCondition{},
		&models.Condition{},
		&models.PokemonSpawnRelation{},

		// Player specific
		&models.PokedexEntry{},
		&models.PlayerItem{},
		&models.PlayerPokemonMove{},
		&models.PlayerPokemon{},

		&models.Player{},

		// Account specific
		&models.PasswordReset{},
	)

	// Insert Pokémon up to ID 809
	fmt.Print("Populating pokemon\t\t...\t")
	for i := uint(1); i < 810; i++ {
		db.FirstOrCreate(&models.Pokemon{ID: i})
	}
	color.Println(color.GREEN, "done")

	// Insert Moves up to ID 676
	fmt.Print("Train pokemon to learn moves\t...\t")
	for i := uint(1); i < 678; i++ {
		db.FirstOrCreate(&models.Move{ID: i})
	}
	color.Println(color.GREEN, "done")

	// Insert Items up to ID 633
	fmt.Print("Buying a bunch of items\t\t...\t")
	for i := uint(1); i < 634; i++ {
		db.FirstOrCreate(&models.Item{ID: i})
	}
	color.Println(color.GREEN, "done")

	// WorkerPool
	poolConf := &worker.PoolConfig{
		NumWorker: 8,
		QueueSize: 80,
	}
	wp := worker.NewPool(poolConf)
	wp.Start()

	// Middleware
	e.Use(middleware.CustomContext(&handler.Context{}))
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format:           "",
		CustomTimeFormat: time.RFC1123Z,
		Output:           os.Stdout,
	}))
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Gzip())
	f := false
	age := time.Hour * 24
	e.Use(middleware.Session(
		db,
		&middleware.SessionConfig{
			Secure:         &f,
			CookieAge:      &age,
			DisableLogging: true,
		},
	))

	// Template rendering
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gohtml")),
	}
	e.Renderer = renderer

	// Routes
	defineRoutes(e, config, db, wp)

	fmt.Printf("\nListening on %s\n", color.Colorize(color.PURPLE, config.GetListenString()))
	execution.SignalStart(e, config.GetListenString(), func() {
		StartServer(configPath)
	})
}
