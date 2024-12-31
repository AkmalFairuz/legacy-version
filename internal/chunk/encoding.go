package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/akmalfairuz/legacy-version/mapping"
	"strings"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	// SubChunkVersion is the current version of the written sub chunks, specifying the format they are
	// written on disk and over network.
	SubChunkVersion = 9
)

type (
	// Encoding is an encoding type used for Chunk encoding. Implementations of this interface are DiskEncoding and
	// NetworkEncoding, which can be used to encode a Chunk to an intermediate disk or network representation respectively.
	Encoding interface {
		EncodePalette(buf *bytes.Buffer, p *Palette, e PaletteEncoding)
		DecodePalette(buf *bytes.Buffer, blockSize paletteSize, e PaletteEncoding) (*Palette, error)
		Network() byte
	}
	// PaletteEncoding is an encoding type used for Chunk encoding. It is used to encode different types of palettes
	// (for example, blocks or biomes) differently.
	PaletteEncoding interface {
		Encode(buf *bytes.Buffer, v uint32)
		Decode(buf *bytes.Buffer) (uint32, error)
	}
	// Encoding is an encoding type used for Chunk encoding. Implementations of this interface are DiskEncoding and
	// NetworkEncoding, which can be used to encode a Chunk to an intermediate disk or network representation respectively.
	subChunkVersion interface {
		EncodeHeader(buf *bytes.Buffer, s *SubChunk, r cube.Range, ind int)
	}
)

var (
	// SubChunkVersion9 subChunkVersion9
	SubChunkVersion9 subChunkVersion9
	// SubChunkVersion8 subChunkVersion8
	SubChunkVersion8 subChunkVersion8

	// NetworkEncoding is the Encoding used for sending a Chunk over network. It does not use NBT and writes varints.
	NetworkEncoding networkEncoding
	// BiomePaletteEncoding is the paletteEncoding used for encoding a palette of biomes.
	BiomePaletteEncoding biomePaletteEncoding
)

// networkEncoding implements the Chunk encoding for sending over network.
type networkEncoding struct{}

func (networkEncoding) Network() byte { return 1 }
func (networkEncoding) EncodePalette(buf *bytes.Buffer, p *Palette, _ PaletteEncoding) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}
	for _, val := range p.values {
		_ = protocol.WriteVarint32(buf, int32(val))
	}
}
func (networkEncoding) DecodePalette(buf *bytes.Buffer, blockSize paletteSize, _ PaletteEncoding) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize != 0 {
		if err := protocol.Varint32(buf, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	blocks, temp := make([]uint32, paletteCount), int32(0)
	for i := int32(0); i < paletteCount; i++ {
		if err := protocol.Varint32(buf, &temp); err != nil {
			return nil, fmt.Errorf("error decoding palette entry: %w", err)
		}
		blocks[i] = uint32(temp)
	}
	return &Palette{values: blocks, size: blockSize}, nil
}

// biomePaletteEncoding implements the encoding of biome palettes to disk.
type biomePaletteEncoding struct{}

func (biomePaletteEncoding) Encode(buf *bytes.Buffer, v uint32) {
	_ = binary.Write(buf, binary.LittleEndian, v)
}
func (biomePaletteEncoding) Decode(buf *bytes.Buffer) (uint32, error) {
	var v uint32
	return v, binary.Read(buf, binary.LittleEndian, &v)
}

// BlockPaletteEncoding implements the encoding of block palettes to disk.
type BlockPaletteEncoding struct {
	block   mapping.Block
	version int32
}

// NewBlockPaletteEncoding returns a new BlockPaletteEncoding using the block and version passed.
func NewBlockPaletteEncoding(block mapping.Block, version int32) BlockPaletteEncoding {
	return BlockPaletteEncoding{block: block, version: version}
}

func (b BlockPaletteEncoding) Encode(buf *bytes.Buffer, v uint32) {
	// Get the block state registered with the runtime IDs we have in the palette of the block storage
	// as we need the name and data value to store.
	state, _ := b.block.RuntimeIDToState(v)
	_ = nbt.NewEncoderWithEncoding(buf, nbt.LittleEndian).Encode(blockupgrader.BlockState{Name: state.Name, Properties: state.Properties, Version: b.version})
}
func (b BlockPaletteEncoding) Decode(buf *bytes.Buffer) (uint32, error) {
	var e blockupgrader.BlockState
	if err := nbt.NewDecoderWithEncoding(buf, nbt.LittleEndian).Decode(&e); err != nil {
		return 0, fmt.Errorf("error decoding block palette entry: %w", err)
	}
	v, ok := b.block.StateToRuntimeID(e)
	if !ok {
		return 0, fmt.Errorf("cannot get runtime ID of block state %v{%+v}", e.Name, e.Properties)
	}
	return v, nil
}

// NewNetworkPersistentEncoding returns a new NetworkPersistentEncoding using the block and version passed.
func NewNetworkPersistentEncoding(block mapping.Block, version int32) NetworkPersistentEncoding {
	return NetworkPersistentEncoding{block: block, version: version}
}

// NetworkPersistentEncoding implements the Chunk encoding for sending over network with a persistent palette.
type NetworkPersistentEncoding struct {
	block   mapping.Block
	version int32
}

func (n NetworkPersistentEncoding) Network() byte { return 1 }
func (n NetworkPersistentEncoding) EncodePalette(buf *bytes.Buffer, p *Palette, _ PaletteEncoding) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}

	enc := nbt.NewEncoderWithEncoding(buf, nbt.NetworkLittleEndian)
	for _, val := range p.values {
		state, _ := n.block.RuntimeIDToState(val)
		_ = enc.Encode(blockupgrader.BlockState{Name: strings.TrimPrefix("minecraft:", state.Name), Properties: state.Properties, Version: n.version})
	}
}
func (n NetworkPersistentEncoding) DecodePalette(buf *bytes.Buffer, blockSize paletteSize, _ PaletteEncoding) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize != 0 {
		err := protocol.Varint32(buf, &paletteCount)
		if err != nil {
			panic(err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	blocks := make([]blockupgrader.BlockState, paletteCount)
	dec := nbt.NewDecoderWithEncoding(buf, nbt.NetworkLittleEndian)
	for i := int32(0); i < paletteCount; i++ {
		if err := dec.Decode(&blocks[i]); err != nil {
			return nil, fmt.Errorf("error decoding block state: %w", err)
		}
	}

	var ok bool
	palette, temp := newPalette(blockSize, make([]uint32, paletteCount)), uint32(0)
	for i, b := range blocks {
		temp, ok = n.block.StateToRuntimeID(blockupgrader.BlockState{Name: "minecraft:" + b.Name, Properties: b.Properties, Version: n.version})
		if !ok {
			return nil, fmt.Errorf("cannot get runtime ID of block state %v{%+v}", b.Name, b.Properties)
		}
		palette.values[i] = temp
	}
	return palette, nil
}

type subChunkVersion8 struct{}

func (subChunkVersion8) EncodeHeader(buf *bytes.Buffer, s *SubChunk, _ cube.Range, _ int) {
	_, _ = buf.Write([]byte{8, byte(len(s.storages))})
}

type subChunkVersion9 struct{}

func (subChunkVersion9) EncodeHeader(buf *bytes.Buffer, s *SubChunk, r cube.Range, ind int) {
	_, _ = buf.Write([]byte{SubChunkVersion, byte(len(s.storages)), uint8(ind + (r[0] >> 4))})
}
