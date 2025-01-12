package main

import "gohub/internal/api"

func main() {
	server := api.Init()
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
