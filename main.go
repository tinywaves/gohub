package main

import (
	"gohub/internal"
	"gohub/internal/api"
	"strconv"
)

func main() {
	server := api.Init()
	err := server.Run(":" + strconv.Itoa(internal.Port))
	if err != nil {
		return
	}
}
