package server

import (
	"github.com/phthaocse/go-gin-demo/utils"
)

type Config struct {
	ServerPort string
	DbDriver   string
	DbAddr     string
	DbName     string
	DbUser     string
	DbPassword string
}

func GetSrvConfig() *Config {
	port := utils.GetEnv("SERVER_PORT", "8000")
	port = ":" + port

	config := Config{
		ServerPort: port,
		DbDriver:   "postgres",
		DbAddr:     utils.GetEnv("DB_ADDRESS", "localhost"),
		DbName:     utils.GetEnv("DB_NAME", "go_gin_demo"),
		DbUser:     utils.GetEnv("DB_USER", "postgres"),
		DbPassword: utils.GetEnv("DB_PASSWORD", ""),
	}
	return &config
}
