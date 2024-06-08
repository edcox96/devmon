package grpc_server

import (
	"context"
	"log"
	"net"

	api "github.com/edcox96/devmon/api/v1"
	model "github.com/edcox96/devmon/internal/server/mvc/model"
	"google.golang.org/grpc"
)

func NewUsbGrpcServer(address string) (*UsbGrpcServer, error) {
	log.Printf("NewUsbGrpcServer")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("net.listen failed: %s", err)
		return nil, err
	}

	gsrv := grpc.NewServer()

	usbGrpcSrv.gsrv = gsrv
	usbGrpcSrv.lis = lis
	usbGrpcSrv.model = model.NewUsbModel()

	api.RegisterUsbServer(gsrv, usbGrpcSrv)

	return usbGrpcSrv, nil
}

func StartServer(srv *UsbGrpcServer) error {
	if err := srv.gsrv.Serve(srv.lis); err != nil {
		log.Printf("gsrv.Serve failed! %s", err)
		return err
	}
	return nil
}

type UsbGrpcServer struct {
	api.UnimplementedUsbServer
	gsrv  *grpc.Server
	lis   net.Listener
	model *model.UsbModel
}

var usbGrpcSrv = &UsbGrpcServer{}

func (s *UsbGrpcServer) PutUsbDevDesc(ctx context.Context,
	req *api.PutUsbDevDescRequest) (*api.PutUsbDevDescResponse, error) {
	log.Printf("PusUsbDevDesc: Spec %x", req.Spec.BcdValue)

	s.model.PutUsbDevDesc(req)
	
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
