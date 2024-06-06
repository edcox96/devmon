package grpc_server

import (
	"context"
	"log"

    model "github.com/edcox96/devmon/internal/mvc/model"
    api "github.com/edcox96/devmon/api/v1"
	"google.golang.org/grpc"
)

func NewGrpcServer() (*UsbGrpcServer, error) {
	gsrv := grpc.NewServer()
	m := model.NewUsbModel()

	usbSrv, err := newGrpcServer(m)
	if err != nil {
		log.Printf("newGrpcSeerver failed. %s", err)
		return nil, err
	}
   	api.RegisterUsbServer(gsrv, usbSrv)

	return usbSrv, nil
}

func newGrpcServer(m *model.UsbModel) (*UsbGrpcServer, error) {
	srv := &UsbGrpcServer { model: m }
	return srv, nil
}

//var _ api.UsbServer = (*UsbGrpcServer) (nil)

type UsbGrpcServer struct {
	api.UnimplementedUsbServer
	model *model.UsbModel
}

func (s *UsbGrpcServer) PutUsbDevDesc(ctx context.Context,
				req *api.PutUsbDevDescRequest) (*api.PutUsbDevDescResponse, error) {
	log.Printf("PusUsbDevDesc")
	resp := &api.PutUsbDevConnResponse{}
	log.Printf("resp %s", resp)
	return nil, nil
}

func (s *UsbGrpcServer) PutUsbDevConnState(ctx context.Context,
				req *api.PutUsbDevConnRequest) (*api.PutUsbDevConnResponse, error) {
	log.Printf("PusUsbDevConnState")
	resp := &api.PutUsbDevConnResponse{}
	log.Printf("resp %s", resp)
	return nil, nil
}
