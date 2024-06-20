package event_logging

import (
	"testing"
	"fmt"
	"math/rand"
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/stretchr/testify/require"
)

var eventLogger EventLogger
var err error
var readKeys []string
var readValues []string
var lock = sync.NewCond(&sync.Mutex{})

const n = 50
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func CreateNewEventLogger(t *testing.T) {

	fmt.Println("Creating new event logger")
	eventLogger, err = NewFileEventLogger("event.log")
	require.Equal(t, err, nil, "Error creating new Event Logger")
}

func TruncateFile(t *testing.T) {
	
	fmt.Println("Truncating log")
	err := os.Truncate("event.log", 0)
	require.Equal(t, err, nil, "Error truncating Event Log")
}

func ReceiveEvents(t *testing.T, eventsChannel chan Event, errorsChannel chan error, receivedKeys []string, receivedValues []string) {
	
	i := 0
	for i < n {
		select {

		case event := <-eventsChannel:
			receivedKeys[i] = event.Key
			receivedValues[i] = event.Value
			i++

		case err := <-errorsChannel:
			fmt.Println("error")
			fmt.Println(err)
			require.Equal(t, err, nil, err)

		default:
		}
	}
}

func WaitForEventLogger(c *sync.Cond) {
	c.L.Lock()
	for eventLogger.GetSequenceNumber() < n {
		c.Wait()
	}
	c.L.Unlock()
}

func ReadEventLog(readKeys []string, readValues []string) {

	fmt.Println("Reading Event Log")
	file, _ := os.OpenFile("event.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	scanner := bufio.NewScanner(file)
	
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		readKeys[i] = fields[3]
		readValues[i] = fields[4]
		i++
	}
}

func TestRun(t *testing.T) {
	
	// Creating the new event logger if needed or clearing the entries from the existing event logger.
	CreateNewEventLogger(t)

	fmt.Println("Running event logger")
	eventLogger.Run()

	fmt.Println("Checking events channel length is zero")
	require.Equal(t, 0, len(eventLogger.GetEvents()), "Events channel length is not 1")
	fmt.Println("Checking errors channel length is zero")
	require.Equal(t, 0, len(eventLogger.GetErrors()), "Errors channel length is not 0")
	fmt.Println("Checking sequence number is zero")
	require.Equal(t, 0, eventLogger.GetSequenceNumber(), "Sequence number is not 1")
}

func TestWritePut(t *testing.T) {
	fmt.Println("\nTestWritePut")

	// Creating the new event logger and clearing the entries from the existing event log.
	CreateNewEventLogger(t)
	TruncateFile(t)

	fmt.Println("Running event logger")
	eventLogger.Run()

	// Getting the events and errors channels from the event logger and creating slices for sent and received event values.
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()
	cond := eventLogger.GetWritten()
	keys := make([]string, n)
	values := make([]string, n)
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)

	fmt.Println("Creating put events and writing them to the event logger")
	go func() {
		for i := 0; i < n; i++ {

			// Getting the next key and generating a new random value.
			key := letterBytes[i]
			value := RandStringBytes(5)
			keys[i] = string(key)
			values[i] = value
			eventLogger.WritePut(keys[i], values[i])
		}
	}()

	// Receiving the events from the test channel and waiting for the event log to finish writing events to file.
	fmt.Println("Receiving the events from the test channel and waiting for the event log to finish writing")
	ReceiveEvents(t, eventsChannel, errorsChannel, receivedKeys, receivedValues)
	WaitForEventLogger(cond)

	// Checking that the sent and received values are equal.
	fmt.Println("Checking that the sent and received values are equal")
	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
	require.Equal(t, values, receivedValues, "Sent and received values not equal.")
}

func TestWriteDelete(t *testing.T) {
	fmt.Println("\nTestWriteDelete")

	// Creating the new event logger and clearing the entries from the existing event logger.
	CreateNewEventLogger(t)
	TruncateFile(t)

	fmt.Println("Running event logger")
	eventLogger.Run()

	// Getting the events and errors channels from the event logger and creating slices for sent and received event values.
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()
	cond := eventLogger.GetWritten()
	keys := make([]string, n)
	values := make([]string, n)
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)

	fmt.Println("Creating delete events and writing them to the event logger")
	go func() {
		for i := 0; i < n; i++ {

			// Getting the next key and setting the value to "Delete" to match expected received results.
			key := letterBytes[i]
			keys[i] = string(key)
			values[i] = "Delete"
			eventLogger.WriteDelete(keys[i])
		}
	}()

	// Receiving the events from the test channel and waiting for the event log to finish writing events to file.
	fmt.Println(" Receiving the events from the test channel and waiting for the event log to finish writing")
	ReceiveEvents(t, eventsChannel, errorsChannel, receivedKeys, receivedValues)
	WaitForEventLogger(cond)

	// Checking that the sent and received values are equal.
	fmt.Println("Checking that the sent and received values are equal")
	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
}

func TestFullRun(t *testing.T) {
	fmt.Println("\nTestFullRun")

	// Creating the new event logger and clearing the entries from the existing event log.
	CreateNewEventLogger(t)
	TruncateFile(t)

	fmt.Println("Running event logger")
	eventLogger.Run()

	// Getting the events and errors channels from the event logger and creating slices for sent, received, and read event values.
	eventsChannel := eventLogger.GetEvents()
	errorsChannel := eventLogger.GetErrors()
	cond := eventLogger.GetWritten()
	keys := make([]string, n)
	values := make([]string, n)
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)
	readKeys = make([]string, n)
	readValues = make([]string, n)

	fmt.Println("Creating random put and delete events and writing them to the event logger")
	go func() {
		for i := 0; i < n; i++ {
			operation := rand.Intn(2)
			key := letterBytes[i]
			keys[i] = string(key)
			
			if operation == 0 {
				value := RandStringBytes(5)
				values[i] = value
				eventLogger.WritePut(keys[i], values[i])
			} else {
				values[i] = "Delete"
				eventLogger.WriteDelete(keys[i])
			}
		}
	}()

	// Receiving the events from the test channel and waiting for the event log to finish writing events to file.
	fmt.Println("Receiving the events from the test channel and waiting for the event log to finish writing")
	ReceiveEvents(t, eventsChannel, errorsChannel, receivedKeys, receivedValues)
	WaitForEventLogger(cond)

	// Reading the event log file storing the key and value from each line.
	fmt.Println("Reading the event log file storing the key and value from each line")
	ReadEventLog(readKeys, readValues)
	
	// Checking that the sent and received values are equal.
	fmt.Println("Checking that the sent and received values are equal")
	require.Equal(t, keys, receivedKeys, "Sent and received keys not equal.")
	require.Equal(t, values, receivedValues, "Sent and received values not equal.")
	require.Equal(t, keys, readKeys, "Sent and read keys not equal.")
	require.Equal(t, receivedKeys, readKeys, "Received and read keys are not equal")
	require.Equal(t, values, readValues, "Sent and read values not equal.")
	require.Equal(t, receivedValues, readValues, "Received and read values are not equal")
}

func TestReadEvents(t *testing.T) {
	fmt.Println("\nTestReadEvents")

	// Creating a new event logger but not truncating the file so it can be read from.
	fmt.Println("Creating event logger")
	CreateNewEventLogger(t)

	// Creating slices for receiving values from channels.
	receivedKeys := make([]string, n)
	receivedValues := make([]string, n)

	// Receiving events from the channels.
	eventsChannel, errorsChannel := eventLogger.ReadEvents()
	ReceiveEvents(t, eventsChannel, errorsChannel, receivedKeys, receivedValues)

	// Checking that the sent and received values are equal.
	fmt.Println("Checking that the sent and received values are equal")
	require.Equal(t, receivedKeys, readKeys, "Received and read keys are not equal")
	require.Equal(t, receivedValues, readValues, "Received and read values are not equal")
}