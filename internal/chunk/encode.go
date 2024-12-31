package chunk

import (
	"bytes"
	"github.com/df-mc/dragonfly/server/block/cube"
	"sync"
)

// pool is used to pool byte buffers used for encoding chunks.
var pool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func NetworkEncode(air uint32, c *Chunk, oldFormat bool, pe PaletteEncoding) ([]byte, error) {
	buf := pool.Get().(*bytes.Buffer)
	for i := 0; i < len(c.sub); i++ {
		_, _ = buf.Write(EncodeSubChunk(c.sub[i], NetworkEncoding, pe, SubChunkVersion8, c.r, i))
	}
	if oldFormat {
		biomes := make([]byte, 256)

		// Make our 3D biomes 2D.
		for x := uint8(0); x < 16; x++ {
			for z := uint8(0); z < 16; z++ {
				biomes[(x&15)|(z&15)<<4] = byte(c.Biome(x, c.HighestBlock(x, z), z))
			}
		}
		_, _ = buf.Write(biomes)
	} else {
		_, _ = buf.Write(EncodeBiomes(c, NetworkEncoding))
	}

	return buf.Bytes(), nil
}

// EncodeSubChunk encodes a sub-chunk from a chunk into bytes. An Encoding may be passed to encode either for network or
// disk purposed, the most notable difference being that the network encoding generally uses varints and no NBT.
func EncodeSubChunk(s *SubChunk, e Encoding, pe PaletteEncoding, subChunkVer subChunkVersion, r cube.Range, ind int) []byte {
	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	subChunkVer.EncodeHeader(buf, s, r, ind)
	for _, storage := range s.storages {
		encodePalettedStorage(buf, storage, nil, e, pe)
	}
	sub := make([]byte, buf.Len())
	_, _ = buf.Read(sub)
	return sub
}

// EncodeBiomes encodes the biomes of a chunk into bytes. An Encoding may be passed to encode either for network or
// disk purposed, the most notable difference being that the network encoding generally uses varints and no NBT.
func EncodeBiomes(c *Chunk, e Encoding) []byte {
	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	var previous *PalettedStorage
	for _, b := range c.biomes {
		encodePalettedStorage(buf, b, previous, e, BiomePaletteEncoding)
		previous = b
	}
	biomes := make([]byte, buf.Len())
	_, _ = buf.Read(biomes)
	return biomes
}

// encodePalettedStorage encodes a PalettedStorage into a bytes.Buffer. The Encoding passed is used to write the Palette
// of the PalettedStorage.
func encodePalettedStorage(buf *bytes.Buffer, storage, previous *PalettedStorage, e Encoding, pe PaletteEncoding) {
	if storage.Equal(previous) {
		_, _ = buf.Write([]byte{0x7f<<1 | e.Network()})
		return
	}
	b := make([]byte, len(storage.indices)*4+1)
	b[0] = byte(storage.bitsPerIndex<<1) | e.Network()

	for i, v := range storage.indices {
		// Explicitly don't use the binary package to greatly improve performance of writing the uint32s.
		b[i*4+1], b[i*4+2], b[i*4+3], b[i*4+4] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	}
	_, _ = buf.Write(b)

	e.EncodePalette(buf, storage.palette, pe)
}
