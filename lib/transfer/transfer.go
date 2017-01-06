package transfer

import (
	"fmt"
	"io"
	"sync/atomic"
)

type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.3fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.1fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.0fKB", b/KB)
	}
	return fmt.Sprintf("%.0fB", b)
}

type Transfer interface {
	Transferred() ByteSize
	Length() ByteSize
}

// ReaderTransfer monitors how much of a reader has been read
type ReaderTransfer struct {
	length ByteSize
	read   uint64
	reader io.Reader
	kill   int32
}

// NewReaderTransfer instantiates a ReadTransfer struct
func NewReaderTransfer(reader io.Reader, length int) *ReaderTransfer {
	rt := new(ReaderTransfer)
	rt.length = ByteSize(length)
	rt.read = 0
	rt.reader = reader
	rt.kill = 0
	return rt
}

func (rt *ReaderTransfer) Read(p []byte) (int, error) {
	if atomic.LoadInt32(&rt.kill) == 1 {
		return 0, fmt.Errorf("transfer killed")
	}
	n, err := rt.reader.Read(p)
	atomic.AddUint64(&rt.read, uint64(n))
	return n, err
}

func (rt *ReaderTransfer) KillTransfer() {
	atomic.StoreInt32(&rt.kill, 1)
}

func (rt *ReaderTransfer) Transferred() ByteSize {
	return ByteSize(atomic.LoadUint64(&rt.read))
}

func (rt *ReaderTransfer) Length() ByteSize {
	return rt.length
}

// WriteCloserTransfer monitors how much of a WriteCloser has been written
type WriteCloserTransfer struct {
	length      ByteSize
	written     uint64
	writeCloser io.WriteCloser
}

// NewWriteCloserTransfer instantiates a new WriteCloserTransfer
func NewWriteCloserTransfer(writeCloser io.WriteCloser, length int) *WriteCloserTransfer {
	wct := new(WriteCloserTransfer)
	wct.length = ByteSize(length)
	wct.written = 0
	wct.writeCloser = writeCloser
	return wct
}

func (wct *WriteCloserTransfer) Write(p []byte) (int, error) {
	n, err := wct.writeCloser.Write(p)
	atomic.AddUint64(&wct.written, uint64(n))
	return n, err
}

func (wct *WriteCloserTransfer) Close() error {
	return wct.writeCloser.Close()
}

func (wct *WriteCloserTransfer) Transferred() ByteSize {
	return ByteSize(atomic.LoadUint64(&wct.written))
}

func (wct *WriteCloserTransfer) Length() ByteSize {
	return wct.length
}
