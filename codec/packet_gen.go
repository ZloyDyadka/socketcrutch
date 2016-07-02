package codec

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Packet) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "type":
			z.pType, err = dc.ReadByte()
			if err != nil {
				return
			}
		case "id":
			z.id, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "event":
			z.eventName, err = dc.ReadString()
			if err != nil {
				return
			}
		case "payload":
			z.payload, err = dc.ReadBytes(z.payload)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Packet) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "type"
	err = en.Append(0x84, 0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteByte(z.pType)
	if err != nil {
		return
	}
	// write "id"
	err = en.Append(0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.id)
	if err != nil {
		return
	}
	// write "event"
	err = en.Append(0xa5, 0x65, 0x76, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.eventName)
	if err != nil {
		return
	}
	// write "payload"
	err = en.Append(0xa7, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.payload)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Packet) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "type"
	o = append(o, 0x84, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendByte(o, z.pType)
	// string "id"
	o = append(o, 0xa2, 0x69, 0x64)
	o = msgp.AppendInt(o, z.id)
	// string "event"
	o = append(o, 0xa5, 0x65, 0x76, 0x65, 0x6e, 0x74)
	o = msgp.AppendString(o, z.eventName)
	// string "payload"
	o = append(o, 0xa7, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	o = msgp.AppendBytes(o, z.payload)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Packet) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "type":
			z.pType, bts, err = msgp.ReadByteBytes(bts)
			if err != nil {
				return
			}
		case "id":
			z.id, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "event":
			z.eventName, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "payload":
			z.payload, bts, err = msgp.ReadBytesBytes(bts, z.payload)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Packet) Msgsize() (s int) {
	s = 1 + 5 + msgp.ByteSize + 3 + msgp.IntSize + 6 + msgp.StringPrefixSize + len(z.eventName) + 8 + msgp.BytesPrefixSize + len(z.payload)
	return
}
