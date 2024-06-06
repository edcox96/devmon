package http_server

import (
    "context"
    "log"
    _ "net/http"
    "os"
    "testing"

    _ "github.com/stretchr/testify/require"
    "github.com/edcox96/devmon/internal/mvc"
    "github.com/edcox96/devmon/client/http_clients"
)

func TestMain(m *testing.M) {
    log.Printf("TestMain running tests\n")

    if err := mvc.NewController(); err != nil {
        log.Fatalf("NewCotroller failed, err %s\n", err)
    }

    os.Exit(m.Run())
}

func TestHttpServer(t *testing.T) {
    setupTest(t)
}

func setupTest(t *testing.T) () {
    t.Helper()

    if err := NewHTTPServer(); err != nil {
        log.Fatalf("NewHTTPServer err %s", err)
    }
    begin := make(chan string)

    // start a mock client app to interact with the server
    go clientAgent(begin)

    begin<- "proceed" // tell client go routine to proceed sending HEAD requests 
    StartServer()
}

// call the http_client to run REST API tests and then reset the server
func clientAgent(begin <-chan string) {
    <-begin

    http_client.ClientRestAPITests()

    log.Printf("gracefully shutdown server\n")
    httpSrv.server.Shutdown(context.Background())
}
