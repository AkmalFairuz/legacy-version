package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CodeBuilderSource is an Education Edition packet sent by the client to the server to run an operation with a
// code builder.
type CodeBuilderSource struct {
	// Operation is used to distinguish the operation performed. It is always one of the constants listed above.
	Operation byte
	// Category is used to distinguish the category of the operation performed. It is always one of the constants
	// listed above.
	Category byte
	// CodeStatus is the status of the code builder. It is always one of the constants listed above.
	CodeStatus byte
}

// ID ...
func (pk *CodeBuilderSource) ID() uint32 {
	return packet.IDCodeBuilderSource
}

func (pk *CodeBuilderSource) Marshal(io protocol.IO) {
	io.Uint8(&pk.Operation)
	io.Uint8(&pk.Category)
	if proto.IsProtoGTE(io, proto.ID685) {
		io.Uint8(&pk.CodeStatus)
	}
}
