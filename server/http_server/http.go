package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"errors"

	"devmon/storage"
)

/* Function for creating http servers based on the http package at a specified address. */
func NewHTTPServer(addr string) *http.Server {
	// Create a new instance of our httpServer type and a new ServeMux to handle request multiplexing.
	httpsrv := newHttpServer()
	if httpsrv == nil {
		return nil
	}

	fmt.Println("Initializing transaction log")
	err := httpsrv.initializeTransactionLog()
	if err != nil {
		fmt.Println("Unable to initialize transaction log, returning")
		return nil
	}
	mux := http.NewServeMux()

	// Handle get and post requests at the root path of the server using the handler functions defined below.
	mux.HandleFunc("/", httpsrv.rootFunc)
	mux.HandleFunc("PUT /v1/key/{key}", httpsrv.keyValuePutHandler)
	mux.HandleFunc("GET /v1/key/{key}", httpsrv.keyValueGetHandler)
	mux.HandleFunc("DELETE /v1/key/{key}", httpsrv.keyValueDeleteHandler)

	// Create a server from the http package at the address passed using the defined ServeMux and return the server.
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &server
}

var logger TransactionLogger

/* Type httpServer used for defining our server. */
type httpServer struct {
	kvStore *storage.Store
}

/* Used for initializing the httpServer with a new log. */
func newHttpServer() *httpServer {
	return &httpServer{
		kvStore: storage.NewStore(),
	}
}

/* Request and Response structs. */
type Request struct {
	Value string `json:"value"`
}

type GetRequest struct {
	Key string `json:"key"`
}

type PutResponse struct {
	Added bool `json:"added"`
}

type GetResponse struct {
	Value string `json:"value"`
}

type DeleteResponse struct {
	Removed bool `json:"removed"`
}

func (s *httpServer) initializeTransactionLog() error {
    var err error

	fmt.Println("creating new file transaction logger")
    logger, err = NewFileTransactionLogger("transaction.log")
    if err != nil {
        return fmt.Errorf("failed to create event logger: %w", err)
    }

	fmt.Println("Reading events")
    events, errors := logger.ReadEvents()

    e := Event{}
    ok := true

    for ok && err == nil {
        select {
        case err, ok = <-errors:            // Retrieve any errors
        case e, ok = <-events:
            switch e.EventType {
            case EventDelete:               // Got a DELETE event!
                err = s.kvStore.Delete(e.Key)
            case EventPut:                  // Got a PUT event!
                err = s.kvStore.Put(e.Key, e.Value)
            }
        }
    }

    logger.Run()

	fmt.Println(err)
    return err
}

func (s *httpServer) rootFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rootFunc")
	if r.URL.Path != "/" {
		http.Error(w, "Wrong URL", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome to / handler")
}

/* Server handler for PUT requests adding key/value pairs to the server's key/value store. */
func (s *httpServer) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {

	// Get the key and value from the JSON request and handle any errors.
	key := r.PathValue("key")
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the key/value pair to the store and handle any errors.
	err = s.kvStore.Put(key, string(value))
	logger.WritePut(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the resulting PutResponse and write it, handling any errors.
	res := PutResponse{true}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* Server hadler for GET requests reading from the server's kvStore. */
func (s *httpServer) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {

	// Get the key from the JSON request, the value from the kvStore and handle any errors.
	key := r.PathValue("key")
	value, err := s.kvStore.Get(key)
	if errors.Is(err, storage.ErrorKeyNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the resulting GetResponse and write it, handling any errors.
	res := GetResponse{string(value)}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (s *httpServer) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the key from the JSON request, check if it is in the kvStore, and handle any errors.
	key := r.PathValue("key")
	_, err := s.kvStore.Get(key)
	if errors.Is(err, storage.ErrorKeyNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the key from the kvStore and handle any errors.
	err = s.kvStore.Delete(key)
	logger.WriteDelete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the resulting DeleteResponse and write it, handling any errors.
	res := DeleteResponse{Removed: true}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
