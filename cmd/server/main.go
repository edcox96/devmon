package main

import (
    _ "fmt"
    "log"

    "github.com/edcox96/devmon/internal/server/mvc"

    httpsrv "github.com/edcox96/devmon/internal/server/http_server"
)

func main() {
    log.Printf("server main\n")

    // Create the mvc Controller which will also create the Model and Views
    if err := mvc.NewController(); err != nil {
        log.Fatalf("NewCotroller failed, err %s\n", err)
    }

    // Create the http_server
    err := httpsrv.NewHTTPServer()
    if err != nil {
        log.Fatalf("NewHTTPServer failed, err %s", err)
    }

    // MVC and http_server ready, Start http processing requests
    httpsrv.StartServer()
}

type DevMonState struct {
}

//var devMonState = DevMonState{}
