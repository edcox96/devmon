package main

import (
	"flag"
	"fmt"
	"log"

	cons_agent "github.com/edcox96/devmon/internal/console_agent"
	usb_devs "github.com/edcox96/devmon/internal/console_agent/usb_devs"
)

var (
	_ = flag.Int("debug", 1, "libusb debug level (0..3)")
)

func main() {
	flag.Parse()
	fmt.Printf("console_agent\n")

	usbClient, err := cons_agent.NewGrpcUsbClient("localhost:8080")
	if err != nil {
		log.Fatalf("NewGrpcUsbClient failed! %s", err)
	}

	err = usb_devs.SendUsbInfoToServer(usbClient)
	if err != nil {
		log.Fatalf("SendUsbInfoToServer failed, %s", err)
	}
}
