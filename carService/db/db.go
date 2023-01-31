package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type database struct {
	db *sql.DB
}

var connection database

func CreateConnection() error {
	var err error

	dsn := fmt.Sprintf(`host=%s port=%v user=%s dbname=%s password=%s sslmode=%s`,
		viper.GetString("carDB.host"),
		viper.GetInt("carDB.port"),
		viper.GetString("carDB.user"),
		viper.GetString("carDB.dbname"),
		viper.GetString("carDB.password"),
		viper.GetString("carDB.sslmode"))

	if connection.db, err = sql.Open("postgres", dsn); err != nil {
		return err // где враппинг ошибки?
	}

	if err = connection.db.Ping(); err != nil {
		return err // где враппинг ошибки?
	}

	return nil
}
