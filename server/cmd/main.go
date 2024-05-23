package main

import (
	"fmt"
	"log"

	"devmon/server/http_server"
)

func main() {
	fmt.Printf("server main\n")
	fmt.Println("Starting server")
	srv := http_server.NewHTTPServer(":8080")
	if srv == nil {
		fmt.Println("Server is nil, returning")
		return
	}
	log.Fatal(srv.ListenAndServe())
}

