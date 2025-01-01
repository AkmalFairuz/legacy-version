package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Transfer is sent by the server to transfer a player from the current server to another. Doing so will
// fully disconnect the client, bring it back to the main menu and make it connect to the next server.
type Transfer struct {
	// Address is the address of the new server, which might be either a hostname or an actual IP address.
	Address string
	// Port is the UDP port of the new server.
	Port uint16
	// ReloadWorld currently has an unknown usage.
	ReloadWorld bool
}

// ID ...
func (*Transfer) ID() uint32 {
	return packet.IDTransfer
}

func (pk *Transfer) Marshal(io protocol.IO) {
	io.String(&pk.Address)
	io.Uint16(&pk.Port)
	if proto.IsProtoGTE(io, proto.ID729) {
		io.Bool(&pk.ReloadWorld)
	}
}
