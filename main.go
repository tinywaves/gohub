package main

import "gohub/internal/web"

func main() {
	server := web.InitWeb()
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
