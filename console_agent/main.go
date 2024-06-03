package main

import (
    _ "context"
    "flag"
    "fmt"
    "log"
    "net"

    api "github.com/edcox96/devmon/api/v1"
    usb_devs "github.com/edcox96/devmon/console_agent/usb_devs"
    "github.com/google/gousb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

var (
    _ = flag.Int("debug", 0, "libusb debug level (0..3)")
)

func main() {
    flag.Parse()
    fmt.Printf("console_agent\n")

     // hack to just get known device vid/pid for camera
    vid, pid := gousb.ID(0x04b4), gousb.ID(0x00f3)
    dev, err := usb_devs.OpenUsbDev(vid, pid)

    usbClient, err := Connect()
    if err != nil {
        log.Fatalf("Connect failed! %s", err)
    }

    err = usb_devs.SendUsbInfoToServer(usbClient, dev)
    if err != nil {
        log.Fatalf("SendUsbInfoToServer failed, %s", err)
    }
}

func Connect() (api.UsbClient, error) {
    l, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        return nil, err
    }
    defer l.Close()

    ccOpts := []grpc.DialOption {
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }
    cc, err := grpc.NewClient(l.Addr().String(), ccOpts...)
    if err != nil {
        return nil, err
    }

    client := api.NewUsbClient(cc)

    return client, nil
}
