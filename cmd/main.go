package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/qwuiemme/ellipsespace-server/internal/server"
)

// @Title EllipseSpace API
// @version 1.3.0
// @description API for the Encyclopedia of Space project

// @host ellipsespace.onrender.com
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env files found")
	}
}

func main() {
	server := new(server.Server)
	addr := server.MakeAddr()
	server.Run(addr)
}
