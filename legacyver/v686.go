package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion686 ...
	ItemVersion686 = 191
	// BlockVersion686 ...
	BlockVersion686 int32 = (1 << 24) | (21 << 16) | (2 << 8)
)

var (
	//go:embed data/item_runtime_ids_686.nbt
	itemRuntimeIDData686 []byte
	//go:embed data/required_item_list_686.json
	requiredItemList686 []byte
	//go:embed data/block_states_686.nbt
	blockStateData686 []byte
)

func New686(direct bool) *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData686, requiredItemList686, ItemVersion686, false)
	blockMapping := mapping.NewBlockMapping(blockStateData686)

	return &Protocol{
		ver:             "1.21.2",
		id:              proto.ID686,
		blockTranslator: NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion686), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion686), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest, blockMapping, blockMappingLatest),
	}
}
