package http_server

import (
	_ "fmt"
	"log"
	"net/http"
	"os"
)

func NewHTTPServer() error {
	log.Printf("NewHttpServer")
	httpSrv.address, ok := os.LookupEnv("DEVMON_SERVER")
	if !ok || len(addr) == 0 {
		httpSrv.address := "localhost:8080"
	}
	return nil
}

type httpServer struct {
	address string;
	server *http.Server;
    mux     http.ServeMux;
}

var httpSrv = &httpServer {}

