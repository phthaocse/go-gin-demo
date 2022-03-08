package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/phthaocse/go-gin-demo/db"
	"log"
)

type Server struct {
	Router     *gin.Engine
	db         *sql.DB
	DbTeardown func()
}

func CreateServer(config *Config) *Server {
	svr := &Server{}
	svr.Router = gin.Default()
	setUpRouter(svr)
	dbCon, dbTeardown, err := db.SetUp(&db.DBConfig{config.DbDriver, config.DbAddr, config.DbName, config.DbUser, config.DbPassword})
	if err != nil {
		log.Panicln("Can't set up a database")
	}
	svr.db = dbCon
	svr.DbTeardown = dbTeardown
	return svr
}

func Start() {
	config := GetSrvConfig()
	svr := CreateServer(config)
	svr.Router.Run(config.ServerPort)
}
