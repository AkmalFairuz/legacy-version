package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion685 ...
	ItemVersion685 = 191
	// BlockVersion685 ...
	BlockVersion685 int32 = (1 << 24) | (21 << 16) | (0 << 8)
)

// New685 uses same data as 686
func New685(direct bool) *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData686, requiredItemList686, ItemVersion685, false)
	blockMapping := mapping.NewBlockMapping(blockStateData686)

	return &Protocol{
		ver:             "1.21.0",
		id:              proto.ID685,
		blockTranslator: NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion685), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion685), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest, blockMapping, blockMappingLatest),
	}
}
