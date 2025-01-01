package proto

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	FilterCauseServerChatPublic = iota
	FilterCauseServerChatWhisper
	FilterCauseSignText
	FilterCauseAnvilText
	FilterCauseBookAndQuillText
	FilterCauseCommandBlockText
	FilterCauseBlockActorDataText
	FilterCauseJoinEventText
	FilterCauseLeaveEventText
	FilterCauseSlashCommandChat
	FilterCauseCartographyText
	FilterCauseKickCommand
	FilterCauseTitleCommand
	FilterCauseSummonCommand
)

// ItemStackRequest represents a single request present in an ItemStackRequest packet sent by the client to
// change an item in an inventory.
// Item stack requests are either approved or rejected by the server using the ItemStackResponse packet.
type ItemStackRequest struct {
	// RequestID is a unique ID for the request. This ID is used by the server to send a response for this
	// specific request in the ItemStackResponse packet.
	RequestID int32
	// Actions is a list of actions performed by the client. The actual type of the actions depends on which
	// ID was present, and is one of the concrete types below.
	Actions []protocol.StackRequestAction
	// FilterStrings is a list of filter strings involved in the request. This is typically filled with one string
	// when an anvil or cartography is used.
	FilterStrings []string
	// FilterCause represents the cause of any potential filtering. This is one of the constants above.
	FilterCause int32
}

func (x *ItemStackRequest) FromLatest(y protocol.ItemStackRequest) ItemStackRequest {
	x.RequestID = y.RequestID
	x.Actions = make([]protocol.StackRequestAction, len(y.Actions))
	for i, v := range y.Actions {
		if z, ok := v.(*protocol.CraftRecipeStackRequestAction); ok {
			x.Actions[i] = (&CraftRecipeStackRequestAction{}).FromLatest(z)
			continue
		}

		if z, ok := v.(*protocol.AutoCraftRecipeStackRequestAction); ok {
			x.Actions[i] = (&AutoCraftRecipeStackRequestAction{}).FromLatest(z)
			continue
		}

		if z, ok := v.(*protocol.CraftCreativeStackRequestAction); ok {
			x.Actions[i] = (&CraftCreativeStackRequestAction{}).FromLatest(z)
			continue
		}

		if z, ok := v.(*protocol.CraftGrindstoneRecipeStackRequestAction); ok {
			x.Actions[i] = (&CraftGrindstoneRecipeStackRequestAction{}).FromLatest(z)
			continue
		}

		if z, ok := v.(*protocol.CraftLoomRecipeStackRequestAction); ok {
			x.Actions[i] = (&CraftLoomRecipeStackRequestAction{}).FromLatest(z)
			continue
		}
	}
	x.FilterStrings = y.FilterStrings
	x.FilterCause = y.FilterCause
	return *x
}

func (x *ItemStackRequest) ToLatest() protocol.ItemStackRequest {
	ret := protocol.ItemStackRequest{
		RequestID:     x.RequestID,
		Actions:       make([]protocol.StackRequestAction, len(x.Actions)),
		FilterStrings: x.FilterStrings,
		FilterCause:   x.FilterCause,
	}
	for i, v := range x.Actions {
		if z, ok := v.(*CraftRecipeStackRequestAction); ok {
			ret.Actions[i] = z.ToLatest()
			continue
		}

		if z, ok := v.(*AutoCraftRecipeStackRequestAction); ok {
			ret.Actions[i] = z.ToLatest()
			continue
		}

		if z, ok := v.(*CraftCreativeStackRequestAction); ok {
			ret.Actions[i] = z.ToLatest()
			continue
		}

		if z, ok := v.(*CraftGrindstoneRecipeStackRequestAction); ok {
			ret.Actions[i] = z.ToLatest()
			continue
		}

		if z, ok := v.(*CraftLoomRecipeStackRequestAction); ok {
			ret.Actions[i] = z.ToLatest()
			continue
		}
	}
	return ret
}

// Marshal encodes/decodes an ItemStackRequest.
func (x *ItemStackRequest) Marshal(r protocol.IO) {
	r.Varint32(&x.RequestID)
	protocol.FuncSlice(r, &x.Actions, func(p *protocol.StackRequestAction) {
		IOStackRequestAction(r, p)
	})
	protocol.FuncSlice(r, &x.FilterStrings, r.String)
	r.Int32(&x.FilterCause)
}

// lookupStackRequestActionType looks up the ID of a StackRequestAction.
func lookupStackRequestActionType(x StackRequestAction, id *uint8) bool {
	switch x.(type) {
	case *protocol.TakeStackRequestAction:
		*id = protocol.StackRequestActionTake
	case *protocol.PlaceStackRequestAction:
		*id = protocol.StackRequestActionPlace
	case *protocol.SwapStackRequestAction:
		*id = protocol.StackRequestActionSwap
	case *protocol.DropStackRequestAction:
		*id = protocol.StackRequestActionDrop
	case *protocol.DestroyStackRequestAction:
		*id = protocol.StackRequestActionDestroy
	case *protocol.ConsumeStackRequestAction:
		*id = protocol.StackRequestActionConsume
	case *protocol.CreateStackRequestAction:
		*id = protocol.StackRequestActionCreate
	case *protocol.LabTableCombineStackRequestAction:
		*id = protocol.StackRequestActionLabTableCombine
	case *protocol.BeaconPaymentStackRequestAction:
		*id = protocol.StackRequestActionBeaconPayment
	case *protocol.MineBlockStackRequestAction:
		*id = protocol.StackRequestActionMineBlock
	case *CraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipe
	case *AutoCraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeAuto
	case *CraftCreativeStackRequestAction:
		*id = protocol.StackRequestActionCraftCreative
	case *protocol.CraftRecipeOptionalStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeOptional
	case *CraftGrindstoneRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftGrindstone
	case *CraftLoomRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftLoom
	case *protocol.CraftNonImplementedStackRequestAction:
		*id = protocol.StackRequestActionCraftNonImplementedDeprecated
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		*id = protocol.StackRequestActionCraftResultsDeprecated
	default:
		return false
	}
	return true
}

// lookupStackRequestAction looks up the StackRequestAction matching an ID.
func lookupStackRequestAction(id uint8, x *protocol.StackRequestAction) bool {
	switch id {
	case protocol.StackRequestActionTake:
		*x = &protocol.TakeStackRequestAction{}
	case protocol.StackRequestActionPlace:
		*x = &protocol.PlaceStackRequestAction{}
	case protocol.StackRequestActionSwap:
		*x = &protocol.SwapStackRequestAction{}
	case protocol.StackRequestActionDrop:
		*x = &protocol.DropStackRequestAction{}
	case protocol.StackRequestActionDestroy:
		*x = &protocol.DestroyStackRequestAction{}
	case protocol.StackRequestActionConsume:
		*x = &protocol.ConsumeStackRequestAction{}
	case protocol.StackRequestActionCreate:
		*x = &protocol.CreateStackRequestAction{}
	case protocol.StackRequestActionPlaceInContainer:
		*x = &protocol.PlaceInContainerStackRequestAction{}
	case protocol.StackRequestActionTakeOutContainer:
		*x = &protocol.TakeOutContainerStackRequestAction{}
	case protocol.StackRequestActionLabTableCombine:
		*x = &protocol.LabTableCombineStackRequestAction{}
	case protocol.StackRequestActionBeaconPayment:
		*x = &protocol.BeaconPaymentStackRequestAction{}
	case protocol.StackRequestActionMineBlock:
		*x = &protocol.MineBlockStackRequestAction{}
	case protocol.StackRequestActionCraftRecipe:
		*x = &CraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeAuto:
		*x = &AutoCraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftCreative:
		*x = &CraftCreativeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeOptional:
		*x = &protocol.CraftRecipeOptionalStackRequestAction{}
	case protocol.StackRequestActionCraftGrindstone:
		*x = &CraftGrindstoneRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftLoom:
		*x = &CraftLoomRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftNonImplementedDeprecated:
		*x = &protocol.CraftNonImplementedStackRequestAction{}
	case protocol.StackRequestActionCraftResultsDeprecated:
		*x = &protocol.CraftResultsDeprecatedStackRequestAction{}
	default:
		return false
	}
	return true
}

const (
	ItemStackResponseStatusOK = iota
	ItemStackResponseStatusError
	ItemStackResponseStatusInvalidRequestActionType
	ItemStackResponseStatusActionRequestNotAllowed
	ItemStackResponseStatusScreenHandlerEndRequestFailed
	ItemStackResponseStatusItemRequestActionHandlerCommitFailed
	ItemStackResponseStatusInvalidRequestCraftActionType
	ItemStackResponseStatusInvalidCraftRequest
	ItemStackResponseStatusInvalidCraftRequestScreen
	ItemStackResponseStatusInvalidCraftResult
	ItemStackResponseStatusInvalidCraftResultIndex
	ItemStackResponseStatusInvalidCraftResultItem
	ItemStackResponseStatusInvalidItemNetId
	ItemStackResponseStatusMissingCreatedOutputContainer
	ItemStackResponseStatusFailedToSetCreatedItemOutputSlot
	ItemStackResponseStatusRequestAlreadyInProgress
	ItemStackResponseStatusFailedToInitSparseContainer
	ItemStackResponseStatusResultTransferFailed
	ItemStackResponseStatusExpectedItemSlotNotFullyConsumed
	ItemStackResponseStatusExpectedAnywhereItemNotFullyConsumed
	ItemStackResponseStatusItemAlreadyConsumedFromSlot
	ItemStackResponseStatusConsumedTooMuchFromSlot
	ItemStackResponseStatusMismatchSlotExpectedConsumedItem
	ItemStackResponseStatusMismatchSlotExpectedConsumedItemNetIdVariant
	ItemStackResponseStatusFailedToMatchExpectedSlotConsumedItem
	ItemStackResponseStatusFailedToMatchExpectedAllowedAnywhereConsumedItem
	ItemStackResponseStatusConsumedItemOutOfAllowedSlotRange
	ItemStackResponseStatusConsumedItemNotAllowed
	ItemStackResponseStatusPlayerNotInCreativeMode
	ItemStackResponseStatusInvalidExperimentalRecipeRequest
	ItemStackResponseStatusFailedToCraftCreative
	ItemStackResponseStatusFailedToGetLevelRecipe
	ItemStackResponseStatusFailedToFindRecipeByNetId
	ItemStackResponseStatusMismatchedCraftingSize
	ItemStackResponseStatusMissingInputSparseContainer
	ItemStackResponseStatusMismatchedRecipeForInputGridItems
	ItemStackResponseStatusEmptyCraftResults
	ItemStackResponseStatusFailedToEnchant
	ItemStackResponseStatusMissingInputItem
	ItemStackResponseStatusInsufficientPlayerLevelToEnchant
	ItemStackResponseStatusMissingMaterialItem
	ItemStackResponseStatusMissingActor
	ItemStackResponseStatusUnknownPrimaryEffect
	ItemStackResponseStatusPrimaryEffectOutOfRange
	ItemStackResponseStatusPrimaryEffectUnavailable
	ItemStackResponseStatusSecondaryEffectOutOfRange
	ItemStackResponseStatusSecondaryEffectUnavailable
	ItemStackResponseStatusDstContainerEqualToCreatedOutputContainer
	ItemStackResponseStatusDstContainerAndSlotEqualToSrcContainerAndSlot
	ItemStackResponseStatusFailedToValidateSrcSlot
	ItemStackResponseStatusFailedToValidateDstSlot
	ItemStackResponseStatusInvalidAdjustedAmount
	ItemStackResponseStatusInvalidItemSetType
	ItemStackResponseStatusInvalidTransferAmount
	ItemStackResponseStatusCannotSwapItem
	ItemStackResponseStatusCannotPlaceItem
	ItemStackResponseStatusUnhandledItemSetType
	ItemStackResponseStatusInvalidRemovedAmount
	ItemStackResponseStatusInvalidRegion
	ItemStackResponseStatusCannotDropItem
	ItemStackResponseStatusCannotDestroyItem
	ItemStackResponseStatusInvalidSourceContainer
	ItemStackResponseStatusItemNotConsumed
	ItemStackResponseStatusInvalidNumCrafts
	ItemStackResponseStatusInvalidCraftResultStackSize
	ItemStackResponseStatusCannotRemoveItem
	ItemStackResponseStatusCannotConsumeItem
	ItemStackResponseStatusScreenStackError
)

// ItemStackResponse is a response to an individual ItemStackRequest.
type ItemStackResponse struct {
	// Status specifies if the request with the RequestID below was successful. If this is the case, the
	// ContainerInfo below will have information on what slots ended up changing. If not, the container info
	// will be empty.
	// A non-0 status means an error occurred and will result in the action being reverted.
	Status uint8
	// RequestID is the unique ID of the request that this response is in reaction to. If rejected, the client
	// will undo the actions from the request with this ID.
	RequestID int32
	// ContainerInfo holds information on the containers that had their contents changed as a result of the
	// request.
	ContainerInfo []StackResponseContainerInfo
}

func (x *ItemStackResponse) FromLatest(y protocol.ItemStackResponse) ItemStackResponse {
	ret := ItemStackResponse{
		Status:        y.Status,
		RequestID:     y.RequestID,
		ContainerInfo: make([]StackResponseContainerInfo, len(y.ContainerInfo)),
	}
	for i, v := range y.ContainerInfo {
		ret.ContainerInfo[i] = StackResponseContainerInfo{
			Container: (&FullContainerName{}).FromLatest(v.Container),
			SlotInfo:  make([]StackResponseSlotInfo, len(v.SlotInfo)),
		}
		for j, w := range v.SlotInfo {
			ret.ContainerInfo[i].SlotInfo[j] = StackResponseSlotInfo{
				Slot:                 w.Slot,
				HotbarSlot:           w.HotbarSlot,
				Count:                w.Count,
				StackNetworkID:       w.StackNetworkID,
				CustomName:           w.CustomName,
				FilteredCustomName:   w.FilteredCustomName,
				DurabilityCorrection: w.DurabilityCorrection,
			}
		}
	}
	return ret
}

func (x *ItemStackResponse) ToLatest() protocol.ItemStackResponse {
	ret := protocol.ItemStackResponse{
		Status:        x.Status,
		RequestID:     x.RequestID,
		ContainerInfo: make([]protocol.StackResponseContainerInfo, len(x.ContainerInfo)),
	}
	for i, v := range x.ContainerInfo {
		ret.ContainerInfo[i] = protocol.StackResponseContainerInfo{
			Container: v.Container.ToLatest(),
			SlotInfo:  make([]protocol.StackResponseSlotInfo, len(v.SlotInfo)),
		}
		for j, w := range v.SlotInfo {
			ret.ContainerInfo[i].SlotInfo[j] = protocol.StackResponseSlotInfo{
				Slot:                 w.Slot,
				HotbarSlot:           w.HotbarSlot,
				Count:                w.Count,
				StackNetworkID:       w.StackNetworkID,
				CustomName:           w.CustomName,
				FilteredCustomName:   w.FilteredCustomName,
				DurabilityCorrection: w.DurabilityCorrection,
			}
		}
	}
	return ret
}

// Marshal encodes/decodes an ItemStackResponse.
func (x *ItemStackResponse) Marshal(r protocol.IO) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	if x.Status == ItemStackResponseStatusOK {
		protocol.Slice(r, &x.ContainerInfo)
	}
}

// StackResponseContainerInfo holds information on what slots in a container have what item stack in them.
type StackResponseContainerInfo struct {
	// Container is the FullContainerName that describes the container that the slots that follow are in. For
	// the main inventory, the ContainerID seems to be 0x1b. Fur the cursor, this value seems to be 0x3a. For
	// the crafting grid, this value seems to be 0x0d.
	Container FullContainerName
	// SlotInfo holds information on what item stack should be present in specific slots in the container.
	SlotInfo []StackResponseSlotInfo
}

// Marshal encodes/decodes a StackResponseContainerInfo.
func (x *StackResponseContainerInfo) Marshal(r protocol.IO) {
	protocol.Single(r, &x.Container)
	protocol.Slice(r, &x.SlotInfo)
}

// StackResponseSlotInfo holds information on what item stack should be present in a specific slot.
type StackResponseSlotInfo struct {
	// Slot and HotbarSlot seem to be the same value every time: The slot that was actually changed. I'm not
	// sure if these slots ever differ.
	Slot, HotbarSlot byte
	// Count is the total count of the item stack. This count will be shown client-side after the response is
	// sent to the client.
	Count byte
	// StackNetworkID is the network ID of the new stack at a specific slot.
	StackNetworkID int32
	// CustomName is the custom name of the item stack. It is used in relation to text filtering.
	CustomName string
	// FilteredCustomName is a filtered version of CustomName with all the profanity removed. The client will
	// use this over CustomName if this field is not empty and they have the "Filter Profanity" setting enabled.
	FilteredCustomName string
	// DurabilityCorrection is the current durability of the item stack. This durability will be shown
	// client-side after the response is sent to the client.
	DurabilityCorrection int32
}

// Marshal encodes/decodes a StackResponseSlotInfo.
func (x *StackResponseSlotInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.Slot)
	r.Uint8(&x.HotbarSlot)
	r.Uint8(&x.Count)
	r.Varint32(&x.StackNetworkID)
	if x.Slot != x.HotbarSlot {
		r.InvalidValue(x.HotbarSlot, "hotbar slot", "hot bar slot must be equal to normal slot")
	}
	r.String(&x.CustomName)
	if IsProtoGTE(r, ID766) {
		r.String(&x.FilteredCustomName)
	}
	r.Varint32(&x.DurabilityCorrection)
}

// StackRequestAction represents a single action related to the inventory present in an ItemStackRequest.
// The action is one of the concrete types below, each of which are indicative of a different action by the
// client, such as moving an item around the inventory or placing a block. It is an alias of Marshaler.
type StackRequestAction interface {
	protocol.Marshaler
}

const (
	StackRequestActionTake = iota
	StackRequestActionPlace
	StackRequestActionSwap
	StackRequestActionDrop
	StackRequestActionDestroy
	StackRequestActionConsume
	StackRequestActionCreate
	StackRequestActionPlaceInContainer
	StackRequestActionTakeOutContainer
	StackRequestActionLabTableCombine
	StackRequestActionBeaconPayment
	StackRequestActionMineBlock
	StackRequestActionCraftRecipe
	StackRequestActionCraftRecipeAuto
	StackRequestActionCraftCreative
	StackRequestActionCraftRecipeOptional
	StackRequestActionCraftGrindstone
	StackRequestActionCraftLoom
	StackRequestActionCraftNonImplementedDeprecated
	StackRequestActionCraftResultsDeprecated
)

// transferStackRequestAction is the structure shared by StackRequestActions that transfer items from one
// slot into another.
type transferStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *transferStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

// TakeStackRequestAction is sent by the client to the server to take x amount of items from one slot in a
// container to the cursor.
type TakeStackRequestAction struct {
	transferStackRequestAction
}

// PlaceStackRequestAction is sent by the client to the server to place x amount of items from one slot into
// another slot, such as when shift clicking an item in the inventory to move it around or when moving an item
// in the cursor into a slot.
type PlaceStackRequestAction struct {
	transferStackRequestAction
}

// SwapStackRequestAction is sent by the client to swap the item in its cursor with an item present in another
// container. The two item stacks swap places.
type SwapStackRequestAction struct {
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *SwapStackRequestAction) Marshal(r protocol.IO) {
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

// DropStackRequestAction is sent by the client when it drops an item out of the inventory when it has its
// inventory opened. This action is not sent when a player drops an item out of the hotbar using the Q button
// (or the equivalent on mobile). The InventoryTransaction packet is still used for that action, regardless of
// whether the item stack network IDs are used or not.
type DropStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source is the source slot from which items were dropped to the ground.
	Source StackRequestSlotInfo
	// Randomly seems to be set to false in most cases. I'm not entirely sure what this does, but this is what
	// vanilla calls this field.
	Randomly bool
}

// Marshal ...
func (a *DropStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	r.Bool(&a.Randomly)
}

// DestroyStackRequestAction is sent by the client when it destroys an item in creative mode by moving it
// back into the creative inventory.
type DestroyStackRequestAction struct {
	// Count is the count of the item in the source slot that was destroyed.
	Count byte
	// Source is the source slot from which items came that were destroyed by moving them into the creative
	// inventory.
	Source StackRequestSlotInfo
}

// Marshal ...
func (a *DestroyStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
}

// ConsumeStackRequestAction is sent by the client when it uses an item to craft another item. The original
// item is 'consumed'.
type ConsumeStackRequestAction struct {
	DestroyStackRequestAction
}

// CreateStackRequestAction is sent by the client when an item is created through being used as part of a
// recipe. For example, when milk is used to craft a cake, the buckets are leftover. The buckets are moved to
// the slot sent by the client here.
// Note that before this is sent, an action for consuming all items in the crafting table/grid is sent. Items
// that are not fully consumed when used for a recipe should not be destroyed there, but instead, should be
// turned into their respective resulting items.
type CreateStackRequestAction struct {
	// ResultsSlot is the slot in the inventory in which the results of the crafting ingredients are to be
	// placed.
	ResultsSlot byte
}

// Marshal ...
func (a *CreateStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.ResultsSlot)
}

// PlaceInContainerStackRequestAction currently has no known purpose.
type PlaceInContainerStackRequestAction struct {
	transferStackRequestAction
}

// TakeOutContainerStackRequestAction currently has no known purpose.
type TakeOutContainerStackRequestAction struct {
	transferStackRequestAction
}

// LabTableCombineStackRequestAction is sent by the client when it uses a lab table to combine item stacks.
type LabTableCombineStackRequestAction struct{}

// Marshal ...
func (a *LabTableCombineStackRequestAction) Marshal(protocol.IO) {}

// BeaconPaymentStackRequestAction is sent by the client when it submits an item to enable effects from a
// beacon. These items will have been moved into the beacon item slot in advance.
type BeaconPaymentStackRequestAction struct {
	// PrimaryEffect and SecondaryEffect are the effects that were selected from the beacon.
	PrimaryEffect, SecondaryEffect int32
}

// Marshal ...
func (a *BeaconPaymentStackRequestAction) Marshal(r protocol.IO) {
	r.Varint32(&a.PrimaryEffect)
	r.Varint32(&a.SecondaryEffect)
}

// MineBlockStackRequestAction is sent by the client when it breaks a block.
type MineBlockStackRequestAction struct {
	// HotbarSlot is the slot held by the player while mining a block.
	HotbarSlot int32
	// PredictedDurability is the durability of the item that the client assumes to be present at the time.
	PredictedDurability int32
	// StackNetworkID is the unique stack ID that the client assumes to be present at the time. The server
	// must check if these IDs match. If they do not match, servers should reject the stack request that the
	// action holding this info was in.
	StackNetworkID int32
}

// Marshal ...
func (a *MineBlockStackRequestAction) Marshal(r protocol.IO) {
	r.Varint32(&a.HotbarSlot)
	r.Varint32(&a.PredictedDurability)
	r.Varint32(&a.StackNetworkID)
}

// CraftRecipeStackRequestAction is sent by the client the moment it begins crafting an item. This is the
// first action sent, before the Consume and Create item stack request actions.
// This action is also sent when an item is enchanted. Enchanting should be treated mostly the same way as
// crafting, where the old item is consumed.
type CraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// NumberOfCrafts is how many times the recipe was crafted. This field appears to be boilerplate and
	// has no effect.
	NumberOfCrafts byte
}

func (a *CraftRecipeStackRequestAction) FromLatest(y *protocol.CraftRecipeStackRequestAction) *CraftRecipeStackRequestAction {
	a.RecipeNetworkID = y.RecipeNetworkID
	a.NumberOfCrafts = y.NumberOfCrafts
	return a
}

func (a *CraftRecipeStackRequestAction) ToLatest() *protocol.CraftRecipeStackRequestAction {
	return &protocol.CraftRecipeStackRequestAction{
		RecipeNetworkID: a.RecipeNetworkID,
		NumberOfCrafts:  a.NumberOfCrafts,
	}
}

// Marshal ...
func (a *CraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&a.NumberOfCrafts)
	}
}

// AutoCraftRecipeStackRequestAction is sent by the client similarly to the CraftRecipeStackRequestAction. The
// only difference is that the recipe is automatically created and crafted by shift clicking the recipe book.
type AutoCraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// NumberOfCrafts is how many times the recipe was crafted. This field is just a duplicate of TimesCrafted.
	NumberOfCrafts byte
	// TimesCrafted is how many times the recipe was crafted.
	TimesCrafted byte
	// Ingredients is a slice of ItemDescriptorCount that contains the ingredients that were used to craft the recipe.
	// It is not exactly clear what this is used for, but it is sent by the vanilla client.
	Ingredients []protocol.ItemDescriptorCount
}

// FromLatest ...
func (a *AutoCraftRecipeStackRequestAction) FromLatest(y *protocol.AutoCraftRecipeStackRequestAction) *AutoCraftRecipeStackRequestAction {
	a.RecipeNetworkID = y.RecipeNetworkID
	a.NumberOfCrafts = y.NumberOfCrafts
	a.TimesCrafted = y.TimesCrafted
	a.Ingredients = y.Ingredients
	return a
}

// ToLatest ...
func (a *AutoCraftRecipeStackRequestAction) ToLatest() *protocol.AutoCraftRecipeStackRequestAction {
	return &protocol.AutoCraftRecipeStackRequestAction{
		RecipeNetworkID: a.RecipeNetworkID,
		NumberOfCrafts:  a.NumberOfCrafts,
		TimesCrafted:    a.TimesCrafted,
		Ingredients:     a.Ingredients,
	}
}

// Marshal ...
func (a *AutoCraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
	r.Uint8(&a.NumberOfCrafts)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&a.TimesCrafted)
	}
	protocol.FuncSlice(r, &a.Ingredients, r.ItemDescriptorCount)
}

// CraftCreativeStackRequestAction is sent by the client when it takes an item out fo the creative inventory.
// The item is thus not really crafted, but instantly created.
type CraftCreativeStackRequestAction struct {
	// CreativeItemNetworkID is the network ID of the creative item that is being created. This is one of the
	// creative item network IDs sent in the CreativeContent packet.
	CreativeItemNetworkID uint32
	// NumberOfCrafts is how many times the recipe was crafted. This field appears to be boilerplate and
	// has no effect.
	NumberOfCrafts byte
}

// FromLatest ...
func (a *CraftCreativeStackRequestAction) FromLatest(y *protocol.CraftCreativeStackRequestAction) *CraftCreativeStackRequestAction {
	a.CreativeItemNetworkID = y.CreativeItemNetworkID
	a.NumberOfCrafts = y.NumberOfCrafts
	return a
}

// ToLatest ...
func (a *CraftCreativeStackRequestAction) ToLatest() *protocol.CraftCreativeStackRequestAction {
	return &protocol.CraftCreativeStackRequestAction{
		CreativeItemNetworkID: a.CreativeItemNetworkID,
		NumberOfCrafts:        a.NumberOfCrafts,
	}
}

// Marshal ...
func (a *CraftCreativeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.CreativeItemNetworkID)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&a.NumberOfCrafts)
	}
}

// CraftRecipeOptionalStackRequestAction is sent when using an anvil. When this action is sent, the
// FilterStrings field in the respective stack request is non-empty and contains the name of the item created
// using the anvil or cartography table.
type CraftRecipeOptionalStackRequestAction struct {
	// RecipeNetworkID is the network ID of the multi-recipe that is about to be crafted. This network ID matches
	// one of the multi-recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// FilterStringIndex is the index of a filter string sent in a ItemStackRequest.
	FilterStringIndex int32
}

// Marshal ...
func (c *CraftRecipeOptionalStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Int32(&c.FilterStringIndex)
}

// CraftGrindstoneRecipeStackRequestAction is sent when a grindstone recipe is crafted. It contains the RecipeNetworkID
// to identify the recipe crafted, and the cost for crafting the recipe.
type CraftGrindstoneRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// NumberOfCrafts is how many times the recipe was crafted. This field appears to be boilerplate and
	// has no effect.
	NumberOfCrafts byte
	// Cost is the cost of the recipe that was crafted.
	Cost int32
}

// FromLatest ...
func (c *CraftGrindstoneRecipeStackRequestAction) FromLatest(y *protocol.CraftGrindstoneRecipeStackRequestAction) *CraftGrindstoneRecipeStackRequestAction {
	c.RecipeNetworkID = y.RecipeNetworkID
	c.NumberOfCrafts = y.NumberOfCrafts
	c.Cost = y.Cost
	return c
}

// ToLatest ...
func (c *CraftGrindstoneRecipeStackRequestAction) ToLatest() *protocol.CraftGrindstoneRecipeStackRequestAction {
	return &protocol.CraftGrindstoneRecipeStackRequestAction{
		RecipeNetworkID: c.RecipeNetworkID,
		NumberOfCrafts:  c.NumberOfCrafts,
		Cost:            c.Cost,
	}
}

// Marshal ...
func (c *CraftGrindstoneRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&c.NumberOfCrafts)
	}
	r.Varint32(&c.Cost)
}

// CraftLoomRecipeStackRequestAction is sent when a loom recipe is crafted. It simply contains the
// pattern identifier to figure out what pattern is meant to be applied to the item.
type CraftLoomRecipeStackRequestAction struct {
	// Pattern is the pattern identifier for the loom recipe.
	Pattern string
	// TimesCrafted is how many times the recipe was crafted.
	TimesCrafted byte
}

// FromLatest ...
func (c *CraftLoomRecipeStackRequestAction) FromLatest(y *protocol.CraftLoomRecipeStackRequestAction) *CraftLoomRecipeStackRequestAction {
	c.Pattern = y.Pattern
	c.TimesCrafted = y.TimesCrafted
	return c
}

// ToLatest ...
func (c *CraftLoomRecipeStackRequestAction) ToLatest() *protocol.CraftLoomRecipeStackRequestAction {
	return &protocol.CraftLoomRecipeStackRequestAction{
		Pattern:      c.Pattern,
		TimesCrafted: c.TimesCrafted,
	}
}

// Marshal ...
func (c *CraftLoomRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.String(&c.Pattern)
	if IsProtoGTE(r, ID712) {
		r.Uint8(&c.TimesCrafted)
	}
}

// CraftNonImplementedStackRequestAction is an action sent for inventory actions that aren't yet implemented
// in the new system. These include, for example, anvils.
type CraftNonImplementedStackRequestAction struct{}

// Marshal ...
func (*CraftNonImplementedStackRequestAction) Marshal(protocol.IO) {}

// CraftResultsDeprecatedStackRequestAction is an additional, deprecated packet sent by the client after
// crafting. It holds the final results and the amount of times the recipe was crafted. It shouldn't be used.
// This action is also sent when an item is enchanted. Enchanting should be treated mostly the same way as
// crafting, where the old item is consumed.
type CraftResultsDeprecatedStackRequestAction struct {
	ResultItems  []protocol.ItemStack
	TimesCrafted byte
}

// Marshal ...
func (a *CraftResultsDeprecatedStackRequestAction) Marshal(r protocol.IO) {
	protocol.FuncSlice(r, &a.ResultItems, r.Item)
	r.Uint8(&a.TimesCrafted)
}

// StackRequestSlotInfo holds information on a specific slot client-side.
type StackRequestSlotInfo struct {
	// Container is the FullContainerName that describes the container that the slot is in.
	Container FullContainerName
	// Slot is the index of the slot within the container with the ContainerID above.
	Slot byte
	// StackNetworkID is the unique stack ID that the client assumes to be present in this slot. The server
	// must check if these IDs match. If they do not match, servers should reject the stack request that the
	// action holding this info was in.
	StackNetworkID int32
}

// StackReqSlotInfo reads/writes a StackRequestSlotInfo x using IO r.
func StackReqSlotInfo(r protocol.IO, x *StackRequestSlotInfo) {
	protocol.Single(r, &x.Container)
	r.Uint8(&x.Slot)
	r.Varint32(&x.StackNetworkID)
}
