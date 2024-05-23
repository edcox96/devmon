package http_server

import (
	_ "fmt"
	"log"
	"net/http"
	"os"
)

func NewHTTPServer() error {
	log.Printf("NewHttpServer")
	host, ok := os.LookupEnv("DEVMON_SERVER")
	if !ok || len(host) == 0 {
		host = "localhost:8080"
	}
	httpSrv.host = host
	
	httpSrv.mux = newServerMux()
	httpSrv.server = nil

	return nil
}

func newServerMux() http.ServeMux {

	
}

type HttpServer struct {
	host    string;
	server *http.Server;
    mux     http.ServeMux;
}

var httpSrv = &HttpServer {}

