package utility

import (
	"bytes"
	"sync"
)

// SafeBuffer is a thread safe version of bytes.Buffer
type SafeBuffer struct {
	m sync.RWMutex
	b bytes.Buffer
}

// Read is a thread safe version of bytes.Buffer::Read
func (b *SafeBuffer) Read(p []byte) (n int, err error) {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.b.Read(p)
}

// Write is a thread safe version of bytes.Buffer::Write
func (b *SafeBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

// String is a thread safe version of bytes.Buffer::String
func (b *SafeBuffer) String() string {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.b.String()
}

// Bytes is a thread safe version of bytes.Buffer::Bytes
func (b *SafeBuffer) Bytes() []byte {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.b.Bytes()
}
