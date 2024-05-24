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
	
	httpSrv.mux = httpSrv.newServerMux()
	httpSrv.server = nil

	return nil
}

func (hs *HttpServer) newServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	// add mux handlers

	return mux;
}

func (hs *HttpServer) newHttpServer() error {
	hs.server = &http.Server {
		Addr:    hs.host,
		Handler: hs.mux,
	}
	return nil
}

func StartServer() {
	if httpSrv.server == nil {
	   log.Fatalf("httpSrv.server is nil")
	}
	log.Fatal(httpSrv.server.ListenAndServe())
}

type HttpServer struct {
	host    string;
	server *http.Server;
    mux    *http.ServeMux;
}

var httpSrv = &HttpServer {}

