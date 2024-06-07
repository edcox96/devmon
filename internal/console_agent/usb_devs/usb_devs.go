package console_agent

import (
	"context"
	"fmt"
	"log"
	"net"

	api "github.com/edcox96/devmon/api/v1"
	"github.com/google/gousb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//    _ "github.com/edcox96/devmon/server/internal/mvc/model"

func Connect(lis net.Listener) (api.UsbClient, error) {
	ccOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	cc, err := grpc.NewClient(lis.Addr().String(), ccOpts...)
	if err != nil {
		return nil, err
	}

	client := api.NewUsbClient(cc)

	return client, nil
}

func OpenUsbDev(vid, pid gousb.ID) (*gousb.Device, error) {
	devs, err := GetUsbDevsWithVidPid(vid, pid)
	if err != nil {
		log.Printf("NewUsbDevs failed! %s", err)
		return nil, err
	}
	if len(devs) == 0 {
		log.Printf("vid: %v, pid %v not found!", vid, pid)
		return nil, nil
	}

	return devs[0], nil
}

func GetUsbDevsWithVidPid(vid, pid gousb.ID) ([]*gousb.Device, error) {
	log.Printf("GetUsbDevsWithVidPid")

	// Initialize a new gousb Context.
	ctx := gousb.NewContext()
	//defer ctx.Close()

	devs, err := GetDevices(ctx, vid, pid)
	// All returned devices are now open and will need to be closed.
	//for _, d := range devs {
	//	defer d.Close()
	//}
	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}
	return devs, nil
}

func GetDevices(ctx *gousb.Context, vid, pid gousb.ID) ([]*gousb.Device, error) {
	// Iterate through available Devices, finding all that match a known VID/PID.
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		fmt.Printf("vid %s pid %s\n", desc.Vendor, desc.Product)
		// this function is called for every device present.
		// Returning true means the device should be opened.
		return desc.Vendor == vid && desc.Product == pid
	})
	// All returned devices are now open and will need to be closed.
	/*for _, d := range devs {
		defer d.Close()
	}*/
	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}
	return devs, nil
}

//func GetUsbDevDesc(ctx *context.Context) error {
	/*dev, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	if dev == nil {
		return nil
	}

	// dev found
	//defer dev.Close()

	log.Printf("Camera device info: %v", dev)
	log.Printf("Camera device desc: %v", dev.Desc)
	log.Printf("bus %d, addres %d\n", dev.Desc.Bus, dev.Desc.Address)
	*/
//	return nil
//}

func SendUsbInfoToServer(usbClient api.UsbClient) error {
	ctx := context.Background()

	// TODO - replace with enumerated device list
	// HACK - for now just hardcode the camera device vid/pid
	vid, pid := gousb.ID(0x04b4), gousb.ID(0x00f3)
	dev, err := OpenUsbDev(vid, pid)
	if err != nil {
		return err
	}

	// create and send PutUsbDevDescReq
	devDescReq, err := NewPutUsbDevDescReq(usbClient, dev)
	if err != nil {
		log.Fatalf("Connect failed! %s", err)
	}

	usbClient.PutUsbDevDesc(ctx, devDescReq)

	// create and send PutUsbDevConnReq

	devConnReq, err := NewPutUsbDevConnReq(usbClient, dev)
	if err != nil {
		log.Fatalf("Connect failed! %s", err)
	}

	usbClient.PutUsbDevConnState(ctx, devConnReq)
	// put usb desc and usb conn state

	return nil
}

func NewPutUsbDevDescReq(client api.UsbClient,
	dev *gousb.Device) (*api.PutUsbDevDescRequest, error) {
	desc := dev.Desc
	spec := api.BCD{BcdValue: uint32(desc.Spec)}
	devVer := api.BCD{BcdValue: uint32(desc.Device)}
/*	man, err := dev.Manufacturer()
	if err != nil {
		return nil, err
	}
	prodstr, err := dev.Product()
	if err != nil {
		return nil, err
	}
	sn, err := dev.SerialNumber()
	if err != nil {
		return nil, err
	}*/

	putUsbDevDesc := api.PutUsbDevDescRequest{
		Spec:           &spec,
		Class:          uint32(desc.Class),
		SubClass:       uint32(desc.SubClass),
		Protocol:       uint32(desc.Protocol),
		MaxCtrlPktSize: uint32(desc.MaxControlPacketSize),
		VendorID:       uint32(desc.Vendor),
		ProductID:      uint32(desc.Product),
		Version:        &devVer,
	/*	Manufacturer:   man,
		Product:        prodstr,
		SerialNum:      sn,*/
		NumConfigs:     uint32(len(desc.Configs)),
	}

	return &putUsbDevDesc, nil
}

func NewPutUsbDevConnReq(client api.UsbClient,
	dev *gousb.Device) (*api.PutUsbDevConnRequest, error) {
	desc := dev.Desc
	path := make([]int32, 1)
	for i, p := range desc.Path {
		path[i] = int32(p)
	}

	putUsbDevConn := api.PutUsbDevConnRequest{
		Bus:     int32(desc.Bus),
		Address: int32(desc.Address),
		Speed:   api.UsbSpeed(desc.Speed),
		HubPort: int32(desc.Port),
	}

	return &putUsbDevConn, nil
}
