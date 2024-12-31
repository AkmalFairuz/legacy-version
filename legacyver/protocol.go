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

	for pkId, cur := range packetPoolClient {
		packetPoolClient[pkId] = convertPacketFunc(pkId, cur)
	}
}

func convertPacketFunc(pid uint32, cur func() packet.Packet) func() packet.Packet {
	switch pid {
	case packet.IDCameraAimAssist:
		return func() packet.Packet { return &legacypacket.CameraAimAssist{} }
	case packet.IDCameraPresets:
		return func() packet.Packet { return &legacypacket.CameraPresets{} }
	case packet.IDInventoryContent:
		return func() packet.Packet { return &legacypacket.InventoryContent{} }
	case packet.IDInventorySlot:
		return func() packet.Packet { return &legacypacket.InventorySlot{} }
	case packet.IDItemStackResponse:
		return func() packet.Packet { return &legacypacket.ItemStackResponse{} }
	case packet.IDMobEffect:
		return func() packet.Packet { return &legacypacket.MobEffect{} }
	case packet.IDPlayerAuthInput:
		return func() packet.Packet { return &legacypacket.PlayerAuthInput{} }
	case packet.IDResourcePacksInfo:
		return func() packet.Packet { return &legacypacket.ResourcePacksInfo{} }
	default:
		return cur
	}
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
		case *packet.CameraPresets:
			presets := make([]proto.CameraPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i].FromLatest(p)
			}
			pks[pkIndex] = &legacypacket.CameraPresets{
				Presets: presets,
			}
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
			packURLs := make([]protocol.PackURL, 0)
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
				if t.DownloadURL != "" {
					packURLs = append(packURLs, protocol.PackURL{
						UUIDVersion: t.UUID.String() + "_" + t.Version,
						URL:         t.DownloadURL,
					})
				}
			}
			pks[pkIndex] = &legacypacket.ResourcePacksInfo{
				TexturePackRequired:  pk.TexturePackRequired,
				HasAddons:            pk.HasAddons,
				HasScripts:           pk.HasScripts,
				WorldTemplateUUID:    pk.WorldTemplateUUID,
				WorldTemplateVersion: pk.WorldTemplateVersion,
				TexturePacks:         texturePacks,
				PackURLs:             packURLs,
			}
		case *packet.InventorySlot:
			pks[pkIndex] = &legacypacket.InventorySlot{
				WindowID:             pk.WindowID,
				Slot:                 pk.Slot,
				Container:            pk.Container,
				DynamicContainerSize: 0,
				StorageItem:          pk.StorageItem,
				NewItem:              pk.NewItem,
			}
		case *packet.InventoryContent:
			pks[pkIndex] = &legacypacket.InventoryContent{
				WindowID:             pk.WindowID,
				Content:              pk.Content,
				Container:            pk.Container,
				DynamicContainerSize: 0,
				StorageItem:          pk.StorageItem,
			}
		case *packet.MobEffect:
			pks[pkIndex] = &legacypacket.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
			}
		case *packet.CameraAimAssist:
			pks[pkIndex] = &legacypacket.CameraAimAssist{
				Preset:     pk.Preset,
				Angle:      pk.Angle,
				Distance:   pk.Distance,
				TargetMode: pk.TargetMode,
				Action:     pk.Action,
			}
		}
	}

	return pks
}

func (p *Protocol) upgradePackets(pks []packet.Packet, conn *minecraft.Conn) []packet.Packet {
	for pkIndex, pk := range pks {
		switch pk := pk.(type) {
		case *packet.ClientCacheStatus:
			pk.Enabled = false // TODO: enable when chunk translation is not broken
		case *legacypacket.CameraPresets:
			presets := make([]protocol.CameraPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = p.ToLatest()
			}
			pks[pkIndex] = &packet.CameraPresets{
				Presets: presets,
			}
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
				if texturePacks[i].DownloadURL == "" {
					for _, u := range pk.PackURLs {
						if u.UUIDVersion == t.UUID.String()+"_"+t.Version {
							texturePacks[i].DownloadURL = u.URL
							break
						}
					}
				}
			}
			pks[pkIndex] = &packet.ResourcePacksInfo{
				TexturePackRequired:  pk.TexturePackRequired,
				HasAddons:            pk.HasAddons,
				HasScripts:           pk.HasScripts,
				WorldTemplateUUID:    pk.WorldTemplateUUID,
				WorldTemplateVersion: pk.WorldTemplateVersion,
				TexturePacks:         texturePacks,
			}
		case *legacypacket.InventorySlot:
			pks[pkIndex] = &packet.InventorySlot{
				WindowID:    pk.WindowID,
				Slot:        pk.Slot,
				Container:   pk.Container,
				StorageItem: pk.StorageItem,
				NewItem:     pk.NewItem,
			}
		case *legacypacket.InventoryContent:
			pks[pkIndex] = &packet.InventoryContent{
				WindowID:    pk.WindowID,
				Content:     pk.Content,
				Container:   pk.Container,
				StorageItem: pk.StorageItem,
			}
		case *legacypacket.MobEffect:
			pks[pkIndex] = &packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
				Tick:            pk.Tick,
			}
		case *legacypacket.CameraAimAssist:
			pks[pkIndex] = &packet.CameraAimAssist{
				Preset:     pk.Preset,
				Angle:      pk.Angle,
				Distance:   pk.Distance,
				TargetMode: pk.TargetMode,
				Action:     pk.Action,
			}
		}
	}
	return pks
}
