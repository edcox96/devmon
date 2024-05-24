package main

import (
	_ "fmt"
	"log"

	http_srv "github.com/edcox96/devmon/server/internal/http_server"
	mvc "github.com/edcox96/devmon/server/internal/mvc"
)

func main() {
	log.Printf("server main\n")

	// Create the mvc to init Controller and Views
	if err := mvc.NewController(); err != nil {
		log.Fatalf("NewCotroller failed, err %s\n", err)
	}

	// Create the http_server
	err := http_srv.NewHTTPServer()
	if err != nil {
		log.Fatalf("NewHTTPServer failed, err %s", err)
	}

	// MVC and http_server ready, Start http processing requests
	http_srv.StartServer()
}

type DevMonState struct {
}

//var devMonState = DevMonState{}
