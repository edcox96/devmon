package http_server

import (
    "context"
    _ "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/edcox96/devmon/server/internal/mvc"
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

func (h *HttpServer) newServerMux() error {
    h.mux = http.NewServeMux()
    if h.mux == nil {
        log.Fatalf("http.NewServeMux returned nil")
    }
    //h.mux.HandleFunc("GET /", httpSrv.getRootHandler)
    h.mux.HandleFunc("GET /index.html", httpSrv.getIndexHandler)
    h.mux.HandleFunc("GET /devmon/v1/consoles/", httpSrv.getConsoleListHandler)
    h.mux.HandleFunc("GET /devmon/v1/console/{consNum}/", httpSrv.getConsoleHandler)
    h.mux.HandleFunc("GET /devmon/v1/console/{consNum}/devices",
                     httpSrv.getDeviceListHandler)
    h.mux.HandleFunc("GET /devmon/v1/console/{consNum}/{devType}/{devNum}/",
                     httpSrv.getDeviceHandler)
    h.mux.HandleFunc("GET /devmon/v1/console/{consNum}/{devType}/{devNum}/{usbObject}/",
                     httpSrv.getDeviceObjectHandler)
    h.mux.HandleFunc("GET /devmon/v1/eventlog/", httpSrv.getUsbEventLogHandler)
    h.mux.HandleFunc("GET /devmon/v1/swlog/", httpSrv.getSoftwareLogHandler)
                 
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
    // start the MVC engine to process http requests
    reqProc = &mvc.HttpReqProc { ReqProcNum: 1 }
    //httpSrv.reqProc[0] = reqProc
    reqProc.Run()

    if httpSrv.server == nil {
       log.Fatalf("httpSrv.server is nil")
    }
    err := httpSrv.server.ListenAndServe()
    if err != http.ErrServerClosed {
        log.Fatalf("ListenAndServe err %s\n", err)
    }
}

type HttpServer struct {
    host    string;
    server *http.Server;
    mux    *http.ServeMux;
    //reqProc []*mvc.HttpRequestProcessor
}

var httpSrv = &HttpServer {}
var reqProc *mvc.HttpReqProc
//var devmonApiVersion = "v1"

func (h *HttpServer) getIndexHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getIndexHandler: url path %s", r.URL.Path)
    var statusCode = http.StatusOK
    // TBD - figure out request to get the index.html and set Content type to html
    // and return the page to client.
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    w.WriteHeader(statusCode)
}

func (h *HttpServer) getConsoleListHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getConsoleListHandler: url path %s\n", r.URL.Path)
    var ctrlReq = newControlGetHttpRequest(r)
    ctrlReq.ReqType = mvc.GetConsoleList
    h.ServeCtrlRequest(w, ctrlReq)
}

func (h *HttpServer) getConsoleHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getConsoleHandler: url path %s\n", r.URL.Path)

    var ctrlReq = newControlGetHttpRequest(r)
    ctrlReq.ConsoleNum = r.PathValue("consNum")
    ctrlReq.ReqType = mvc.GetConsoleInfo

    h.ServeCtrlRequest(w, ctrlReq)
}

func (h *HttpServer) getDeviceListHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getDeviceListHandler: url path %s\n", r.URL.Path)

    var ctrlReq = newControlGetHttpRequest(r)
    ctrlReq.ConsoleNum = r.PathValue("consNum")
    ctrlReq.ReqType = mvc.GetDeviceList

    h.ServeCtrlRequest(w, ctrlReq)
}

func (h *HttpServer) getDeviceHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getDeviceHandler: url path %s\n", r.URL.Path)

    var ctrlReq = newControlGetHttpRequest(r)
    ctrlReq.ConsoleNum = r.PathValue("consNum")
    ctrlReq.DeviceType = r.PathValue("devType")
    ctrlReq.DeviceNum  = r.PathValue("devNum")
    ctrlReq.ReqType = mvc.GetDeviceInfo

    h.ServeCtrlRequest(w, ctrlReq)
}

var validUsbObjs = []string { "device_desc", "config_desc" }

func (h *HttpServer) getDeviceObjectHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getDeviceObjectHandler: url path %s\n", r.URL.Path)

    var v bool
    var usbObject = r.PathValue("usbObject")
    
    // check list of valid usbObjects
    for _, obj := range validUsbObjs {
        if strings.Compare(usbObject, obj) == 0 {
            v = true
            break
        }
    }
    // not valid usbObject then return bad request
    if !v {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // create new mvcRequest and send to request processor
    var ctrlReq = newControlGetHttpRequest(r)
    ctrlReq.ConsoleNum = r.PathValue("consNum")
    ctrlReq.DeviceType = r.PathValue("devType")
    ctrlReq.DeviceNum  = r.PathValue("devNum")
    ctrlReq.UsbObject  = r.PathValue("usbObject")
    ctrlReq.ReqType = mvc.GetDeviceObject

    h.ServeCtrlRequest(w, ctrlReq)
}

func (h *HttpServer) getUsbEventLogHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getEventLogHandler: url path %s\n", r.URL.Path)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    w.WriteHeader(http.StatusOK)
}

func (h *HttpServer) getSoftwareLogHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getSoftwareLogHandler: url path %s\n", r.URL.Path)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    w.WriteHeader(http.StatusOK)
}

func newControlGetHttpRequest(r *http.Request) *mvc.HttpMvcRequest {
    var ctrlReq = &mvc.HttpMvcRequest { Method: "GET" }
    ctrlReq.Ctx = context.Background()
    ctrlReq.HttpReq = r
    return ctrlReq
}

func (h *HttpServer) ServeCtrlRequest(w http.ResponseWriter,
                                      ctrlReq *mvc.HttpMvcRequest) {
    //log.Printf("newControlGetHttpRequest: ctrlReq %v", ctrlReq)
    var sc int

    // send mvc request to cotroller request processor
    reqProc.Requests <- ctrlReq

    select { // wait for response, error, or done cancel ctx
    case err := <- reqProc.Errors:
        log.Printf("getHandler received an error, %s", err)
        sc = http.StatusNotFound
    case <- ctrlReq.Ctx.Done():
        log.Printf("getHandler: request canceled\n")
        sc = http.StatusNotFound
    case resp := <- reqProc.Response:
        log.Printf("getHandler: Response: %s", resp)
        sc = http.StatusOK
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    w.WriteHeader(sc)
}
