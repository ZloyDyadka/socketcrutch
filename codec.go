// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	// BinaryMessage denotes a binary data message.
	BinaryMessage = websocket.BinaryMessage
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = websocket.TextMessage

	// The ID for not ack packet
	WithoutID = -1
)

type PacketType byte

const (
	EventType PacketType = iota
	AckType
)

func (t PacketType) String() string {
	switch t {
	case EventType:
		return "event"
	case AckType:
		return "ack"
	default:
		return fmt.Sprintf("unknown (%d)", t)
	}
}

type (
	// Codec defines the interface that a codec should implement.
	// A codec should be able to marshal/unmarshal messages from the
	// given bytes.
	Codec interface {
		// NewPacket returns the empty new packet.
		NewPacket() Packet
		// Returns read mode, that should be used for this codec.
		ReadMode() int
		// Returns write mode, that should be used for this codec
		WriteMode() int
	}
)
