package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// EntityLink is a link between two entities, typically being one entity riding another.
type EntityLink struct {
	// RiddenEntityUniqueID is the entity unique ID of the entity that is being ridden. For a player sitting
	// in a boat, this is the unique ID of the boat.
	RiddenEntityUniqueID int64
	// RiderEntityUniqueID is the entity unique ID of the entity that is riding. For a player sitting in a
	// boat, this is the unique ID of the player.
	RiderEntityUniqueID int64
	// Type is one of the types above. It specifies the way the entity is linked to another entity.
	Type byte
	// Immediate is set to immediately dismount an entity from another. This should be set when the mount of
	// an entity is killed.
	Immediate bool
	// RiderInitiated specifies if the link was created by the rider, for example the player starting to ride
	// a horse by itself. This is generally true in vanilla environment for players.
	RiderInitiated bool
	// VehicleAngularVelocity is the angular velocity of the vehicle that the rider is riding.
	VehicleAngularVelocity float32
}

func (x *EntityLink) FromLatest(link protocol.EntityLink) EntityLink {
	x.RiddenEntityUniqueID = link.RiddenEntityUniqueID
	x.RiderEntityUniqueID = link.RiderEntityUniqueID
	x.Type = link.Type
	x.Immediate = link.Immediate
	x.RiderInitiated = link.RiderInitiated
	x.VehicleAngularVelocity = link.VehicleAngularVelocity
	return *x
}

func (x *EntityLink) ToLatest() protocol.EntityLink {
	return protocol.EntityLink{
		RiddenEntityUniqueID:   x.RiddenEntityUniqueID,
		RiderEntityUniqueID:    x.RiderEntityUniqueID,
		Type:                   x.Type,
		Immediate:              x.Immediate,
		RiderInitiated:         x.RiderInitiated,
		VehicleAngularVelocity: x.VehicleAngularVelocity,
	}
}

// Marshal encodes/decodes a single entity link.
func (x *EntityLink) Marshal(r protocol.IO) {
	r.Varint64(&x.RiddenEntityUniqueID)
	r.Varint64(&x.RiderEntityUniqueID)
	r.Uint8(&x.Type)
	r.Bool(&x.Immediate)
	r.Bool(&x.RiderInitiated)
	if IsProtoGTE(r, ID712) {
		r.Float32(&x.VehicleAngularVelocity)
	}
}
