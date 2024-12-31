package legacyver

import "github.com/sandertv/gophertunnel/minecraft/protocol"

func fitBitset(b protocol.Bitset, size int) protocol.Bitset {
	if b.Len() != size {
		ret := protocol.NewBitset(size)

		copySize := b.Len()
		if copySize > size {
			copySize = size
		}

		for i := 0; i < copySize; i++ {
			if b.Load(i) {
				ret.Set(i)
			}
		}
		return ret
	}
	return b
}
