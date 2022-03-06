package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
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
	teardownFunc := func() {
		db.Close()
	}
	return db, teardownFunc, nil
}
