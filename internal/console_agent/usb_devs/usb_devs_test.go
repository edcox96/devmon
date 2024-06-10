package console_agent_test

import (
	"flag"
	"log"
	"os"
	"testing"

	cons_agent "github.com/edcox96/devmon/internal/console_agent"
	devsim "github.com/edcox96/devmon/internal/devsim"
	usb_devs "github.com/edcox96/devmon/internal/console_agent/usb_devs"
	grpc_srv "github.com/edcox96/devmon/internal/server/grpc_server"
	"github.com/google/gousb"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.Printf("TestMain running tests\n")

	address := "localhost:8080"

	usbGrpcSrv, err := grpc_srv.NewUsbGrpcServer(address)
	if err != nil {
		log.Fatalf("setupTestServer failed")
	}

	go func() {
		log.Printf("server go routine: Serve(lis)")
		grpc_srv.StartServer(usbGrpcSrv)
		log.Printf("server go routine: after StartServer)")
	}()

	usbClient, err := cons_agent.NewGrpcUsbClient(address)
	if err != nil {
		log.Fatalf("NewGrpcUsbClient failed! %s", err)
	}

	// init the devsim consoles
	err = devsim.InitDevSim()
	if err != nil {
		log.Fatalf("devsim.InitDevSim failed! %v", err)
	}

	err = usb_devs.SendUsbInfoToServer(usbClient, 1)
	if err != nil {
		log.Fatalf("SendUsbInfoToServer failed, %s", err)
	}

	os.Exit(m.Run())
}

func TestOpenUsbDev(t *testing.T) {
	// hack to just get known device vid/pid for camera
	vid, pid := gousb.ID(0x04b4), gousb.ID(0x00f9)
	dev, err := usb_devs.OpenUsbDev(vid, pid)
	require.NoError(t, err)
	log.Printf("Open camera dev %v", dev)
	//require.Equal(t, vid, dev
}
