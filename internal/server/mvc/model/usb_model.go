package model

import (
	"log"
	"errors"

    api "github.com/edcox96/devmon/api/v1"
)

func NewUsbModel() *UsbModel {
	um := &usbModel
	um.Consoles = make([]*Console, 0, 2)
	return um
}

var usbModel = UsbModel {}

type UsbModel struct {
	Consoles []*Console
}

func (um *UsbModel) RegisterConsole(pbCons *api.RegisterConsoleRequest) error {
	con := um.NewConsole()
	
	con.Name = pbCons.Name
	con.SN = pbCons.SN
	con.ID = pbCons.ID
	
	con.UsbDevs = make([]*UsbDevice, 0, 2)

	um.Consoles = append(um.Consoles, con)
	return nil
}

func (um *UsbModel) NewConsole() *Console {
	return new(Console)
}

type Console struct {
	Name string
	SN   string
	ID   uint64
	UsbDevs []*UsbDevice
}

type UsbDevice struct {
	consId  uint64
	devType UsbDevType
	devTypeNum uint32
	vid     uint32
	pid 	uint32
	devDesc *UsbDevDesc
	devConn *UsbDevConnState
}

type UsbDevType byte

const (
	UsbDevTypeUnknown UsbDevType = iota
	UsbDevTypeCamera  UsbDevType = iota
)

func (con *Console) RegisterUsbDevice(pbDev *api.RegisterUsbDeviceRequest) error {
	if pbDev.ConsID != con.ID {
		log.Printf("RegisterUsbDevice: dev conID %d, conID %d mismatch!",
				   pbDev.ConsID, con.ID)
		return errors.New("invalid console ID")
	}

	dev := new(UsbDevice)
	dev.consId = pbDev.ConsID
	dev.devType = UsbDevType(pbDev.DevType)
	dev.pid = pbDev.Pid
	dev.vid = pbDev.Vid
	dev.devTypeNum = pbDev.DevTypeNum
	
	con.UsbDevs = append(con.UsbDevs, dev)
	return nil
}

func (dev *UsbDevice) PutUsbDevDesc(pbDesc *api.PutUsbDevDescRequest) error {
	desc := newUsbDevDesc()
	
	desc.usbVer.bcdValue = pbDesc.Spec.BcdValue
	desc.class = pbDesc.Class
	desc.subClass = pbDesc.SubClass
	desc.protocol = pbDesc.Protocol
	desc.maxCtrlPktSize = pbDesc.MaxCtrlPktSize
	desc.idVendor = uint16(pbDesc.VendorID)
	desc.idProduct = uint16(pbDesc.ProductID)
	desc.devVer.bcdValue = pbDesc.Version.BcdValue
	desc.manufacturer = pbDesc.Manufacturer
	desc.product = pbDesc.Product
	desc.numConfigs = pbDesc.NumConfigs

	dev.devDesc = desc
	return nil
}

func newUsbDevDesc() *UsbDevDesc {
	return &UsbDevDesc {}
}

type UsbDevDesc struct {
    usbVer   BCD // USB Spec device is compliant with
    class    uint32
    subClass uint32
    protocol uint32
    maxCtrlPktSize uint32
    idVendor  uint16
    idProduct uint16
    devVer    BCD
    manufacturer  string
    product       string
    numConfigs uint32
}

func (dev *UsbDevice) PutUsbDevConnState(pbConn *api.PutUsbDevConnRequest) error {
	conn := newUsbDevConnState()
	
	conn.bus     = pbConn.Bus
	conn.address = pbConn.Address
	conn.port    = pbConn.HubPort
	conn.speed   = int32(pbConn.Speed.Enum().Number())
	conn.path    = pbConn.Path

	dev.devConn = conn
	return nil
}

func newUsbDevConnState() *UsbDevConnState {
	return &UsbDevConnState {}
}

type UsbDevConnState struct {
    bus     int32
    address int32
    port    int32
    speed   int32 //UsbSpeed
    path    []int32
}

type BCD struct {
    bcdValue  uint32
}
