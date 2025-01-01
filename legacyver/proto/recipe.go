package proto

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	RecipeUnlockContextNone = iota
	RecipeUnlockContextAlwaysUnlocked
	RecipeUnlockContextPlayerInWater
	RecipeUnlockContextPlayerHasManyItems
)

// RecipeUnlockRequirement represents a requirement that must be met in order to unlock a recipe. This is used
// for both shaped and shapeless recipes.
type RecipeUnlockRequirement struct {
	// Context is the context in which the recipe is unlocked. This is one of the constants above.
	Context byte
	// Ingredients are the ingredients required to unlock the recipe and only used if Context is set to none.
	Ingredients []protocol.ItemDescriptorCount
}

func (x *RecipeUnlockRequirement) FromLatest(requirement protocol.RecipeUnlockRequirement) RecipeUnlockRequirement {
	x.Context = requirement.Context
	x.Ingredients = requirement.Ingredients
	return *x
}

func (x *RecipeUnlockRequirement) ToLatest() protocol.RecipeUnlockRequirement {
	return protocol.RecipeUnlockRequirement{
		Context:     x.Context,
		Ingredients: x.Ingredients,
	}
}

// Marshal ...
func (x *RecipeUnlockRequirement) Marshal(r protocol.IO) {
	r.Uint8(&x.Context)
	if x.Context == RecipeUnlockContextNone {
		protocol.FuncSlice(r, &x.Ingredients, r.ItemDescriptorCount)
	}
}

const (
	RecipeShapeless int32 = iota
	RecipeShaped
	RecipeFurnace
	RecipeFurnaceData
	RecipeMulti
	RecipeShulkerBox
	RecipeShapelessChemistry
	RecipeShapedChemistry
	RecipeSmithingTransform
	RecipeSmithingTrim
)

// Recipe represents a recipe that may be sent in a CraftingData packet to let the client know what recipes
// are available server-side.
type Recipe interface {
	// Marshal encodes the recipe data to its binary representation into buf.
	Marshal(w *protocol.Writer)
	// Unmarshal decodes a serialised recipe from Reader r into the recipe instance.
	Unmarshal(r *protocol.Reader)
}

func takePtr[T any](x T) *T {
	return &x
}

func RecipeToLatest(x Recipe) protocol.Recipe {
	switch x := x.(type) {
	case *ShapelessRecipe:
		return takePtr(x.ToLatest())
	case *ShulkerBoxRecipe:
		return takePtr(x.ToLatest())
	case *ShapelessChemistryRecipe:
		return takePtr(x.ToLatest())
	case *ShapedRecipe:
		return takePtr(x.ToLatest())
	case *ShapedChemistryRecipe:
		return takePtr(x.ToLatest())
	case *FurnaceRecipe:
		return takePtr(x.ToLatest())
	case *FurnaceDataRecipe:
		return takePtr(x.ToLatest())
	case *MultiRecipe:
		return takePtr(x.ToLatest())
	case *SmithingTransformRecipe:
		return takePtr(x.ToLatest())
	case *SmithingTrimRecipe:
		return takePtr(x.ToLatest())
	}
	fmt.Println("RecipeToLatest: unknown recipe type")
	return x
}

func RecipeFromLatest(x protocol.Recipe) Recipe {
	switch x := x.(type) {
	case *protocol.ShapelessRecipe:
		return (&ShapelessRecipe{}).FromLatest(*x)
	case *protocol.ShulkerBoxRecipe:
		return (&ShulkerBoxRecipe{}).FromLatest(*x)
	case *protocol.ShapelessChemistryRecipe:
		return (&ShapelessChemistryRecipe{}).FromLatest(*x)
	case *protocol.ShapedRecipe:
		return (&ShapedRecipe{}).FromLatest(*x)
	case *protocol.ShapedChemistryRecipe:
		return (&ShapedChemistryRecipe{}).FromLatest(*x)
	case *protocol.FurnaceRecipe:
		return (&FurnaceRecipe{}).FromLatest(*x)
	case *protocol.FurnaceDataRecipe:
		return (&FurnaceDataRecipe{}).FromLatest(*x)
	case *protocol.MultiRecipe:
		return (&MultiRecipe{}).FromLatest(*x)
	case *protocol.SmithingTransformRecipe:
		return (&SmithingTransformRecipe{}).FromLatest(*x)
	case *protocol.SmithingTrimRecipe:
		return (&SmithingTrimRecipe{}).FromLatest(*x)
	}
	fmt.Println("RecipeFromLatest: unknown recipe type")
	return x
}

// lookupRecipe looks up the Recipe for a recipe type. False is returned if not
// found.
func lookupRecipe(recipeType int32, x *Recipe) bool {
	switch recipeType {
	case RecipeShapeless:
		*x = &ShapelessRecipe{}
	case RecipeShaped:
		*x = &ShapedRecipe{}
	case RecipeFurnace:
		*x = &FurnaceRecipe{}
	case RecipeFurnaceData:
		*x = &FurnaceDataRecipe{}
	case RecipeMulti:
		*x = &MultiRecipe{}
	case RecipeShulkerBox:
		*x = &ShulkerBoxRecipe{}
	case RecipeShapelessChemistry:
		*x = &ShapelessChemistryRecipe{}
	case RecipeShapedChemistry:
		*x = &ShapedChemistryRecipe{}
	case RecipeSmithingTransform:
		*x = &SmithingTransformRecipe{}
	case RecipeSmithingTrim:
		*x = &SmithingTrimRecipe{}
	default:
		return false
	}
	return true
}

// lookupRecipeType looks up the recipe type for a Recipe. False is returned if
// none was found.
func lookupRecipeType(x Recipe, recipeType *int32) bool {
	switch x.(type) {
	case *ShapelessRecipe:
		*recipeType = RecipeShapeless
	case *ShapedRecipe:
		*recipeType = RecipeShaped
	case *FurnaceRecipe:
		*recipeType = RecipeFurnace
	case *FurnaceDataRecipe:
		*recipeType = RecipeFurnaceData
	case *MultiRecipe:
		*recipeType = RecipeMulti
	case *ShulkerBoxRecipe:
		*recipeType = RecipeShulkerBox
	case *ShapelessChemistryRecipe:
		*recipeType = RecipeShapelessChemistry
	case *ShapedChemistryRecipe:
		*recipeType = RecipeShapedChemistry
	case *SmithingTransformRecipe:
		*recipeType = RecipeSmithingTransform
	case *SmithingTrimRecipe:
		*recipeType = RecipeSmithingTrim
	default:
		return false
	}
	return true
}

// ShapelessRecipe is a recipe that has no particular shape. Its functionality is shared with the
// RecipeShulkerBox and RecipeShapelessChemistry types.
type ShapelessRecipe struct {
	// RecipeID is a unique ID of the recipe. This ID must be unique amongst all other types of recipes too,
	// but its functionality is not exactly known.
	RecipeID string
	// Input is a list of items that serve as the input of the shapeless recipe. These items are the items
	// required to craft the output.
	Input []protocol.ItemDescriptorCount
	// Output is a list of items that are created as a result of crafting the recipe.
	Output []protocol.ItemStack
	// UUID is a UUID identifying the recipe. Since the CraftingEvent packet no longer exists, this can always be empty.
	UUID uuid.UUID
	// Block is the block name that is required to craft the output of the recipe. The block is not prefixed
	// with 'minecraft:', so it will look like 'crafting_table' as an example.
	// The available blocks are:
	// - crafting_table
	// - cartography_table
	// - stonecutter
	// - furnace
	// - blast_furnace
	// - smoker
	// - campfire
	Block string
	// Priority ...
	Priority int32
	// UnlockRequirement is a requirement that must be met in order to unlock the recipe.
	UnlockRequirement RecipeUnlockRequirement
	// RecipeNetworkID is a unique ID used to identify the recipe over network. Each recipe must have a unique
	// network ID. Recommended is to just increment a variable for each unique recipe registered.
	// This field must never be 0.
	RecipeNetworkID uint32
}

// ShulkerBoxRecipe is a shapeless recipe made specifically for shulker box crafting, so that they don't lose
// their user data when dyeing a shulker box.
type ShulkerBoxRecipe struct {
	ShapelessRecipe
}

// ShapelessChemistryRecipe is a recipe specifically made for chemistry related features, which exist only in
// the Education Edition. They function the same as shapeless recipes do.
type ShapelessChemistryRecipe struct {
	ShapelessRecipe
}

// ShapedRecipe is a recipe that has a specific shape that must be used to craft the output of the recipe.
// Trying to craft the item in any other shape will not work. The ShapedRecipe is of the same structure as the
// ShapedChemistryRecipe.
type ShapedRecipe struct {
	// RecipeID is a unique ID of the recipe. This ID must be unique amongst all other types of recipes too,
	// but its functionality is not exactly known.
	RecipeID string
	// Width is the width of the recipe's shape.
	Width int32
	// Height is the height of the recipe's shape.
	Height int32
	// Input is a list of items that serve as the input of the shapeless recipe. These items are the items
	// required to craft the output. The amount of input items must be exactly equal to Width * Height.
	Input []protocol.ItemDescriptorCount
	// Output is a list of items that are created as a result of crafting the recipe.
	Output []protocol.ItemStack
	// UUID is a UUID identifying the recipe. Since the CraftingEvent packet no longer exists, this can always be empty.
	UUID uuid.UUID
	// Block is the block name that is required to craft the output of the recipe. The block is not prefixed
	// with 'minecraft:', so it will look like 'crafting_table' as an example.
	Block string
	// Priority ...
	Priority int32
	// AssumeSymmetry specifies if the recipe is symmetrical. If this is set to true, the recipe will be
	// mirrored along the diagonal axis. This means that the recipe will be the same if rotated 180 degrees.
	AssumeSymmetry bool
	// UnlockRequirement is a requirement that must be met in order to unlock the recipe.
	UnlockRequirement RecipeUnlockRequirement
	// RecipeNetworkID is a unique ID used to identify the recipe over network. Each recipe must have a unique
	// network ID. Recommended is to just increment a variable for each unique recipe registered.
	// This field must never be 0.
	RecipeNetworkID uint32
}

// ShapedChemistryRecipe is a recipe specifically made for chemistry related features, which exist only in the
// Education Edition. It functions the same as a normal ShapedRecipe.
type ShapedChemistryRecipe struct {
	ShapedRecipe
}

// FurnaceRecipe is a recipe that is specifically used for all kinds of furnaces. These recipes don't just
// apply to furnaces, but also blast furnaces and smokers.
type FurnaceRecipe struct {
	// InputType is the item type of the input item. The metadata value of the item is not used in the
	// FurnaceRecipe. Use FurnaceDataRecipe to allow an item with only one metadata value.
	InputType protocol.ItemType
	// Output is the item that is created as a result of smelting/cooking an item in the furnace.
	Output protocol.ItemStack
	// Block is the block name that is required to create the output of the recipe. The block is not prefixed
	// with 'minecraft:', so it will look like 'furnace' as an example.
	Block string
}

// FurnaceDataRecipe is a recipe specifically used for furnace-type crafting stations. It is equal to
// FurnaceRecipe, except it has an input item with a specific metadata value, instead of any metadata value.
type FurnaceDataRecipe struct {
	FurnaceRecipe
}

// MultiRecipe serves as an 'enable' switch for multi-shape recipes.
type MultiRecipe struct {
	// UUID is a UUID identifying the recipe. Since the CraftingEvent packet no longer exists, this can always be empty.
	UUID uuid.UUID
	// RecipeNetworkID is a unique ID used to identify the recipe over network. Each recipe must have a unique
	// network ID. Recommended is to just increment a variable for each unique recipe registered.
	// This field must never be 0.
	RecipeNetworkID uint32
}

// SmithingTransformRecipe is a recipe specifically used for smithing tables. It has three input items and adds them
// together, resulting in a new item.
type SmithingTransformRecipe struct {
	// RecipeNetworkID is a unique ID used to identify the recipe over network. Each recipe must have a unique
	// network ID. Recommended is to just increment a variable for each unique recipe registered.
	// This field must never be 0.
	RecipeNetworkID uint32
	// RecipeID is a unique ID of the recipe. This ID must be unique amongst all other types of recipes too,
	// but its functionality is not exactly known.
	RecipeID string
	// Template is the item that is used to shape the Base item based on the Addition being applied.
	Template protocol.ItemDescriptorCount
	// Base is the item that the Addition is being applied to in the smithing table.
	Base protocol.ItemDescriptorCount
	// Addition is the item that is being added to the Base item to result in a modified item.
	Addition protocol.ItemDescriptorCount
	// Result is the resulting item from the two items being added together.
	Result protocol.ItemStack
	// Block is the block name that is required to create the output of the recipe. The block is not prefixed with
	// 'minecraft:', so it will look like 'smithing_table' as an example.
	Block string
}

// SmithingTrimRecipe is a recipe specifically used for applying armour trims to an armour piece inside a smithing table.
type SmithingTrimRecipe struct {
	// RecipeNetworkID is a unique ID used to identify the recipe over network. Each recipe must have a unique
	// network ID. Recommended is to just increment a variable for each unique recipe registered.
	// This field must never be 0.
	RecipeNetworkID uint32
	// RecipeID is a unique ID of the recipe. This ID must be unique amongst all other types of recipes too,
	// but its functionality is not exactly known.
	RecipeID string
	// Template is the item that is used to shape the Base item based on the Addition being applied.
	Template protocol.ItemDescriptorCount
	// Base is the item that the Addition is being applied to in the smithing table.
	Base protocol.ItemDescriptorCount
	// Addition is the item that is being added to the Base item to result in a modified item.
	Addition protocol.ItemDescriptorCount
	// Block is the block name that is required to create the output of the recipe. The block is not prefixed with
	// 'minecraft:', so it will look like 'smithing_table' as an example.
	Block string
}

// FromLatest ...
func (recipe *ShapelessRecipe) FromLatest(v protocol.ShapelessRecipe) *ShapelessRecipe {
	recipe.RecipeID = v.RecipeID
	recipe.Input = v.Input
	recipe.Output = v.Output
	recipe.UUID = v.UUID
	recipe.Block = v.Block
	recipe.Priority = v.Priority
	recipe.UnlockRequirement = (&RecipeUnlockRequirement{}).FromLatest(v.UnlockRequirement)
	recipe.RecipeNetworkID = v.RecipeNetworkID
	return recipe
}

// ToLatest ...
func (recipe *ShapelessRecipe) ToLatest() protocol.ShapelessRecipe {
	return protocol.ShapelessRecipe{
		RecipeID:          recipe.RecipeID,
		Input:             recipe.Input,
		Output:            recipe.Output,
		UUID:              recipe.UUID,
		Block:             recipe.Block,
		Priority:          recipe.Priority,
		UnlockRequirement: recipe.UnlockRequirement.ToLatest(),
		RecipeNetworkID:   recipe.RecipeNetworkID,
	}
}

// Marshal ...
func (recipe *ShapelessRecipe) Marshal(w *protocol.Writer) {
	marshalShapeless(w, recipe)
}

// Unmarshal ...
func (recipe *ShapelessRecipe) Unmarshal(r *protocol.Reader) {
	marshalShapeless(r, recipe)
}

// FromLatest ...
func (recipe *ShulkerBoxRecipe) FromLatest(v protocol.ShulkerBoxRecipe) *ShulkerBoxRecipe {
	recipe.ShapelessRecipe = *(&ShapelessRecipe{}).FromLatest(v.ShapelessRecipe)
	return recipe
}

// ToLatest ...
func (recipe *ShulkerBoxRecipe) ToLatest() protocol.ShulkerBoxRecipe {
	return protocol.ShulkerBoxRecipe{
		ShapelessRecipe: recipe.ShapelessRecipe.ToLatest(),
	}
}

// Marshal ...
func (recipe *ShulkerBoxRecipe) Marshal(w *protocol.Writer) {
	marshalShapeless(w, &recipe.ShapelessRecipe)
}

// Unmarshal ...
func (recipe *ShulkerBoxRecipe) Unmarshal(r *protocol.Reader) {
	marshalShapeless(r, &recipe.ShapelessRecipe)
}

// FromLatest ...
func (recipe *ShapelessChemistryRecipe) FromLatest(v protocol.ShapelessChemistryRecipe) *ShapelessChemistryRecipe {
	recipe.ShapelessRecipe = *(&ShapelessRecipe{}).FromLatest(v.ShapelessRecipe)
	return recipe
}

// ToLatest ...
func (recipe *ShapelessChemistryRecipe) ToLatest() protocol.ShapelessChemistryRecipe {
	return protocol.ShapelessChemistryRecipe{
		ShapelessRecipe: recipe.ShapelessRecipe.ToLatest(),
	}
}

// Marshal ...
func (recipe *ShapelessChemistryRecipe) Marshal(w *protocol.Writer) {
	marshalShapeless(w, &recipe.ShapelessRecipe)
}

// Unmarshal ...
func (recipe *ShapelessChemistryRecipe) Unmarshal(r *protocol.Reader) {
	marshalShapeless(r, &recipe.ShapelessRecipe)
}

// FromLatest ...
func (recipe *ShapedRecipe) FromLatest(v protocol.ShapedRecipe) *ShapedRecipe {
	recipe.RecipeID = v.RecipeID
	recipe.Width = v.Width
	recipe.Height = v.Height
	recipe.Input = v.Input
	recipe.Output = v.Output
	recipe.UUID = v.UUID
	recipe.Block = v.Block
	recipe.AssumeSymmetry = v.AssumeSymmetry
	recipe.UnlockRequirement = (&RecipeUnlockRequirement{}).FromLatest(v.UnlockRequirement)
	recipe.Priority = v.Priority
	recipe.RecipeNetworkID = v.RecipeNetworkID
	return recipe
}

// ToLatest ...
func (recipe *ShapedRecipe) ToLatest() protocol.ShapedRecipe {
	return protocol.ShapedRecipe{
		RecipeID:          recipe.RecipeID,
		Width:             recipe.Width,
		Height:            recipe.Height,
		Input:             recipe.Input,
		Output:            recipe.Output,
		UUID:              recipe.UUID,
		Block:             recipe.Block,
		AssumeSymmetry:    recipe.AssumeSymmetry,
		UnlockRequirement: recipe.UnlockRequirement.ToLatest(),
		Priority:          recipe.Priority,
		RecipeNetworkID:   recipe.RecipeNetworkID,
	}
}

// Marshal ...
func (recipe *ShapedRecipe) Marshal(w *protocol.Writer) {
	marshalShaped(w, recipe)
}

// Unmarshal ...
func (recipe *ShapedRecipe) Unmarshal(r *protocol.Reader) {
	marshalShaped(r, recipe)
}

// FromLatest ...
func (recipe *ShapedChemistryRecipe) FromLatest(v protocol.ShapedChemistryRecipe) *ShapedChemistryRecipe {
	recipe.ShapedRecipe = *(&ShapedRecipe{}).FromLatest(v.ShapedRecipe)
	return recipe
}

// ToLatest ...
func (recipe *ShapedChemistryRecipe) ToLatest() protocol.ShapedChemistryRecipe {
	return protocol.ShapedChemistryRecipe{
		ShapedRecipe: recipe.ShapedRecipe.ToLatest(),
	}
}

// Marshal ...
func (recipe *ShapedChemistryRecipe) Marshal(w *protocol.Writer) {
	marshalShaped(w, &recipe.ShapedRecipe)
}

// Unmarshal ...
func (recipe *ShapedChemistryRecipe) Unmarshal(r *protocol.Reader) {
	marshalShaped(r, &recipe.ShapedRecipe)
}

// FromLatest ...
func (recipe *FurnaceRecipe) FromLatest(v protocol.FurnaceRecipe) *FurnaceRecipe {
	recipe.InputType = v.InputType
	recipe.Output = v.Output
	recipe.Block = v.Block
	return recipe
}

// ToLatest ...
func (recipe *FurnaceRecipe) ToLatest() protocol.FurnaceRecipe {
	return protocol.FurnaceRecipe{
		InputType: recipe.InputType,
		Output:    recipe.Output,
		Block:     recipe.Block,
	}
}

// Marshal ...
func (recipe *FurnaceRecipe) Marshal(w *protocol.Writer) {
	w.Varint32(&recipe.InputType.NetworkID)
	w.Item(&recipe.Output)
	w.String(&recipe.Block)
}

// Unmarshal ...
func (recipe *FurnaceRecipe) Unmarshal(r *protocol.Reader) {
	r.Varint32(&recipe.InputType.NetworkID)
	r.Item(&recipe.Output)
	r.String(&recipe.Block)
}

// FromLatest ...
func (recipe *FurnaceDataRecipe) FromLatest(v protocol.FurnaceDataRecipe) *FurnaceDataRecipe {
	recipe.FurnaceRecipe = *(&FurnaceRecipe{}).FromLatest(v.FurnaceRecipe)
	return recipe
}

// ToLatest ...
func (recipe *FurnaceDataRecipe) ToLatest() protocol.FurnaceDataRecipe {
	return protocol.FurnaceDataRecipe{
		FurnaceRecipe: recipe.FurnaceRecipe.ToLatest(),
	}
}

// Marshal ...
func (recipe *FurnaceDataRecipe) Marshal(w *protocol.Writer) {
	w.Varint32(&recipe.InputType.NetworkID)
	aux := int32(recipe.InputType.MetadataValue)
	w.Varint32(&aux)
	w.Item(&recipe.Output)
	w.String(&recipe.Block)
}

// Unmarshal ...
func (recipe *FurnaceDataRecipe) Unmarshal(r *protocol.Reader) {
	var dataValue int32
	r.Varint32(&recipe.InputType.NetworkID)
	r.Varint32(&dataValue)
	recipe.InputType.MetadataValue = uint32(dataValue)
	r.Item(&recipe.Output)
	r.String(&recipe.Block)
}

// FromLatest ...
func (recipe *MultiRecipe) FromLatest(v protocol.MultiRecipe) *MultiRecipe {
	recipe.UUID = v.UUID
	recipe.RecipeNetworkID = v.RecipeNetworkID
	return recipe
}

// ToLatest ...
func (recipe *MultiRecipe) ToLatest() protocol.MultiRecipe {
	return protocol.MultiRecipe{
		UUID:            recipe.UUID,
		RecipeNetworkID: recipe.RecipeNetworkID,
	}
}

// Marshal ...
func (recipe *MultiRecipe) Marshal(w *protocol.Writer) {
	w.UUID(&recipe.UUID)
	w.Varuint32(&recipe.RecipeNetworkID)
}

// Unmarshal ...
func (recipe *MultiRecipe) Unmarshal(r *protocol.Reader) {
	r.UUID(&recipe.UUID)
	r.Varuint32(&recipe.RecipeNetworkID)
}

// FromLatest ...
func (recipe *SmithingTransformRecipe) FromLatest(v protocol.SmithingTransformRecipe) *SmithingTransformRecipe {
	recipe.RecipeNetworkID = v.RecipeNetworkID
	recipe.RecipeID = v.RecipeID
	recipe.Template = v.Template
	recipe.Base = v.Base
	recipe.Addition = v.Addition
	recipe.Result = v.Result
	recipe.Block = v.Block
	return recipe
}

// ToLatest ...
func (recipe *SmithingTransformRecipe) ToLatest() protocol.SmithingTransformRecipe {
	return protocol.SmithingTransformRecipe{
		RecipeNetworkID: recipe.RecipeNetworkID,
		RecipeID:        recipe.RecipeID,
		Template:        recipe.Template,
		Base:            recipe.Base,
		Addition:        recipe.Addition,
		Result:          recipe.Result,
		Block:           recipe.Block,
	}
}

// Marshal ...
func (recipe *SmithingTransformRecipe) Marshal(w *protocol.Writer) {
	w.String(&recipe.RecipeID)
	w.ItemDescriptorCount(&recipe.Template)
	w.ItemDescriptorCount(&recipe.Base)
	w.ItemDescriptorCount(&recipe.Addition)
	w.Item(&recipe.Result)
	w.String(&recipe.Block)
	w.Varuint32(&recipe.RecipeNetworkID)
}

// Unmarshal ...
func (recipe *SmithingTransformRecipe) Unmarshal(r *protocol.Reader) {
	r.String(&recipe.RecipeID)
	r.ItemDescriptorCount(&recipe.Template)
	r.ItemDescriptorCount(&recipe.Base)
	r.ItemDescriptorCount(&recipe.Addition)
	r.Item(&recipe.Result)
	r.String(&recipe.Block)
	r.Varuint32(&recipe.RecipeNetworkID)
}

// FromLatest ...
func (recipe *SmithingTrimRecipe) FromLatest(v protocol.SmithingTrimRecipe) *SmithingTrimRecipe {
	recipe.RecipeNetworkID = v.RecipeNetworkID
	recipe.RecipeID = v.RecipeID
	recipe.Template = v.Template
	recipe.Base = v.Base
	recipe.Addition = v.Addition
	recipe.Block = v.Block
	return recipe
}

// ToLatest ...
func (recipe *SmithingTrimRecipe) ToLatest() protocol.SmithingTrimRecipe {
	return protocol.SmithingTrimRecipe{
		RecipeNetworkID: recipe.RecipeNetworkID,
		RecipeID:        recipe.RecipeID,
		Template:        recipe.Template,
		Base:            recipe.Base,
		Addition:        recipe.Addition,
		Block:           recipe.Block,
	}
}

// Marshal ...
func (recipe *SmithingTrimRecipe) Marshal(w *protocol.Writer) {
	w.String(&recipe.RecipeID)
	w.ItemDescriptorCount(&recipe.Template)
	w.ItemDescriptorCount(&recipe.Base)
	w.ItemDescriptorCount(&recipe.Addition)
	w.String(&recipe.Block)
	w.Varuint32(&recipe.RecipeNetworkID)
}

// Unmarshal ...
func (recipe *SmithingTrimRecipe) Unmarshal(r *protocol.Reader) {
	r.String(&recipe.RecipeID)
	r.ItemDescriptorCount(&recipe.Template)
	r.ItemDescriptorCount(&recipe.Base)
	r.ItemDescriptorCount(&recipe.Addition)
	r.String(&recipe.Block)
	r.Varuint32(&recipe.RecipeNetworkID)
}

// marshalShaped ...
func marshalShaped(r protocol.IO, recipe *ShapedRecipe) {
	r.String(&recipe.RecipeID)
	r.Varint32(&recipe.Width)
	r.Varint32(&recipe.Height)
	protocol.FuncSliceOfLen(r, uint32(recipe.Width*recipe.Height), &recipe.Input, r.ItemDescriptorCount)
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	r.Bool(&recipe.AssumeSymmetry)
	if IsProtoGTE(r, ID685) {
		protocol.Single(r, &recipe.UnlockRequirement)
	}
	r.Varuint32(&recipe.RecipeNetworkID)
}

// marshalShapeless ...
func marshalShapeless(r protocol.IO, recipe *ShapelessRecipe) {
	r.String(&recipe.RecipeID)
	protocol.FuncSlice(r, &recipe.Input, r.ItemDescriptorCount)
	protocol.FuncSlice(r, &recipe.Output, r.Item)
	r.UUID(&recipe.UUID)
	r.String(&recipe.Block)
	r.Varint32(&recipe.Priority)
	if IsProtoGTE(r, ID685) {
		protocol.Single(r, &recipe.UnlockRequirement)
	}
	r.Varuint32(&recipe.RecipeNetworkID)
}
