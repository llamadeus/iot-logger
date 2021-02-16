package tailbuf

import (
	"sync"
)

type TailBuf struct {
	buffer []interface{}
	size   int
	length int
	cursor int
	mu     sync.Mutex
}

func New(size int) *TailBuf {
	return &TailBuf{
		buffer: make([]interface{}, size),
		size:   size,
	}
}

func (tb *TailBuf) Read() []interface{} {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	var result []interface{}

	// Copy from cursor until end
	if tb.cursor != tb.length {
		result = append(result, tb.buffer[tb.cursor:tb.length]...)
	}

	// Copy from start until cursor
	if tb.cursor > 0 {
		result = append(result, tb.buffer[0:tb.cursor]...)
	}

	return result
}

func (tb *TailBuf) Write(value interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.buffer[tb.cursor] = value
	tb.cursor = (tb.cursor + 1) % tb.size

	if tb.length < tb.size {
		tb.length += 1
	}
}
