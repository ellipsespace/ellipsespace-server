package main

import (
	"github.com/qwuiemme/ellipsespace-server/internal/server"
)

// @Title EllipseSpace API
// @version 0.0
// @description API for the Encyclopedia of Space project

// @host localhost:8888
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	server := new(server.Server)
	server.Run("localhost:8888")
}
