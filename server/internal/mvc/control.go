package mvc

import (
    "context"
    "errors"
    "fmt"
    "log"
    "net/http"
    _ "strings"
    _ "sync"
)

var ErrorModelCreationFailed = errors.New("unable to create model")
var ErrorViewsCreationFailed = errors.New("unable to create views")

func NewController() error {
    log.Printf("NewController")

    // Create the Model to handle request access to storage
    if err := NewModel(); err != nil {
        return ErrorModelCreationFailed
    }

    // Create the Views to handle request view generation
    if err := NewViews(); err != nil {
        return ErrorModelCreationFailed
    }
    return nil
}

type MvcController struct {
    //httpReqProc []HttpRequestProcessor
}

type RequestProcessor interface {
    Run()
}

type HttpReqProc struct {
    ReqProcNum int
    Requests   chan *HttpMvcRequest
    Errors     chan error
    Response   chan string
}

func (rp *HttpReqProc) Run() {
    requests := make(chan *HttpMvcRequest, 1)
    rp.Requests = requests

    errors := make(chan error, 1)
    rp.Errors = errors

    /* temp until controller sends response */
    rp.Response = make(chan string, 1)

    go func() {
        defer close(requests)
        defer close(errors)
        defer close(rp.Response) //temp - move to more appropriate place

        var reqNum int

        for r := range requests {
            printCtrlReq(r)
            if err := rp.validateRequest(r); err != nil {
                log.Printf("Control.Run: invalid request")
                errors <- err
                continue
            }
            resp, err := rp.Process(r, reqNum)
            if err != nil {
                log.Printf("rp.Run: errror %s", err)
                rp.Errors <- err
            }
            // send response to http_server
            rp.Response <- resp
            reqNum++
        }
    }()
}

func (rp *HttpReqProc) validateRequest(_ *HttpMvcRequest) error {
    return nil // decide what validation of the Control mvcRequest is needed if any
}

func printCtrlReq(ctrlReq *HttpMvcRequest) {
    rt := ctrlReq.ReqType
    var cons string
    if len(ctrlReq.ConsoleNum) > 0 {
        cons = "console"
    }
    log.Printf("ctrl.Run: %s %s %s %s %s %s %s",
                 ctrlReq.Method, rt.String(rt), cons, ctrlReq.ConsoleNum, ctrlReq.DeviceType, 
                 ctrlReq.DeviceNum, ctrlReq.UsbObject)
}
func (rp *HttpReqProc) Process(ctrlReq *HttpMvcRequest, reqNum int) (string, error) {
    response := fmt.Sprintf("Now serving request: %d.\n", reqNum)
    return response, nil
}

type RequestType byte

const (
    _				RequestType = iota
    GetConsoleList  RequestType = iota
    GetConsoleInfo
    GetDeviceList
    GetDeviceInfo
    GetDeviceObject
)

func (rt *RequestType) String(reqType RequestType) string {
    switch reqType {
    case GetConsoleList: return "ConsList"
    case GetConsoleInfo: return "ConsInfo"
    case GetDeviceList:  return "DevList"
    case GetDeviceInfo:  return "DevInfo"
    case GetDeviceObject: return "DevObj"
    }
    return ""
}

type HttpMvcRequest struct {
    Ctx        context.Context
    HttpReq    *http.Request
    ReqType    RequestType
    // http request parameters
    Method     string
    ConsoleNum string
    DeviceType string
    DeviceNum  string
    UsbObject  string
}
