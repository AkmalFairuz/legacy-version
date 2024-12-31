package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CameraPresets gives the client a list of custom camera presets.
type CameraPresets struct {
	// Presets is a list of camera presets that can be used by other cameras. The order of this list is important because
	// the index of presets is used as a pointer in the CameraInstruction packet.
	Presets []proto.CameraPreset
}

// ID ...
func (*CameraPresets) ID() uint32 {
	return packet.IDCameraPresets
}

func (pk *CameraPresets) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Presets)
}
