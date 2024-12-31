package proto

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TexturePackInfo represents a texture pack's info sent over network. It holds information about the
// texture pack such as its name, description and version.
type TexturePackInfo struct {
	// UUID is the UUID of the texture pack. Each texture pack downloaded must have a different UUID in
	// order for the client to be able to handle them properly.
	UUID uuid.UUID
	// Version is the version of the texture pack. The client will cache texture packs sent by the server as
	// long as they carry the same version. Sending a texture pack with a different version than previously
	// will force the client to re-download it.
	Version string
	// Size is the total size in bytes that the texture pack occupies. This is the size of the compressed
	// archive (zip) of the texture pack.
	Size uint64
	// ContentKey is the key used to decrypt the behaviour pack if it is encrypted. This is generally the case
	// for marketplace texture packs.
	ContentKey string
	// SubPackName ...
	SubPackName string
	// ContentIdentity is another UUID for the resource pack, and is generally set for marketplace texture
	// packs. It is also required for client-side validations when the resource pack is encrypted.
	ContentIdentity string
	// HasScripts specifies if the texture packs has any scripts in it. A client will only download the
	// behaviour pack if it supports scripts, which, up to 1.11, only includes Windows 10.
	HasScripts bool
	// AddonPack specifies if the texture pack is from an addon.
	AddonPack bool
	// RTXEnabled specifies if the texture pack uses the raytracing technology introduced in 1.16.200.
	RTXEnabled bool
	// DownloadURL is a URL that the client can use to download the pack instead of the server sending it in
	// chunks, which it will continue to do if this field is left empty.
	DownloadURL string
}

func (x *TexturePackInfo) ToLatest() protocol.TexturePackInfo {
	return protocol.TexturePackInfo{
		UUID:            x.UUID,
		Version:         x.Version,
		Size:            x.Size,
		ContentKey:      x.ContentKey,
		SubPackName:     x.SubPackName,
		ContentIdentity: x.ContentIdentity,
		HasScripts:      x.HasScripts,
		AddonPack:       x.AddonPack,
		RTXEnabled:      x.RTXEnabled,
		DownloadURL:     x.DownloadURL,
	}
}

// Marshal encodes/decodes a TexturePackInfo.
func (x *TexturePackInfo) Marshal(r protocol.IO) {
	if IsProtoGTE(r, ID766) {
		r.UUID(&x.UUID)
	} else {
		if IsReader(r) {
			uuidStr := ""
			r.String(&uuidStr)
			x.UUID = uuid.MustParse(uuidStr)
		} else {
			uuidStr := x.UUID.String()
			r.String(&uuidStr)
		}
	}
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
	r.Bool(&x.AddonPack)
	r.Bool(&x.RTXEnabled)
	if IsProtoGTE(r, ID748) {
		r.String(&x.DownloadURL)
	}
}
