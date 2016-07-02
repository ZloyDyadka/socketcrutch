package codec

import "github.com/ZloyDyadka/socketcrutch"

type MsgPackCodec struct{}

func (_ *MsgPackCodec) NewPacket() socketcrutch.Packet {
	return &Packet{}
}

func (_ *MsgPackCodec) ReadMode() int {
	return socketcrutch.BinaryMessage
}

func (_ *MsgPackCodec) WriteMode() int {
	return socketcrutch.BinaryMessage
}
