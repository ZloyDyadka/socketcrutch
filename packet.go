package socketcrutch

type (
	// Packet is a container for the raw data and therefore
	// it holds the event name, pointer to data and id of the event type.
	Packet interface {
		// EventName returns the event name from packet.
		EventName() string
		// SetEventName sets the name for this packet.
		SetEventName(string)
		// Type returns the type from the packet.
		Type() PacketType
		// SetType sets the name for this packet.
		SetType(PacketType)
		// Payload returns the payload from the packet.
		Payload() []byte
		// SetPayload sets the payload for this packet.
		SetPayload([]byte)
		// ID returns the id from the packet.
		ID() int
		// SetID sets the id for this packet.
		SetID(int)
		//Unmarshal parses the encoded data structure and stores the result in the Packet.
		Unmarshal([]byte) error
		//Marshal returns the data structure encoding of Packet.
		Marshal() ([]byte, error)
	}
)
