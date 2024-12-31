package proto

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type IO interface {
	protocol.IO

	SetProtocolID(protocolID int32)
	ProtocolID() int32
}

type Reader struct {
	*protocol.Reader

	protocolID int32
}

func NewReader(r *protocol.Reader, protocolID int32) *Reader {
	return &Reader{
		Reader:     r,
		protocolID: protocolID,
	}
}

func (r *Reader) SetProtocolID(protocolID int32) { r.protocolID = protocolID }
func (r *Reader) ProtocolID() int32              { return r.protocolID }
func (r *Reader) Reads() bool                    { return true }

type Writer struct {
	*protocol.Writer

	protocolID int32
}

func NewWriter(w *protocol.Writer, protocolID int32) *Writer {
	return &Writer{w, protocolID}
}

func (w *Writer) SetProtocolID(protocolID int32) { w.protocolID = protocolID }
func (w *Writer) ProtocolID() int32              { return w.protocolID }
func (w *Writer) Reads() bool                    { return false }

func IsReader(r protocol.IO) bool {
	_, ok := r.(*Reader)
	return ok
}

func IsWriter(w protocol.IO) bool {
	_, ok := w.(*Writer)
	return ok
}

func EmptySlice[T any](io protocol.IO, slice *[]T) {
	if IsReader(io) {
		*slice = make([]T, 0)
	}
}
