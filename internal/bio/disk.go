package bio

import (
	"encoding/binary"
	"io"

	"github.com/chzyer/logex"
)

var (
	ErrReaderBufferFull = logex.Define("reader buffer is full")
	ErrWriterBufferFull = logex.Define("reader writer is full")
)

type Diskable interface {
	Size() int
	ReadDisk(r *Reader) error
	WriteDisk(w *Writer)
}

type Writer struct {
	data   []byte
	offset int
}

func NewWriter(data []byte) *Writer {
	return &Writer{data: data}
}

func (w *Writer) Available() int {
	return len(w.data) - w.offset
}

func (w *Writer) WriteDisk(d Diskable) error {
	if w.Available() < d.Size() {
		return ErrWriterBufferFull.Trace()
	}
	d.WriteDisk(w)
	return nil
}

func (w *Writer) Int32(n int32) {
	binary.BigEndian.PutUint32(w.data[w.offset:], uint32(n))
	w.offset += 4
	return
}

func (w *Writer) Byte(b []byte) {
	copy(w.data[w.offset:], b)
	w.offset += len(b)
}

func (w *Writer) Int64(n int64) {
	binary.BigEndian.PutUint64(w.data[w.offset:], uint64(n))
	w.offset += 8
}

type Reader struct {
	data   []byte
	offset int
}

func ReadAt(r io.ReaderAt, offset int64, d Diskable) error {
	blk := make([]byte, d.Size())
	n, err := r.ReadAt(blk, offset)
	if err != nil {
		return logex.Trace(err)
	}
	if n != len(blk) {
		return logex.Trace(io.EOF)
	}
	return logex.Trace(d.ReadDisk(NewReader(blk)))
}

func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

func (r *Reader) Byte(n int) []byte {
	ret := r.data[r.offset:r.offset:4]
	r.offset += 4
	return ret
}

func (r *Reader) Available() int {
	return len(r.data) - r.offset
}

func (r *Reader) Check(d Diskable, n int) error {
	if r.Available() < d.Size()*n {
		return ErrReaderBufferFull.Trace()
	}
	return nil
}

func (r *Reader) ReadDisk(d Diskable) error {
	if r.Available() < d.Size() {
		return ErrReaderBufferFull.Trace()
	}
	return d.ReadDisk(r)
}

func (r *Reader) Int32() int32 {
	ret := int32(binary.BigEndian.Uint32(r.data[r.offset:]))
	r.offset += 4
	return ret
}

func (r *Reader) Int64() int64 {
	ret := int64(binary.BigEndian.Uint64(r.data[r.offset:]))
	r.offset += 8
	return ret
}