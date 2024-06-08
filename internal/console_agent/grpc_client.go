package console_agent

import (
	_ "log"

	api "github.com/edcox96/devmon/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcUsbClient(addr string) (api.UsbClient, error) {
	// connect using an insecure tcp link to server at addr
	ccOpts := []grpc.DialOption {
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.NewClient(addr, ccOpts...)
	if err != nil {
		return nil, err
	}

	client := api.NewUsbClient(cc)

	return client, nil
}
