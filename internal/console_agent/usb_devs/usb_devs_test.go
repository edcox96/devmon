package console_agent_test

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"testing"

	api "github.com/edcox96/devmon/api/v1"
	usb_devs "github.com/edcox96/devmon/internal/console_agent/usb_devs"
	"github.com/google/gousb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	log.Printf("TestMain running tests\n")

	flag.Parse()
	log.Printf("console_agent\n")

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("net.listen failed: %s", err)
	}

	gsrv, err := setupTestServer()
	if err != nil {
		log.Fatalf("setupTestServer failed")
	}
	//defer lis.Close()

	usbClient, err := usb_devs.Connect(lis)
	if err != nil {
		log.Fatalf("Connect failed! %s", err)
	}

	err = usb_devs.SendUsbInfoToServer(usbClient)
	if err != nil {
		log.Fatalf("SendUsbInfoToServer failed, %s", err)
	}

	go func() {
		log.Printf("server go routine: Serve(lis)")
		gsrv.Serve(lis)
		log.Printf("server go routine: after Serve(lis)")
	}()

	os.Exit(m.Run())
}

func TestOpenUsbDev(t *testing.T) {
	// hack to just get known device vid/pid for camera
	vid, pid := gousb.ID(0x04b4), gousb.ID(0x00f3)
	dev, err := usb_devs.OpenUsbDev(vid, pid)
	require.NoError(t, err)
	log.Printf("Open camera dev %v", dev)
	//require.Equal(t, vid, dev
}

func setupTestServer() (*grpc.Server, error) {
	log.Printf("setupTestServer")
	//_, err := net.Listen("tcp", ":0")
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}

	gsrv := grpc.NewServer()
	usbSrv, err := newGrpcServer(&usbModel)
	if err != nil {
		log.Printf("newGrpcSeerver failed. %s", err)
		return nil, err
	}
	api.RegisterUsbServer(gsrv, usbSrv)

	return gsrv, nil
}

func newGrpcServer(m *UsbModel) (*UsbGrpcServer, error) {
	srv := &UsbGrpcServer{model: m}
	return srv, nil
}

type UsbModel struct {
	usbData int
}

var usbModel = UsbModel{usbData: 10}

type UsbGrpcServer struct {
	api.UnimplementedUsbServer
	model *UsbModel
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
