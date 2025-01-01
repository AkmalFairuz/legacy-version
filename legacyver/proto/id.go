package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	ID766 = 766 // v1.21.50
	ID748 = 748 // v1.21.40
	ID729 = 729 // v1.21.30
	ID712 = 712 // v1.21.20
	ID686 = 686 // v1.21.2
	ID685 = 685 // v1.21.0
	ID671 = 671 // v1.20.80
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
