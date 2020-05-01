package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

var db *gorm.DB

type DBConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func NewDBConfig() DBConfig {
	return DBConfig{Host: "127.0.0.1", Port: "5432", Username: "postgres"}
}

// InitDB : Initialize DB client
func InitDB(cfg *DBConfig) error {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Host, cfg.Port,
		cfg.Username, cfg.DatabaseName, cfg.Password,
	)

	dbClient, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	db = dbClient

	if viper.GetString("GO_ENV") == "development" {
		dbClient.LogMode(viper.GetBool("database.debug_log"))
	}

	return nil
}

func Client() *gorm.DB {
	if db == nil {
		panic("Database is not initialized")
	}
	return db
}

func CloseDB() error {
	return db.Close()
}
