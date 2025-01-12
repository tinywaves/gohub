package main

import "gohub/internal"

func main() {
	server := internal.InitWeb()
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
