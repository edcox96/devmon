package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/* Function for creating http servers based on the http package at a specified address. */
func NewHTTPServer(addr string) *http.Server {
	// Create a new instance of our httpServer type and a new ServeMux to handle request multiplexing.
	httpsrv := newHttpServer()
	if httpsrv == nil {
		return nil
	}
	mux := http.NewServeMux()

	// Handle get and post requests at the root path of the server using the handler functions defined below.
	mux.HandleFunc("/", httpsrv.rootFunc)
	mux.HandleFunc("GET /task", httpsrv.handleConsume)
	mux.HandleFunc("POST /task", httpsrv.handleProduce)

	// Create a server from the http package at the address passed using the defined ServeMux and return the server.
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &server
}

/* Type httpServer used for defining our server. */
type httpServer struct {
	Log *Log
}

/* Used for initializing the httpServer with a new log. */
func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

/* Request and Response structs. */
type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) rootFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rootFunc")
	if r.URL.Path != "/" {
		http.Error(w, "Wrong URL", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome to / handler")
}

/* Server handler for POST requests writing to the server's log. */
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleProduce")

	// Unmarshal the JSON request's body and handle any errors.
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Append the request record to the log and handle any errors.
	offset, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the resulting ProduceResponse and write to the response, handling any errors.
	res := ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* Server hadler for GET requests reading from the server's log. */
func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleConsume")

	// Unmarshal the JSON request's body and handle any errors.
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read the requested record from the log and handle any errors.
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the resulting ConsumeResponse and write to the response, handling any errors.
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
