package event_logging

import (
	"testing"
	"fmt"
	"math/rand"

	"github.com/stretchr/testify/require"
)

var eventLogger EventLogger
var err error

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

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
	 
}

func TestWritePut(t *testing.T) {

	eventLogger.Run()
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()

	n := 50

	keys := make([]string, n)
	values := make([]string, n)
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)

	go func() {
		for i := 0; i < n; i++ {
			key := letterBytes[i]
			value := RandStringBytes(5)
			keys[i] = string(key)
			values[i] = value
			eventLogger.WritePut(keys[i], values[i])
		}
	}()


	i := 0
	for i < len(keys) {
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

	fmt.Println(receivedKeys)
	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
	require.Equal(t, values, receivedValues, "Sent and received values not equal.")
	
}

func TestWriteDelete(t *testing.T) {
	eventLogger.Run()
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()

	n := 50

	keys := make([]string, n)
	receivedKeys := make([]string, n)

	go func() {
		for i := 0; i < n; i++ {
			key := letterBytes[i]
			keys[i] = string(key)
			eventLogger.WriteDelete(keys[i])
		}
	}()


	i := 0
	for i < len(keys) {
		select {
		case event := <-eventsChannel:
			fmt.Println(event)
			receivedKeys[i] = event.Key
			i++
		case err := <-errorsChannel:
			fmt.Println(err)
			require.Equal(t, err, nil, err)
		default:
		}
	}

	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
}

func TestFullRun(t *testing.T) {
	eventLogger.Run()
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()

	n := 50

	keys := make([]string, n)
	values := make([]string, n)
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)

	go func() {
		for i := 0; i < n; i++ {
			operation := rand.Intn(1)
			key := letterBytes[i]
			keys[i] = string(key)
			
			if operation == 0 {
				value := RandStringBytes(5)
				values[i] = value
				eventLogger.WritePut(keys[i], values[i])
			} else {
				values[i] = ""
				eventLogger.WriteDelete(keys[i])
			}
		}
	}()

	i := 0
	for i < len(keys) {
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

	fmt.Println(receivedKeys)
	fmt.Println(receivedValues)
	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
	require.Equal(t, values, receivedValues, "Sent and received values not equal.")
}

func TestReadEvents(t *testing.T) {
	
}