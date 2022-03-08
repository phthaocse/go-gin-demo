package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/db"
	"log"
)

type server struct {
	router     *gin.Engine
	db         *sql.DB
	dbTeardown func()
}

func createServer(config *Config) *server {
	svr := &server{}
	svr.router = gin.Default()
	setUpRouter(svr)
	dbCon, dbTeardown, err := db.SetUp(&db.DBConfig{config.DbDriver, config.DbAddr, config.DbName, config.DbUser, config.DbPassword})
	if err != nil {
		log.Panicln("Can't set up a database")
	}
	svr.db = dbCon
	svr.dbTeardown = dbTeardown
	return svr
}

func Start() {
	config := GetSrvConfig()
	svr := createServer(config)
	defer svr.dbTeardown()
	svr.router.Run(config.ServerPort)
}
