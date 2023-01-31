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
		viper.GetString("userDB.host"),
		viper.GetInt("userDB.port"),
		viper.GetString("userDB.user"),
		viper.GetString("userDB.dbname"),
		viper.GetString("userDB.password"),
		viper.GetString("userDB.sslmode"))

	if connection.db, err = sql.Open("postgres", dsn); err != nil {
		return err // TODO: где враппинг ошибок?
	}

	if err = connection.db.Ping(); err != nil {
		return err // TODO: где враппинг ошибок?
	}

	return nil
}
