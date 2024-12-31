package legacyver

import (
	"github.com/akmalfairuz/legacy-version/legacyver/legacypacket"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	packetPoolClient packet.Pool
	packetPoolServer packet.Pool
)

func init() {
	packetPoolClient = packet.NewClientPool()
	packetPoolServer = packet.NewServerPool()

	packetPoolServer[packet.IDItemStackResponse] = func() packet.Packet { return &legacypacket.ItemStackResponse{} }
	packetPoolServer[packet.IDResourcePacksInfo] = func() packet.Packet { return &legacypacket.ResourcePacksInfo{} }

	packetPoolClient[packet.IDPlayerAuthInput] = func() packet.Packet { return &legacypacket.PlayerAuthInput{} }
}

type Protocol struct {
	ver string
	id  int32

	blockTranslator BlockTranslator
	itemTranslator  ItemTranslator
}

func (p *Protocol) Ver() string {
	return p.ver
}

func (p *Protocol) ID() int32 {
	return p.id
}

func (p *Protocol) Packets(listener bool) packet.Pool {
	if listener {
		return packetPoolClient
	}
	return packetPoolServer
}

func (p *Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return proto.NewReader(protocol.NewReader(r, shieldID, enableLimits), p.id)
}

func (p *Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return proto.NewWriter(protocol.NewWriter(w, shieldID), p.id)
}

func (p *Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return p.blockTranslator.UpgradeBlockPackets(
		p.itemTranslator.UpgradeItemPackets(p.upgradePackets([]packet.Packet{pk}, conn), conn),
		conn)
}

func (p *Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	return p.downgradePackets(p.blockTranslator.DowngradeBlockPackets(
		p.itemTranslator.DowngradeItemPackets([]packet.Packet{pk}, conn),
		conn), conn)
}

func (p *Protocol) downgradePackets(pks []packet.Packet, conn *minecraft.Conn) []packet.Packet {
	for pkIndex, pk := range pks {
		switch pk := pk.(type) {
		case *packet.StartGame:
			pk.GameVersion = p.ver
			pk.BaseGameVersion = p.ver
		case *packet.PlayerAuthInput:
			inputData := pk.InputData
			if p.ID() < proto.ID766 {
				inputData = fitBitset(inputData, 64)
			}
			pks[pkIndex] = &legacypacket.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              inputData,
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				InteractPitch:          pk.InteractPitch,
				InteractYaw:            pk.InteractYaw,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				CameraOrientation:      pk.CameraOrientation,
				RawMoveVector:          pk.RawMoveVector,
			}
		case *packet.ItemStackResponse:
			responses := make([]proto.ItemStackResponse, len(pk.Responses))
			for i, r := range pk.Responses {
				containerInfo := make([]proto.StackResponseContainerInfo, len(r.ContainerInfo))
				for j, c := range r.ContainerInfo {
					slotInfo := make([]proto.StackResponseSlotInfo, len(c.SlotInfo))
					for k, s := range c.SlotInfo {
						slotInfo[k] = proto.StackResponseSlotInfo{
							Slot:                 s.Slot,
							HotbarSlot:           s.HotbarSlot,
							Count:                s.Count,
							StackNetworkID:       s.StackNetworkID,
							CustomName:           s.CustomName,
							FilteredCustomName:   s.FilteredCustomName,
							DurabilityCorrection: s.DurabilityCorrection,
						}
					}
					containerInfo[j] = proto.StackResponseContainerInfo{
						Container: c.Container,
						SlotInfo:  slotInfo,
					}
				}

				responses[i] = proto.ItemStackResponse{
					Status:        r.Status,
					RequestID:     r.RequestID,
					ContainerInfo: containerInfo,
				}
			}
			pks[pkIndex] = &legacypacket.ItemStackResponse{Responses: responses}
		case *packet.ResourcePacksInfo:
			texturePacks := make([]proto.TexturePackInfo, len(pk.TexturePacks))
			for i, t := range pk.TexturePacks {
				texturePacks[i] = proto.TexturePackInfo{
					UUID:            t.UUID,
					Version:         t.Version,
					Size:            t.Size,
					ContentKey:      t.ContentKey,
					SubPackName:     t.SubPackName,
					ContentIdentity: t.ContentIdentity,
					HasScripts:      t.HasScripts,
					AddonPack:       t.AddonPack,
					RTXEnabled:      t.RTXEnabled,
					DownloadURL:     t.DownloadURL,
				}
			}
			pks[pkIndex] = &legacypacket.ResourcePacksInfo{
				TexturePackRequired:  pk.TexturePackRequired,
				HasAddons:            pk.HasAddons,
				HasScripts:           pk.HasScripts,
				WorldTemplateUUID:    pk.WorldTemplateUUID,
				WorldTemplateVersion: pk.WorldTemplateVersion,
				TexturePacks:         texturePacks,
			}
		}
	}

	return pks
}

func (p *Protocol) upgradePackets(pks []packet.Packet, conn *minecraft.Conn) []packet.Packet {
	for pkIndex, pk := range pks {
		switch pk := pk.(type) {
		case *packet.StartGame:
			pk.GameVersion = p.ver
			pk.BaseGameVersion = p.ver
		case *legacypacket.PlayerAuthInput:
			pks[pkIndex] = &packet.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              fitBitset(pk.InputData, packet.PlayerAuthInputBitsetSize),
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				InteractPitch:          pk.InteractPitch,
				InteractYaw:            pk.InteractYaw,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				CameraOrientation:      pk.CameraOrientation,
				RawMoveVector:          pk.RawMoveVector,
			}
		case *legacypacket.ItemStackResponse:
			responses := make([]protocol.ItemStackResponse, len(pk.Responses))
			for i, r := range pk.Responses {
				responses[i] = r.ToLatest()
			}
			pks[pkIndex] = &packet.ItemStackResponse{Responses: responses}
		case *legacypacket.ResourcePacksInfo:
			texturePacks := make([]protocol.TexturePackInfo, len(pk.TexturePacks))
			for i, t := range pk.TexturePacks {
				texturePacks[i] = t.ToLatest()
			}
			pks[pkIndex] = &packet.ResourcePacksInfo{
				TexturePackRequired:  pk.TexturePackRequired,
				HasAddons:            pk.HasAddons,
				HasScripts:           pk.HasScripts,
				WorldTemplateUUID:    pk.WorldTemplateUUID,
				WorldTemplateVersion: pk.WorldTemplateVersion,
				TexturePacks:         texturePacks,
			}
		}
	}
	return pks
}
