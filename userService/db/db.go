package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "failed opening database")
	}

	if err = connection.db.Ping(); err != nil {
		return errors.Wrap(err, "failed veryfing database connection")
	}

	return nil
}
