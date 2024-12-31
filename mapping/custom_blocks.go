package mapping

import (
	"github.com/akmalfairuz/legacy-version/internal"
	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"golang.org/x/exp/maps"
)

func convert(entries []protocol.BlockEntry) (states []blockupgrader.BlockState) {
	for _, entry := range entries {
		propertiesMap := map[string][]any{}
		if props := jsonCheck[[]any](entry.Properties, "properties"); props != nil {
			for _, prop := range *props {
				prop := prop.(map[string]any)
				name := jsonCheck[string](prop, "name")
				enum := jsonCheck[[]any](prop, "enum")
				if enum == nil {
					int32Enum := jsonCheck[[]int32](prop, "enum")
					if int32Enum != nil {
						enum = &[]any{}
						for _, i := range *int32Enum {
							*enum = append(*enum, i)
						}
					}
				}
				if name == nil || enum == nil {
					panic("could not find field `name` and `enum`")
				}
				propertiesMap[*name] = *enum
			}
		}

		combinations := make([]map[string]any, 0)
		generateCombinationsRecursively(propertiesMap, internal.NewIterator(maps.Keys(propertiesMap)), map[string]any{}, &combinations)
		if len(combinations) == 0 {
			combinations = append(combinations, map[string]any{})
		}

		for _, combination := range combinations {
			blockState := blockupgrader.BlockState{
				Name:       entry.Name,
				Properties: combination,
			}

			states = append(states, blockState)
		}
	}
	return
}

func generateCombinationsRecursively[K comparable, V any](all map[K][]V, iterator *internal.Iterator[K], current map[K]V, output *[]map[K]V) {
	if !iterator.HasNext() {
		entry := map[K]V{}
		for k, v := range current {
			entry[k] = v
		}
		out := append(*output, entry)
		*output = out
	} else {
		key := iterator.Next()
		set := all[key]

		for _, value := range set {
			current[key] = value
			generateCombinationsRecursively(all, iterator, current, output)
			delete(current, key)
		}

		iterator.Previous()
	}
}

func jsonCheck[T any](json map[string]any, field string) *T {
	fieldValue, found := json[field]
	if !found {
		return nil
	}
	castedValue, ok := fieldValue.(T)
	if !ok {
		return nil
	}
	return &castedValue
}
