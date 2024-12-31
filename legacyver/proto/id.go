package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	ID766 = 766
	ID748 = 748
	ID729 = 729
	ID712 = 712
	ID686 = 686
	ID685 = 685
)

func IsProtoGTE(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() >= proto
}

func IsProtoLTE(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() <= proto
}

func IsProtoLT(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() < proto
}

func IsProtoGT(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() > proto
}

func FetchProtoID(io protocol.IO) int32 {
	return io.(IO).ProtocolID()
}
