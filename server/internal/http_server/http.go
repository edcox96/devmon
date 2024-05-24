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
	
	if err := httpSrv.newServerMux(); err != nil {
		log.Printf("httpSrv.newServeMux failed, err %s\n", err)
		return err
	}

	if err := httpSrv.newHttpServer(); err != nil {
		log.Printf("httpSrv.newServeMux failed, err %s\n", err)
		return err
	}

	return nil
}

func (hs *HttpServer) newServerMux() error {
	hs.mux = http.NewServeMux()

	// add mux handlers here

	return nil;
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

