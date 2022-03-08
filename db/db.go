package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DBConfig struct {
	DbDriver   string
	DbAddr     string
	DbName     string
	DbUser     string
	DbPassword string
}

func SetUp(config *DBConfig) (*sql.DB, func(), error) {
	var err error
	var db *sql.DB
	dbDsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", config.DbAddr, config.DbName, config.DbUser, config.DbPassword)
	db, err = sql.Open(config.DbDriver, dbDsn)
	if err != nil {
		log.Fatal("Unable to connect to db with err: ", err)
		return nil, nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to check the health of db due to: ", err)
		return nil, nil, err
	}
	log.Println("Connect to db successfully")
	curr_dir, _ := os.Getwd()
	migrateSrc := fmt.Sprintf("file://%s/db/migrations", curr_dir)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrateSrc,
		config.DbName, driver)
	m.Up()
	teardownFunc := func() {
		db.Close()
	}
	return db, teardownFunc, nil
}
