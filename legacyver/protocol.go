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
	case packet.IDContainerRegistryCleanup:
		return func() packet.Packet { return &legacypacket.ContainerRegistryCleanup{} }
	case packet.IDEmote:
		return func() packet.Packet { return &legacypacket.Emote{} }
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
	case packet.IDTransfer:
		return func() packet.Packet { return &legacypacket.Transfer{} }
	case packet.IDUpdateAttributes:
		return func() packet.Packet { return &legacypacket.UpdateAttributes{} }
	case packet.IDAddPlayer:
		return func() packet.Packet { return &legacypacket.AddPlayer{} }
	case packet.IDAddActor:
		return func() packet.Packet { return &legacypacket.AddActor{} }
	case packet.IDSetActorLink:
		return func() packet.Packet { return &legacypacket.SetActorLink{} }
	case packet.IDCameraInstruction:
		return func() packet.Packet { return &legacypacket.CameraInstruction{} }
	case packet.IDChangeDimension:
		return func() packet.Packet { return &legacypacket.ChangeDimension{} }
	case packet.IDCorrectPlayerMovePrediction:
		return func() packet.Packet { return &legacypacket.CorrectPlayerMovePrediction{} }
	case packet.IDDisconnect:
		return func() packet.Packet { return &legacypacket.Disconnect{} }
	case packet.IDEditorNetwork:
		return func() packet.Packet { return &legacypacket.EditorNetwork{} }
	case packet.IDMobArmourEquipment:
		return func() packet.Packet { return &legacypacket.MobArmourEquipment{} }
	case packet.IDPlayerArmourDamage:
		return func() packet.Packet { return &legacypacket.PlayerArmourDamage{} }
	case packet.IDSetTitle:
		return func() packet.Packet { return &legacypacket.SetTitle{} }
	case packet.IDStopSound:
		return func() packet.Packet { return &legacypacket.StopSound{} }
	case packet.IDInventoryTransaction:
		return func() packet.Packet { return &legacypacket.InventoryTransaction{} }
	case packet.IDItemStackRequest:
		return func() packet.Packet { return &legacypacket.ItemStackRequest{} }
	case packet.IDCraftingData:
		return func() packet.Packet { return &legacypacket.CraftingData{} }
	case packet.IDContainerClose:
		return func() packet.Packet { return &legacypacket.ContainerClose{} }
	case packet.IDText:
		return func() packet.Packet { return &legacypacket.Text{} }
	case packet.IDStartGame:
		return func() packet.Packet { return &legacypacket.StartGame{} }
	case packet.IDCodeBuilderSource:
		return func() packet.Packet { return &legacypacket.CodeBuilderSource{} }
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
		case *packet.ClientCacheStatus:
			pk.Enabled = false // TODO: enable when chunk translation is not broken
		case *packet.CameraPresets:
			presets := make([]proto.CameraPreset, len(pk.Presets))
			for i, p := range pk.Presets {
				presets[i] = (&proto.CameraPreset{}).FromLatest(p)
			}
			pks[pkIndex] = &legacypacket.CameraPresets{
				Presets: presets,
			}
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
				responses[i] = (&proto.ItemStackResponse{}).FromLatest(r)
			}
			pks[pkIndex] = &legacypacket.ItemStackResponse{Responses: responses}
		case *packet.ResourcePacksInfo:
			texturePacks := make([]proto.TexturePackInfo, len(pk.TexturePacks))
			packURLs := make([]protocol.PackURL, 0)
			for i, t := range pk.TexturePacks {
				texturePacks[i] = (&proto.TexturePackInfo{}).FromLatest(t)
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
				Container:            (&proto.FullContainerName{}).FromLatest(pk.Container),
				DynamicContainerSize: 0,
				StorageItem:          pk.StorageItem,
				NewItem:              pk.NewItem,
			}
		case *packet.InventoryContent:
			pks[pkIndex] = &legacypacket.InventoryContent{
				WindowID:             pk.WindowID,
				Content:              pk.Content,
				Container:            (&proto.FullContainerName{}).FromLatest(pk.Container),
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
		case *packet.UpdateAttributes:
			attributes := make([]proto.Attribute, len(pk.Attributes))
			for i, a := range pk.Attributes {
				attributes[i] = (&proto.Attribute{}).FromLatest(a)
			}
			pks[pkIndex] = &legacypacket.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		case *packet.ContainerRegistryCleanup:
			removedContainers := make([]proto.FullContainerName, len(pk.RemovedContainers))
			for i, c := range pk.RemovedContainers {
				removedContainers[i] = (&proto.FullContainerName{}).FromLatest(c)
			}
			pks[pkIndex] = &legacypacket.ContainerRegistryCleanup{
				RemovedContainers: removedContainers,
			}
		case *packet.Emote:
			pks[pkIndex] = &legacypacket.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteLength:     pk.EmoteLength,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *packet.Transfer:
			pks[pkIndex] = &legacypacket.Transfer{
				Address:     pk.Address,
				Port:        pk.Port,
				ReloadWorld: pk.ReloadWorld,
			}
		case *packet.AddActor:
			links := make([]proto.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = (&proto.EntityLink{}).FromLatest(l)
			}
			pks[pkIndex] = &legacypacket.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      links,
			}
		case *packet.AddPlayer:
			links := make([]proto.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = (&proto.EntityLink{}).FromLatest(l)
			}
			pks[pkIndex] = &legacypacket.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      pk.AbilityData,
				EntityLinks:      links,
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			}
		case *packet.SetActorLink:
			pks[pkIndex] = &legacypacket.SetActorLink{
				EntityLink: (&proto.EntityLink{}).FromLatest(pk.EntityLink),
			}
		case *packet.CameraInstruction:
			pks[pkIndex] = &legacypacket.CameraInstruction{
				Set:          pk.Set,
				Clear:        pk.Clear,
				Fade:         pk.Fade,
				Target:       pk.Target,
				RemoveTarget: pk.RemoveTarget,
			}
		case *packet.ChangeDimension:
			pks[pkIndex] = &legacypacket.ChangeDimension{
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				Respawn:         pk.Respawn,
				LoadingScreenID: pk.LoadingScreenID,
			}
		case *packet.CorrectPlayerMovePrediction:
			pks[pkIndex] = &legacypacket.CorrectPlayerMovePrediction{
				PredictionType:         pk.PredictionType,
				Position:               pk.Position,
				Delta:                  pk.Delta,
				Rotation:               pk.Rotation,
				VehicleAngularVelocity: pk.VehicleAngularVelocity,
				OnGround:               pk.OnGround,
				Tick:                   pk.Tick,
			}
		case *packet.Disconnect:
			pks[pkIndex] = &legacypacket.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
				FilteredMessage:         pk.FilteredMessage,
			}
		case *packet.EditorNetwork:
			pks[pkIndex] = &legacypacket.EditorNetwork{
				RouteToManager: pk.RouteToManager,
				Payload:        pk.Payload,
			}
		case *packet.MobArmourEquipment:
			pks[pkIndex] = &legacypacket.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
				Body:            pk.Body,
			}
		case *packet.PlayerArmourDamage:
			pks[pkIndex] = &legacypacket.PlayerArmourDamage{
				Bitset:           pk.Bitset,
				HelmetDamage:     pk.HelmetDamage,
				ChestplateDamage: pk.ChestplateDamage,
				LeggingsDamage:   pk.LeggingsDamage,
				BootsDamage:      pk.BootsDamage,
				BodyDamage:       pk.BodyDamage,
			}
		case *packet.SetTitle:
			pks[pkIndex] = &legacypacket.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *packet.StopSound:
			pks[pkIndex] = &legacypacket.StopSound{
				SoundName:       pk.SoundName,
				StopAll:         pk.StopAll,
				StopMusicLegacy: pk.StopMusicLegacy,
			}
		case *packet.InventoryTransaction:
			trData := pk.TransactionData
			if x, ok := trData.(*protocol.UseItemTransactionData); ok {
				trData = (&proto.UseItemTransactionData{}).FromLatest(x)
			}
			pks[pkIndex] = &legacypacket.InventoryTransaction{
				LegacyRequestID:    pk.LegacyRequestID,
				LegacySetItemSlots: pk.LegacySetItemSlots,
				Actions:            pk.Actions,
				TransactionData:    trData,
			}
		case *packet.ItemStackRequest:
			requests := make([]proto.ItemStackRequest, len(pk.Requests))
			for i, r := range pk.Requests {
				requests[i] = (&proto.ItemStackRequest{}).FromLatest(r)
			}
			pks[pkIndex] = &legacypacket.ItemStackRequest{Requests: requests}
		case *packet.Text:
			pks[pkIndex] = &legacypacket.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *packet.CraftingData:
			recipes := make([]proto.Recipe, len(pk.Recipes))
			for i, r := range pk.Recipes {
				recipes[i] = proto.RecipeFromLatest(r)
			}
			pks[pkIndex] = &legacypacket.CraftingData{
				Recipes:                      recipes,
				PotionRecipes:                pk.PotionRecipes,
				PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
				MaterialReducers:             pk.MaterialReducers,
				ClearRecipes:                 pk.ClearRecipes,
			}
		case *packet.StartGame:
			// Adjust game version
			pk.GameVersion = p.ver
			pk.BaseGameVersion = p.ver

			pks[pkIndex] = &legacypacket.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Hardcore:                       pk.Hardcore,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         pk.PlayerMovementSettings,
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				Items:                          pk.Items,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				ServerID:                       pk.ServerID,
				WorldID:                        pk.WorldID,
				ScenarioID:                     pk.ScenarioID,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			}
		case *packet.CodeBuilderSource:
			pks[pkIndex] = &legacypacket.CodeBuilderSource{
				Operation:  pk.Operation,
				Category:   pk.Category,
				CodeStatus: pk.CodeStatus,
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
				Container:   pk.Container.ToLatest(),
				StorageItem: pk.StorageItem,
				NewItem:     pk.NewItem,
			}
		case *legacypacket.InventoryContent:
			pks[pkIndex] = &packet.InventoryContent{
				WindowID:    pk.WindowID,
				Content:     pk.Content,
				Container:   pk.Container.ToLatest(),
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
		case *legacypacket.UpdateAttributes:
			attributes := make([]protocol.Attribute, len(pk.Attributes))
			for i, a := range pk.Attributes {
				attributes[i] = a.ToLatest()
			}
			pks[pkIndex] = &packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes:      attributes,
				Tick:            pk.Tick,
			}
		case *legacypacket.ContainerRegistryCleanup:
			removedContainers := make([]protocol.FullContainerName, len(pk.RemovedContainers))
			for i, c := range pk.RemovedContainers {
				removedContainers[i] = c.ToLatest()
			}
			pks[pkIndex] = &packet.ContainerRegistryCleanup{
				RemovedContainers: removedContainers,
			}
		case *legacypacket.Emote:
			pks[pkIndex] = &packet.Emote{
				EntityRuntimeID: pk.EntityRuntimeID,
				EmoteLength:     pk.EmoteLength,
				EmoteID:         pk.EmoteID,
				XUID:            pk.XUID,
				PlatformID:      pk.PlatformID,
				Flags:           pk.Flags,
			}
		case *legacypacket.Transfer:
			pks[pkIndex] = &packet.Transfer{
				Address:     pk.Address,
				Port:        pk.Port,
				ReloadWorld: pk.ReloadWorld,
			}
		case *legacypacket.AddActor:
			links := make([]protocol.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = l.ToLatest()
			}
			pks[pkIndex] = &packet.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      links,
			}
		case *legacypacket.AddPlayer:
			links := make([]protocol.EntityLink, len(pk.EntityLinks))
			for i, l := range pk.EntityLinks {
				links[i] = l.ToLatest()
			}
			pks[pkIndex] = &packet.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      pk.AbilityData,
				EntityLinks:      links,
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			}
		case *legacypacket.SetActorLink:
			pks[pkIndex] = &packet.SetActorLink{
				EntityLink: pk.EntityLink.ToLatest(),
			}
		case *legacypacket.CameraInstruction:
			pks[pkIndex] = &packet.CameraInstruction{
				Set:          pk.Set,
				Clear:        pk.Clear,
				Fade:         pk.Fade,
				Target:       pk.Target,
				RemoveTarget: pk.RemoveTarget,
			}
		case *legacypacket.ChangeDimension:
			pks[pkIndex] = &packet.ChangeDimension{
				Dimension:       pk.Dimension,
				Position:        pk.Position,
				Respawn:         pk.Respawn,
				LoadingScreenID: pk.LoadingScreenID,
			}
		case *legacypacket.CorrectPlayerMovePrediction:
			pks[pkIndex] = &packet.CorrectPlayerMovePrediction{
				PredictionType:         pk.PredictionType,
				Position:               pk.Position,
				Delta:                  pk.Delta,
				Rotation:               pk.Rotation,
				VehicleAngularVelocity: pk.VehicleAngularVelocity,
				OnGround:               pk.OnGround,
				Tick:                   pk.Tick,
			}
		case *legacypacket.Disconnect:
			pks[pkIndex] = &packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
				FilteredMessage:         pk.FilteredMessage,
			}
		case *legacypacket.EditorNetwork:
			pks[pkIndex] = &packet.EditorNetwork{
				RouteToManager: pk.RouteToManager,
				Payload:        pk.Payload,
			}
		case *legacypacket.MobArmourEquipment:
			pks[pkIndex] = &packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
				Body:            pk.Body,
			}
		case *legacypacket.PlayerArmourDamage:
			pks[pkIndex] = &packet.PlayerArmourDamage{
				Bitset:           pk.Bitset,
				HelmetDamage:     pk.HelmetDamage,
				ChestplateDamage: pk.ChestplateDamage,
				LeggingsDamage:   pk.LeggingsDamage,
				BootsDamage:      pk.BootsDamage,
				BodyDamage:       pk.BodyDamage,
			}
		case *legacypacket.SetTitle:
			pks[pkIndex] = &packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *legacypacket.StopSound:
			pks[pkIndex] = &packet.StopSound{
				SoundName:       pk.SoundName,
				StopAll:         pk.StopAll,
				StopMusicLegacy: pk.StopMusicLegacy,
			}
		case *legacypacket.InventoryTransaction:
			trData := pk.TransactionData
			if x, ok := trData.(*proto.UseItemTransactionData); ok {
				trData = x.ToLatest()
			}
			pks[pkIndex] = &packet.InventoryTransaction{
				LegacyRequestID:    pk.LegacyRequestID,
				LegacySetItemSlots: pk.LegacySetItemSlots,
				Actions:            pk.Actions,
				TransactionData:    trData,
			}
		case *legacypacket.ItemStackRequest:
			requests := make([]protocol.ItemStackRequest, len(pk.Requests))
			for i, r := range pk.Requests {
				requests[i] = r.ToLatest()
			}
			pks[pkIndex] = &packet.ItemStackRequest{Requests: requests}
		case *legacypacket.CraftingData:
			recipes := make([]protocol.Recipe, len(pk.Recipes))
			for i, r := range pk.Recipes {
				recipes[i] = proto.RecipeToLatest(r)
			}
			pks[pkIndex] = &packet.CraftingData{
				Recipes:                      recipes,
				PotionRecipes:                pk.PotionRecipes,
				PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
				MaterialReducers:             pk.MaterialReducers,
				ClearRecipes:                 pk.ClearRecipes,
			}
		case *legacypacket.ContainerClose:
			pks[pkIndex] = &packet.ContainerClose{
				WindowID:      pk.WindowID,
				ContainerType: pk.ContainerType,
				ServerSide:    pk.ServerSide,
			}
		case *legacypacket.Text:
			pks[pkIndex] = &packet.Text{
				TextType:         pk.TextType,
				NeedsTranslation: pk.NeedsTranslation,
				SourceName:       pk.SourceName,
				Message:          pk.Message,
				Parameters:       pk.Parameters,
				XUID:             pk.XUID,
				PlatformChatID:   pk.PlatformChatID,
				FilteredMessage:  pk.FilteredMessage,
			}
		case *legacypacket.StartGame:
			pks[pkIndex] = &packet.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Hardcore:                       pk.Hardcore,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      pk.ForceExperimentalGameplay,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         pk.PlayerMovementSettings,
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				Items:                          pk.Items,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				ServerID:                       pk.ServerID,
				WorldID:                        pk.WorldID,
				ScenarioID:                     pk.ScenarioID,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			}
		case *legacypacket.CodeBuilderSource:
			pks[pkIndex] = &packet.CodeBuilderSource{
				Operation:  pk.Operation,
				Category:   pk.Category,
				CodeStatus: pk.CodeStatus,
			}
		}
	}
	return pks
}
