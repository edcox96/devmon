package console_agent

import (
	"context"
	"fmt"
	"log"

	api "github.com/edcox96/devmon/api/v1"
	devsim "github.com/edcox96/devmon/internal/devsim"
	model "github.com/edcox96/devmon/internal/server/mvc/model"
	"github.com/google/gousb"
)

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

func SendUsbInfoToServer(usbClient api.UsbClient, conId uint64) error {
	ctx := context.Background()

	regConReq, err := NewRegisterConsoleReq(usbClient, conId)
	if err != nil {
		log.Printf("NewRegisterConsoleReq failed! %s", err)
		return err
	}

	usbClient.RegisterConsole(ctx, regConReq)

	// TODO - replace with enumerated device list
	// HACK - for now just hardcode the camera device vid/pid
	vid, pid := gousb.ID(0x04b4), gousb.ID(0x00f9)
	dev, err := OpenUsbDev(vid, pid)
	if err != nil {
		return err
	}

	regDevReq, err := NewRegisterUsbDeviceReq(usbClient,
		conId, uint32(vid), uint32(pid))
	if err != nil {
		log.Printf("NewRegisterDeviceReq failed! %s", err)
		return err
	}

	usbClient.RegisterUsbDevice(ctx, regDevReq)

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

	return nil
}

func NewRegisterConsoleReq(client api.UsbClient,
	    id uint64) (*api.RegisterConsoleRequest, error) {
	regCons := &api.RegisterConsoleRequest{}

	con, err := devsim.GetConsole(id)
	if err != nil {
		log.Printf("devsim.GetConsole(%d) failed! %s", id, err)
		return nil, err
	}

	regCons.ID = con.ID
	regCons.Name = con.Name
	regCons.SN = con.SN

	return regCons, nil
}

func NewRegisterUsbDeviceReq(client api.UsbClient,
	id uint64, vid uint32, pid uint32) (*api.RegisterUsbDeviceRequest, error) {
	regDev := &api.RegisterUsbDeviceRequest{}

	con, err := devsim.GetConsole(id)
	if err != nil {
		log.Printf("devsim.GetConsole(%d) failed! %s", id, err)
		return nil, err
	}

	regDev.ConsID = con.ID
	regDev.DevType = int32(model.UsbDevTypeCamera)
	regDev.DevTypeNum = 1
	regDev.Vid = vid
	regDev.Pid = pid

	return regDev, nil
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
		NumConfigs: uint32(len(desc.Configs)),
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
