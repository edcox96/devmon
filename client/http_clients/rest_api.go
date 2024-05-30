package http_client

import (
	"log"
	"net/http"
)

type UrlRecord struct {
	url string;
	sc int; // status code
}

var urls = []UrlRecord {
	// test getRootHandler
	{ "http://localhost:8080/", 404},			// path "/"
	{ "http://localhost:8080/index.html", 200}, // path "index.html"
	{ "http://localhost:8080/consoles", 404},    // path must have /devmon/v1/ prefix
	// test getConsolesListHandler
	{ "http://localhost:8080/devmon/v1/consoles/", 200 },
	// test getConsoleHandler
	{ "http://localhost:8080/devmon/v1/console/1", 200 },
	// test getDeviceListHandler
	{ "http://localhost:8080/devmon/v1/console/1/devices", 200 },
	// test getDeviceHandler
	{ "http://localhost:8080/devmon/v1/console/1/camera/2", 200 },
	// test getDeviceObjectHandler
	{ "http://localhost:8080/devmon/v1/console/1/camera/2/device_desc", 200 },
	// test getEventLogHandler
	{ "http://localhost:8080/devmon/v1/eventlog/", 200 },
	// test getSoftwareLogHandler
	{ "http://localhost:8080/devmon/v1/swlog/", 200 },
}

func ClientRestAPITests() {
	log.Printf("Client: Run devmon/v1 server REST API tests.")

	for _, u := range urls {
		resp, err := http.Head(u.url)
		if err != nil {
			log.Fatalf("http.Head err %s", err)
		} else if resp.StatusCode != u.sc {
			log.Printf("URL: %s, expectng status code %d got %d", u.url, u.sc, resp.StatusCode)
		}
	}
}
