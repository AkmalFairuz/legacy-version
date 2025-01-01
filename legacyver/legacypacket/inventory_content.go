package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// InventoryContent is sent by the server to update the full content of a particular inventory. It is usually
// sent for the main inventory of the player, but also works for other inventories that are currently opened
// by the player.
type InventoryContent struct {
	// WindowID is the ID that identifies one of the windows that the client currently has opened, or one of
	// the consistent windows such as the main inventory.
	WindowID uint32
	// Content is the new content of the inventory. The length of this slice must be equal to the full size of
	// the inventory window updated.
	Content []protocol.ItemInstance
	// Container is the protocol.FullContainerName that describes the container that the content is for.
	Container proto.FullContainerName
	// DynamicContainerSize ...
	DynamicContainerSize uint32
	// StorageItem is the item that is acting as the storage container for the inventory. If the inventory is
	// not a dynamic container then this field should be left empty. When set, only the item type is used by
	// the client and none of the other stack info.
	StorageItem protocol.ItemInstance
}

// ID ...
func (*InventoryContent) ID() uint32 {
	return packet.IDInventoryContent
}

func (pk *InventoryContent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	protocol.FuncSlice(io, &pk.Content, io.ItemInstance)
	if proto.IsProtoGTE(io, proto.ID729) {
		protocol.Single(io, &pk.Container)
	}
	if proto.IsProtoGTE(io, proto.ID748) {
		io.ItemInstance(&pk.StorageItem)
	} else {
		io.Varuint32(&pk.DynamicContainerSize)
	}
}
