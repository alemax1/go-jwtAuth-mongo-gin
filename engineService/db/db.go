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

	if connection.db, err = sql.Open("postgres", fmt.Sprintf(`host=%s port=%v user=%s dbname=%s password=%s sslmode=%s`,
		viper.GetString("engineDB.host"),
		viper.GetInt("engineDB.port"),
		viper.GetString("engineDB.user"),
		viper.GetString("engineDB.dbname"),
		viper.GetString("engineDB.password"),
		viper.GetString("engineDB.sslmode"))); err != nil {
		return err
	}

	if err = connection.db.Ping(); err != nil {
		return err
	}

	return nil
}
