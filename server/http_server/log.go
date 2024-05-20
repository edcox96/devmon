package http_server

import (
	"fmt"
	"sync"
)

/* Struct for representing our Record type used in the Log type. */
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

/* Struct for representing our Log type containing the list of records and a mutex for controlling access to the records. */
type Log struct {
	mu      sync.Mutex
	records []Record
}

/* Func for creating and returning a pointer to a new Log. */
func NewLog() *Log {
	return &Log{}
}

/* Func for appending to the Log using the mutex to control access. */
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

/* Func for reading from the Log using the mutex to control access. */
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}
