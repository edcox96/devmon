package event_logging

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/require"
)

var eventLogger EventLogger
var err error

func TestRun(t *testing.T) {
	fmt.Println("Creating new event logger")
	eventLogger, err = NewFileEventLogger("event.log")
	require.Equal(t, err, nil, "Error creating new Event Logger")

	fmt.Println("Running the event logger")
	eventLogger.Run()

	fmt.Println("Checking events channel length is zero")
	require.Equal(t, 0, len(eventLogger.GetEvents()), "Events channel length is not 1")
	fmt.Println("Checking errors channel length is zero")
	require.Equal(t, 0, len(eventLogger.GetErrors()), "Errors channel length is not 0")
	fmt.Println("Checking sequence number is zero")
	require.Equal(t, 0, eventLogger.GetSequenceNumber(), "Sequence number is not 1")

	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()

	keys := [4]string{"a", "b", "c", "d"}
	values := [4]string{"1", "b", "c", "d"}
	var receivedKeys [4]string
	var receivedValues [4]string

	for i := 0; i < 4; i++ {
		eventLogger.WritePut(keys[i], values[i])
	}

	i := 0
	for receivedKeys[3] == "" {
		select {
		case event := <-eventsChannel:
			fmt.Println(event)
			receivedKeys[i] = event.Key
			receivedValues[i] = event.Value
			i++
		case err := <-errorsChannel:
			fmt.Println(err)
			require.Equal(t, err, nil, err)
		default:
		}
	}

	require.Equal(t, keys, receivedKeys, "Sent and reeived keys not equal.")
	require.Equal(t, values, receivedValues, "Sent and received values not equal.")
	
}

func TestReadEvents(t *testing.T) {
	
}