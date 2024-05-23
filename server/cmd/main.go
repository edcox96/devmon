package main

import (
	_ "fmt"
	"log"

	http_srv "github.com/edcox96/devmon/server/internl/http_server"
	mvc "github.com/edcox96/devmon/server/internal/mvc"
	store "github.com/edcox96/devmon/server/internal/storage"
)

func main() {
	log.Printf("server main\n")

	// Make sure the storage is ready for transactions
	if err := store.InitStorage(); err != nil {
		log.Fatalf("InitStorage failed, err %s\n", err)
	}

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
	if err := http_srv.StartServer(); err != nil {
		// for now just treat this as fatal, later restart
		log.Fatalf("Exiting: HTTP Seerver failed! err %s\n", err)
	}
}

type DevMonState struct {
}

var devMonState = DevMonState {}
