package app

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/RustymonBackend/models"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/database"
	"github.com/myOmikron/echotools/utilitymodels"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"net/url"
	"strconv"
)

func InitializeDatabase(config *configs.RustymonConfig) (db *gorm.DB) {
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
	db = database.Initialize(
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

	return
}
