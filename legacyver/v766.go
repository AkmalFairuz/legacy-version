package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion766 ...
	ItemVersion766 = 231
	// BlockVersion766 ...
	BlockVersion766 int32 = (1 << 24) | (21 << 16) | (50 << 8)
)

var (
	//go:embed item_runtime_ids_766.nbt
	itemRuntimeIDData766 []byte
	//go:embed required_item_list_766.json
	requiredItemList766 []byte
	//go:embed block_states_766.nbt
	blockStateData766 []byte

	itemMappingLatest  = mapping.NewItemMapping(itemRuntimeIDData766, requiredItemList766, ItemVersion766, false)
	blockMappingLatest = mapping.NewBlockMapping(blockStateData766)
)
