package codec

import "github.com/ZloyDyadka/socketcrutch"

// packet is container for raw data and therefore
// holds the event name, the pointer to payload,
// the event type id, the packet id.
type Packet struct {
	pType     byte   `msg:"type"`
	id        int    `msg:"id"`
	eventName string `msg:"event"`
	payload   []byte `msg:"payload"`
}

func (p *Packet) EventName() string {
	return p.eventName
}

func (p *Packet) SetEventName(name string) {
	p.eventName = name
}

func (p *Packet) Type() socketcrutch.PacketType {
	return socketcrutch.PacketType(p.pType)
}

func (p *Packet) SetType(pType socketcrutch.PacketType) {
	p.pType = byte(pType)
}

func (p *Packet) Payload() []byte {
	return p.payload
}

func (p *Packet) SetPayload(payload []byte) {
	p.payload = payload
}

func (p *Packet) ID() int {
	return p.id
}

func (p *Packet) SetID(id int) {
	p.id = id
}

func (p *Packet) Unmarshal(data []byte) error {
	_, err := p.UnmarshalMsg(data)

	return err
}

func (p *Packet) Marshal() ([]byte, error) {
	out, err := p.MarshalMsg(nil)

	return out, err
}
