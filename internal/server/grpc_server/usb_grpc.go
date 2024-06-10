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

func (s *UsbGrpcServer) RegisterConsole(ctx context.Context,
	req *api.RegisterConsoleRequest) (*api.RegisterConsoleResponse, error) {
	log.Printf("RegisterConsole: ")

	s.model.RegisterConsole(req)
	resp := &api.RegisterConsoleResponse{}
	log.Printf("resp %v", resp)

	return resp, nil
}

func (s *UsbGrpcServer) RegisterUsbDevice(ctx context.Context,
	req *api.RegisterUsbDeviceRequest) (*api.RegisterUsbDeviceResponse,
	error) {
	log.Printf("RegisterDevice: ")

	con := s.model.Consoles[req.ConsID - 1]
	con.RegisterUsbDevice(req)
	resp := &api.RegisterUsbDeviceResponse{}
	log.Printf("resp %v", resp)

	return resp, nil
}

func (s *UsbGrpcServer) PutUsbDevDesc(ctx context.Context,
	req *api.PutUsbDevDescRequest) (*api.PutUsbDevDescResponse, error) {
	log.Printf("PusUsbDevDesc: Spec %x", req.Spec.BcdValue)
	con := s.model.Consoles[0]

	con.UsbDevs[0].PutUsbDevDesc(req)

	resp := &api.PutUsbDevDescResponse{}
	log.Printf("resp %v", resp)
	return resp, nil
}

func (s *UsbGrpcServer) PutUsbDevConnState(ctx context.Context,
	req *api.PutUsbDevConnRequest) (*api.PutUsbDevConnResponse, error) {
	log.Printf("PusUsbDevConnState")
	con := s.model.Consoles[0]

	con.UsbDevs[0].PutUsbDevConnState(req)

	resp := &api.PutUsbDevConnResponse{}
	log.Printf("resp %v", resp)
	return nil, nil
}
