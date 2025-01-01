package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion712 ...
	ItemVersion712 = 201
	// BlockVersion712 ...
	BlockVersion712 int32 = (1 << 24) | (21 << 16) | (20 << 8)
)

var (
	//go:embed data/item_runtime_ids_712.nbt
	itemRuntimeIDData712 []byte
	//go:embed data/required_item_list_712.json
	requiredItemList712 []byte
	//go:embed data/block_states_712.nbt
	blockStateData712 []byte
)

func New712() *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData712, requiredItemList712, ItemVersion712, false)
	blockMapping := mapping.NewBlockMapping(blockStateData712)

	return &Protocol{
		ver:             "1.21.20",
		id:              proto.ID712,
		blockTranslator: NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion712), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion712), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest, blockMapping, blockMappingLatest),
	}
}
