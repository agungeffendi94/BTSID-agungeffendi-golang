package core

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectDB() *sqlx.DB {
	// //------------ initial koneksi postgre ---------------
	username := viper.GetString("database.user")
	password := viper.GetString("database.pass")
	database := viper.GetString("database.name")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error Connecting DB => ", err)
		os.Exit(0)
	}

	maxLifetime, _ := time.ParseDuration(viper.GetString("database.max_lifetime_connection") + "s")
	db.SetMaxIdleConns(viper.GetInt("database.max_idle_dbection"))
	db.SetConnMaxLifetime(maxLifetime)
	dbs := db

	return dbs
}
