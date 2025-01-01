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

type Writer struct {
	*protocol.Writer

	protocolID int32
}

func NewWriter(w *protocol.Writer, protocolID int32) *Writer {
	return &Writer{w, protocolID}
}

func (w *Writer) SetProtocolID(protocolID int32) { w.protocolID = protocolID }
func (w *Writer) ProtocolID() int32              { return w.protocolID }

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

func IOStackRequestAction(io protocol.IO, x *protocol.StackRequestAction) {
	if IsReader(io) {
		var id uint8
		io.Uint8(&id)
		if !lookupStackRequestAction(id, x) {
			io.UnknownEnumOption(id, "stack request action type")
			return
		}
	} else {
		var id byte
		if !lookupStackRequestActionType(*x, &id) {
			io.UnknownEnumOption(fmt.Sprintf("%T", *x), "stack request action type")
		}
		io.Uint8(&id)
	}
	(*x).Marshal(io)
}

func IORecipe(io protocol.IO, recipe *Recipe) {
	if IsReader(io) {
		var recipeType int32
		io.Varint32(&recipeType)
		if !lookupRecipe(recipeType, recipe) {
			io.UnknownEnumOption(recipeType, "crafting data recipe type")
			return
		}
		(*recipe).Unmarshal(io.(*Reader).Reader)
	} else {
		var recipeType int32
		if !lookupRecipeType(*recipe, &recipeType) {
			io.UnknownEnumOption(fmt.Sprintf("%T", *recipe), "crafting recipe type")
		}
		io.Varint32(&recipeType)
		(*recipe).Marshal(io.(*Writer).Writer)
	}
}
