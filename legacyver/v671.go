package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion671 ...
	ItemVersion671 = 181
	// BlockVersion671 ...
	BlockVersion671 int32 = (1 << 24) | (20 << 16) | (80 << 8)
)

var (
	//go:embed data/item_runtime_ids_671.nbt
	itemRuntimeIDData671 []byte
	//go:embed data/required_item_list_671.json
	requiredItemList671 []byte
	//go:embed data/block_states_671.nbt
	blockStateData671 []byte
)

// New671 ...
func New671() *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData671, requiredItemList671, ItemVersion671, false)
	blockMapping := mapping.NewBlockMapping(blockStateData671)

	return &Protocol{
		ver:             "1.20.80",
		id:              proto.ID671,
		blockTranslator: NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion671), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion671), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest, blockMapping, blockMappingLatest),
	}
}
