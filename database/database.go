package database

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Migrate() (err error) {
	connStr, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		return errors.New("connection string Not Found")
	}
	m, err := migrate.New(
		"file://database/migrations",
		connStr)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Println(err.Error())
	}

	return err
}

func GetDb(w http.ResponseWriter) (db *sqlx.DB, err error) {
	if DB != nil {
		err = DB.Ping()
		if err == nil {
			db = DB
			return db, nil
		}

	}
	log.Println("Opening db")
	connStr, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		return nil, errors.New("connection string Not Found")
	}
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	DB = db

	return db, err
}
