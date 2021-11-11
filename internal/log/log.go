package log

import (
	"errors"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	Records []Record
}

var (
	ErrOutOffset = errors.New("Offset out of range")
)

type Record struct {
	Value  []byte
	Offset uint64
}

func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Offset = uint64(len(l.Records))
	l.Records = append(l.Records, record)
	return record.Offset, nil
}

func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var record Record

	if offset > uint64(len(l.Records)) || len(l.Records) == 0 {
		return Record{}, ErrOutOffset
	}

	record = l.Records[offset]
	return record, nil
}
