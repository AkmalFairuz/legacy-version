package legacypacket

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// CameraInstruction gives a custom camera specific instructions to operate.
type CameraInstruction struct {
	// Set is a camera instruction that sets the camera to a specified preset.
	Set protocol.Optional[protocol.CameraInstructionSet]
	// Clear can be set to true to clear all the current camera instructions.
	Clear protocol.Optional[bool]
	// Fade is a camera instruction that fades the screen to a specified colour.
	Fade protocol.Optional[protocol.CameraInstructionFade]
	// Target is a camera instruction that targets a specific entity.
	Target protocol.Optional[protocol.CameraInstructionTarget]
	// RemoveTarget can be set to true to remove the current aim assist target.
	RemoveTarget protocol.Optional[bool]
}

// ID ...
func (*CameraInstruction) ID() uint32 {
	return packet.IDCameraInstruction
}

func (pk *CameraInstruction) Marshal(io protocol.IO) {
	protocol.OptionalMarshaler(io, &pk.Set)
	protocol.OptionalFunc(io, &pk.Clear, io.Bool)
	protocol.OptionalMarshaler(io, &pk.Fade)
	if proto.IsProtoGTE(io, proto.ID712) {
		protocol.OptionalMarshaler(io, &pk.Target)
		protocol.OptionalFunc(io, &pk.RemoveTarget, io.Bool)
	}
}
