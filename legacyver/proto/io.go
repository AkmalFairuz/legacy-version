package proto

import (
	"fmt"
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

func TransactionDataType(io protocol.IO, x *InventoryTransactionData) {
	if IsReader(io) {
		var transactionType uint32
		io.Varuint32(&transactionType)
		if !lookupTransactionData(transactionType, x) {
			io.UnknownEnumOption(transactionType, "inventory transaction data type")
		}
	} else {
		var id uint32
		if !lookupTransactionDataType(*x, &id) {
			io.UnknownEnumOption(fmt.Sprintf("%T", x), "inventory transaction data type")
		}
		io.Varuint32(&id)
	}
}

func PlayerInventoryAction(io protocol.IO, x *protocol.UseItemTransactionData) {
	io.Varint32(&x.LegacyRequestID)
	if x.LegacyRequestID < -1 && (x.LegacyRequestID&1) == 0 {
		protocol.Slice(io, &x.LegacySetItemSlots)
	}
	protocol.Slice(io, &x.Actions)
	io.Varuint32(&x.ActionType)
	if IsProtoGTE(io, ID712) {
		io.Varuint32(&x.TriggerType)
	}
	io.BlockPos(&x.BlockPosition)
	io.Varint32(&x.BlockFace)
	io.Varint32(&x.HotBarSlot)
	io.ItemInstance(&x.HeldItem)
	io.Vec3(&x.Position)
	io.Vec3(&x.ClickedPosition)
	io.Varuint32(&x.BlockRuntimeID)
	if IsProtoGTE(io, ID712) {
		io.Varuint32(&x.ClientPrediction)
	}
}
