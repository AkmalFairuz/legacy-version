package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion729 ...
	ItemVersion729 = 211
	// BlockVersion729 ...
	BlockVersion729 int32 = (1 << 24) | (21 << 16) | (30 << 8)
)

var (
	//go:embed data/item_runtime_ids_729.nbt
	itemRuntimeIDData729 []byte
	//go:embed data/required_item_list_729.json
	requiredItemList729 []byte
	//go:embed data/block_states_729.nbt
	blockStateData729 []byte
)

func New729(direct bool) *Protocol {
	itemMapping := mapping.NewItemMapping(itemRuntimeIDData729, requiredItemList729, ItemVersion729, false)
	blockMapping := mapping.NewBlockMapping(blockStateData729)

	return &Protocol{
		ver:             "1.21.30",
		id:              proto.ID729,
		blockTranslator: NewBlockTranslator(blockMapping, latestBlockMapping, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion729), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion729), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest, blockMapping, blockMappingLatest),
	}
}
