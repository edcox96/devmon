package event_logging

import (
	"fmt"
    "os"
	"bufio"
	"time"
	"strings"
)

type EventLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

type Event struct {
	SequenceNumber uint64
	EventTime time.Time
	EventType EventType
	Key string
	Value string
}

type EventType byte

const (
	_						= iota
	EventDelete EventType	= iota
	EventPut
)

type FileEventLogger struct {
	events chan<-Event
	errors <-chan error
	lastSequence uint64
	file *os.File
}

func NewFileEventLogger(filename string) (EventLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open event log file: %w", err)
	}
	return &FileEventLogger{file: file}, nil
}

func (l *FileEventLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value, EventTime: time.Now()}
}

func (l *FileEventLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: "Delete", EventTime: time.Now()}
}

func (l *FileEventLogger) Err() <- chan error {
	return l.errors
}

func (l *FileEventLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(l.file, "%s\t%d\t%d\t%s\t%s\n", e.EventTime.UTC().Format(time.UnixDate), l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileEventLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, "\t")
			e.EventTime, _ = time.Parse("Thu May 30 19:50:02 UTC 2024", fields[1])
			fmt.Println(fields[1:])

			if _, err := fmt.Sscanf(strings.Join(fields[1:], "\t"), "%d\t%d\t%s\t%s",
                &e.SequenceNumber, &e.EventType, &e.Key, &e.Value); err != nil {

                outError <- fmt.Errorf("input parse error: %w", err)
                return
            }

			if l.lastSequence >= e.SequenceNumber {
                outError <- fmt.Errorf("event numbers out of sequence")
                return
            }

			l.lastSequence = e.SequenceNumber
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
            outError <- fmt.Errorf("event log read failure: %w", err)
            return
        }
	}()

	return outEvent, outError
}