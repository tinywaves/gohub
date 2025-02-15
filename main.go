package main

import (
	"gohub/internal"
	"gohub/internal/api"
)

func main() {
	server := api.Init()
	err := server.Run(":" + string(rune(internal.Port)))
	if err != nil {
		return
	}
}
