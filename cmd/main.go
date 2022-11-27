package main

import (
	"github.com/qwuiemme/ellipsespace-server/internal/server"
)

func main() {
	server := new(server.Server)
	server.Run("localhost:8888")
}
