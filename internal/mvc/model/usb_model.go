package model

import (
	_ "log"

    api "github.com/edcox96/devmon/api/v1"
)

type UsbModel struct {
	devDesc *UsbDevDesc
	devConn *UsbDevConnState
}

var usbModel = UsbModel {}

func NewUsbModel() *UsbModel {
	return &usbModel
}

func (um *UsbModel) PutUsbDevDesc(pbDesc *api.PutUsbDevDescRequest) error {
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

	um.devDesc = desc

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

func (um *UsbModel) PutUsbDevConnState(pbConn *api.PutUsbDevConnRequest) error {
	conn := newUsbDevConnState()
	
	conn.bus     = pbConn.Bus
	conn.address = pbConn.Address
	conn.port    = pbConn.HubPort
	conn.speed   = int32(pbConn.Speed.Enum().Number())
	conn.path    = pbConn.Path

	um.devConn = conn
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
