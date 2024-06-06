package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	usb_devs "github.com/edcox96/devmon/console_agent/usb_devs"
)

var (
	_ = flag.Int("debug", 1, "libusb debug level (0..3)")
)

func main() {
	flag.Parse()
	fmt.Printf("console_agent\n")

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("net.listen failed: %s", err)
	}

	usbClient, err := usb_devs.Connect(lis)
	if err != nil {
		log.Fatalf("Connect failed! %s", err)
	}

	err = usb_devs.SendUsbInfoToServer(usbClient)
	if err != nil {
		log.Fatalf("SendUsbInfoToServer failed, %s", err)
	}
}
