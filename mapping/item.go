package mapping

import (
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"sync"
)

type Item interface {
	// ItemRuntimeIDToName converts an item runtime ID to a string ID.
	ItemRuntimeIDToName(int32) (string, bool)
	// ItemNameToRuntimeID converts a string ID to an item runtime ID.
	ItemNameToRuntimeID(string) (int32, bool)
	RegisterEntry(string) int32
	Air() int32
	ItemVersion() uint16
}

type DefaultItemMapping struct {
	mu sync.Mutex
	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames map[int32]string
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs map[string]int32
	airRID                int32
	itemVersion           uint16
}

func NewItemMapping(itemRuntimeIDData []byte, requiredItemList []byte, itemVersion uint16, direct bool) *DefaultItemMapping {
	itemRuntimeIDsToNames := make(map[int32]string)
	itemNamesToRuntimeIDs := make(map[string]int32)
	var airRID *int32

	if direct {
		var items map[string]int32
		if err := nbt.Unmarshal(itemRuntimeIDData, &items); err != nil {
			panic(err)
		}
		for name, rid := range items {
			if name == "minecraft:air" {
				airRID = &rid
			}

			itemNamesToRuntimeIDs[name] = rid
			itemRuntimeIDsToNames[rid] = name
		}
	} else {
		var m map[string]struct {
			RuntimeID      int16 `json:"runtime_id"`
			ComponentBased bool  `json:"component_based"`
		}
		if err := json.Unmarshal(requiredItemList, &m); err != nil {
			panic(err)
		}
		for name, data := range m {
			rid := int32(data.RuntimeID)
			if name == "minecraft:air" {
				airRID = &rid
			}

			itemNamesToRuntimeIDs[name] = rid
			itemRuntimeIDsToNames[rid] = name
		}
	}

	if airRID == nil {
		panic("couldn't find air")
	}

	return &DefaultItemMapping{itemRuntimeIDsToNames: itemRuntimeIDsToNames, itemNamesToRuntimeIDs: itemNamesToRuntimeIDs, itemVersion: itemVersion}
}

func (m *DefaultItemMapping) ItemRuntimeIDToName(runtimeID int32) (name string, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	name, ok := m.itemRuntimeIDsToNames[runtimeID]
	return name, ok
}

func (m *DefaultItemMapping) ItemNameToRuntimeID(name string) (runtimeID int32, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	rid, ok := m.itemNamesToRuntimeIDs[name]
	return rid, ok
}

func (m *DefaultItemMapping) RegisterEntry(name string) int32 {
	defer m.mu.Unlock()
	m.mu.Lock()
	if rid, ok := m.itemNamesToRuntimeIDs[name]; ok {
		return rid
	}
	nextRID := int32(len(m.itemRuntimeIDsToNames))
	m.itemNamesToRuntimeIDs[name] = nextRID
	m.itemRuntimeIDsToNames[nextRID] = name
	return nextRID
}

func (m *DefaultItemMapping) Air() int32 {
	defer m.mu.Unlock()
	m.mu.Lock()
	return m.airRID
}

func (m *DefaultItemMapping) ItemVersion() uint16 {
	defer m.mu.Unlock()
	m.mu.Lock()
	return m.itemVersion
}
