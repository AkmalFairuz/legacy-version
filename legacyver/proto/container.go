package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// FullContainerName contains information required to identify a container in a StackRequestSlotInfo.
type FullContainerName struct {
	// ContainerID is the ID of the container that the slot was in.
	ContainerID byte
	// DynamicContainerID is the ID of the container if it is dynamic. If the container is not dynamic, this
	// field should be left empty. A non-optional value of 0 is assumed to be non-empty.
	DynamicContainerID protocol.Optional[uint32]
}

func (x *FullContainerName) FromLatest(v protocol.FullContainerName) FullContainerName {
	x.ContainerID = v.ContainerID
	x.DynamicContainerID = v.DynamicContainerID
	return *x
}

func (x *FullContainerName) ToLatest() protocol.FullContainerName {
	return protocol.FullContainerName{
		ContainerID:        x.ContainerID,
		DynamicContainerID: x.DynamicContainerID,
	}
}

func (x *FullContainerName) Marshal(r protocol.IO) {
	r.Uint8(&x.ContainerID)
	if IsProtoGTE(r, ID729) {
		protocol.OptionalFunc(r, &x.DynamicContainerID, r.Uint32)
	} else {
		dynamicContainerID, _ := x.DynamicContainerID.Value()
		r.Uint32(&dynamicContainerID)
		if dynamicContainerID != 0 {
			x.DynamicContainerID = protocol.Option(dynamicContainerID)
		}
	}
}
