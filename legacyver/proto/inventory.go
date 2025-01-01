package proto

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// InventoryTransactionData represents an object that holds data specific to an inventory transaction type.
// The data it holds depends on the type.
type InventoryTransactionData interface {
	// Marshal encodes/decodes a serialised inventory transaction data object.
	Marshal(r protocol.IO)
}

// lookupTransactionData looks up inventory transaction data for the ID passed.
func lookupTransactionData(id uint32, x *InventoryTransactionData) bool {
	switch id {
	case protocol.InventoryTransactionTypeNormal:
		*x = &protocol.NormalTransactionData{}
	case protocol.InventoryTransactionTypeMismatch:
		*x = &protocol.MismatchTransactionData{}
	case protocol.InventoryTransactionTypeUseItem:
		*x = &UseItemTransactionData{}
	case protocol.InventoryTransactionTypeUseItemOnEntity:
		*x = &protocol.UseItemOnEntityTransactionData{}
	case protocol.InventoryTransactionTypeReleaseItem:
		*x = &protocol.ReleaseItemTransactionData{}
	default:
		return false
	}
	return true
}

// lookupTransactionDataType looks up an ID for a specific transaction data.
func lookupTransactionDataType(x InventoryTransactionData, id *uint32) bool {
	switch x.(type) {
	case *protocol.NormalTransactionData:
		*id = protocol.InventoryTransactionTypeNormal
	case *protocol.MismatchTransactionData:
		*id = protocol.InventoryTransactionTypeMismatch
	case *UseItemTransactionData:
		*id = protocol.InventoryTransactionTypeUseItem
	case *protocol.UseItemOnEntityTransactionData:
		*id = protocol.InventoryTransactionTypeUseItemOnEntity
	case *protocol.ReleaseItemTransactionData:
		*id = protocol.InventoryTransactionTypeReleaseItem
	default:
		return false
	}
	return true
}

// UseItemTransactionData represents an inventory transaction data object sent when the client uses an item on
// a block.
type UseItemTransactionData struct {
	// LegacyRequestID is an ID that is only non-zero at times when sent by the client. The server should
	// always send 0 for this. When this field is not 0, the LegacySetItemSlots slice below will have values
	// in it.
	// LegacyRequestID ties in with the ItemStackResponse packet. If this field is non-0, the server should
	// respond with an ItemStackResponse packet. Some inventory actions such as dropping an item out of the
	// hotbar are still one using this packet, and the ItemStackResponse packet needs to tie in with it.
	LegacyRequestID int32
	// LegacySetItemSlots are only present if the LegacyRequestID is non-zero. These item slots inform the
	// server of the slots that were changed during the inventory transaction, and the server should send
	// back an ItemStackResponse packet with these slots present in it. (Or false with no slots, if rejected.)
	LegacySetItemSlots []protocol.LegacySetItemSlot
	// Actions is a list of actions that took place, that form the inventory transaction together. Each of
	// these actions hold one slot in which one item was changed to another. In general, the combination of
	// all of these actions results in a balanced inventory transaction. This should be checked to ensure that
	// no items are cheated into the inventory.
	Actions []protocol.InventoryAction
	// ActionType is the type of the UseItem inventory transaction. It is one of the action types found above,
	// and specifies the way the player interacted with the block.
	ActionType uint32
	// TriggerType is the type of the trigger that caused the inventory transaction. It is one of the trigger
	// types found in the constants above. If TriggerType is TriggerTypePlayerInput, the transaction is from
	// the initial input of the player. If it is TriggerTypeSimulationTick, the transaction is from a simulation
	// tick when the player is holding down the input.
	TriggerType uint32
	// BlockPosition is the position of the block that was interacted with. This is only really a correct
	// block position if ActionType is not UseItemActionClickAir.
	BlockPosition protocol.BlockPos
	// BlockFace is the face of the block that was interacted with. When clicking the block, it is the face
	// clicked. When breaking the block, it is the face that was last being hit until the block broke.
	BlockFace int32
	// HotBarSlot is the hot bar slot that the player was holding while clicking the block. It should be used
	// to ensure that the hot bar slot and held item are correctly synchronised with the server.
	HotBarSlot int32
	// HeldItem is the item that was held to interact with the block. The server should check if this item
	// is actually present in the HotBarSlot.
	HeldItem protocol.ItemInstance
	// Position is the position of the player at the time of interaction. For clicking a block, this is the
	// position at that time, whereas for breaking the block it is the position at the time of breaking.
	Position mgl32.Vec3
	// ClickedPosition is the position that was clicked relative to the block's base coordinate. It can be
	// used to find out exactly where a player clicked the block.
	ClickedPosition mgl32.Vec3
	// BlockRuntimeID is the runtime ID of the block that was clicked. It may be used by the server to verify
	// that the player's world client-side is synchronised with the server's.
	BlockRuntimeID uint32
	// ClientPrediction is the client's prediction on the output of the transaction. It is one of the client
	// prediction found in the constants above.
	ClientPrediction uint32
}

func (x *UseItemTransactionData) FromLatest(l *protocol.UseItemTransactionData) *UseItemTransactionData {
	return &UseItemTransactionData{
		LegacyRequestID:    l.LegacyRequestID,
		LegacySetItemSlots: l.LegacySetItemSlots,
		Actions:            l.Actions,
		ActionType:         l.ActionType,
		TriggerType:        l.TriggerType,
		BlockPosition:      l.BlockPosition,
		BlockFace:          l.BlockFace,
		HotBarSlot:         l.HotBarSlot,
		HeldItem:           l.HeldItem,
		Position:           l.Position,
		ClickedPosition:    l.ClickedPosition,
		BlockRuntimeID:     l.BlockRuntimeID,
		ClientPrediction:   l.ClientPrediction,
	}
}

func (x *UseItemTransactionData) ToLatest() *protocol.UseItemTransactionData {
	return &protocol.UseItemTransactionData{
		LegacyRequestID:    x.LegacyRequestID,
		LegacySetItemSlots: x.LegacySetItemSlots,
		Actions:            x.Actions,
		ActionType:         x.ActionType,
		TriggerType:        x.TriggerType,
		BlockPosition:      x.BlockPosition,
		BlockFace:          x.BlockFace,
		HotBarSlot:         x.HotBarSlot,
		HeldItem:           x.HeldItem,
		Position:           x.Position,
		ClickedPosition:    x.ClickedPosition,
		BlockRuntimeID:     x.BlockRuntimeID,
		ClientPrediction:   x.ClientPrediction,
	}
}

// Marshal ...
func (x *UseItemTransactionData) Marshal(r protocol.IO) {
	r.Varuint32(&x.ActionType)
	if IsProtoGTE(r, ID712) {
		r.Varuint32(&x.TriggerType)
	}
	r.UBlockPos(&x.BlockPosition)
	r.Varint32(&x.BlockFace)
	r.Varint32(&x.HotBarSlot)
	r.ItemInstance(&x.HeldItem)
	r.Vec3(&x.Position)
	r.Vec3(&x.ClickedPosition)
	r.Varuint32(&x.BlockRuntimeID)
	if IsProtoGTE(r, ID712) {
		r.Varuint32(&x.ClientPrediction)
	}
}
