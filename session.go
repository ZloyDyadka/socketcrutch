// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

import (
	"net"
	"net/http"
	"sync"
)

type (
	// Session represents an active SocketCrutch session.
	Session struct {
		connection Conn
		id         string
		details    map[string]interface{}
		mutex      sync.RWMutex
		sc         *SocketCrutch
		request    *http.Request
	}
)

func newSession(c Conn, id string, req *http.Request, sc *SocketCrutch) *Session {
	s := &Session{
		connection: c,
		id:         id,
		request:    req,
		details:    make(map[string]interface{}),
		sc:         sc,
	}

	go s.handleMessagesLoop()

	return s
}

// Codec returns the established codec.
func (s *Session) codec() Codec {
	return s.connection.Codec()
}

// Request returns the first HTTP request when established connection.
func (s *Session) Request() *http.Request {
	return s.request
}

// RemoteAddr returns the remote network address.
func (s *Session) RemoteAddr() net.Addr {
	return s.connection.RemoteAddr()
}

// Request returns the session ID.
func (s *Session) ID() string {
	return s.id
}

// Close closes the connection.
func (s *Session) Close() error {
	return s.connection.Close()
}

// SetDetail sets a value of the details.
// SetDetail is thread-safe.
func (s *Session) SetDetail(key string, value interface{}) {
	s.mutex.Lock()
	s.details[key] = value
	s.mutex.Unlock()
}

// GetDetail returns a value of the details.
// GetDetail is thread-safe
func (s *Session) GetDetail(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, ok := s.details[key]

	return val, ok
}

// Publish publishes the event with provided data.
// The data will be automatically packed.
func (s *Session) Publish(event string, payload []byte) error {
	packet := s.formatEventPacket(event, payload)

	rawData, err := packet.Marshal()

	if err != nil {
		return err
	}

	if err := s.connection.Send(rawData); err != nil {
		return err
	}

	return nil
}

func (s *Session) formatEventPacket(event string, payload []byte) Packet {
	packet := s.codec().NewPacket()

	packet.SetType(EventType)
	packet.SetID(WithoutID)
	packet.SetEventName(event)
	packet.SetPayload(payload)

	return packet
}

// handleMessagesLoop reads messages from the connection and process them.
func (s *Session) handleMessagesLoop() {
	for {
		msg, ok := s.connection.NextMessage()

		if !ok {
			return
		}

		s.processMessage(msg)
	}
}

func (s *Session) processMessage(msg []byte) {
	log.Debug("unmarshall the message")

	packet := s.codec().NewPacket()

	if err := packet.Unmarshal(msg); err != nil {
		log.Debug("failed unmarshal the message: %v", err)
		return
	}

	log.Debugf("the packet type %s", packet.Type().String())

	// TODO: ack packet type
	switch packet.Type() {
	case EventType:
		s.processEventPacket(packet)
	default:
		log.Debugf("unknown the packet type %d", packet.Type())
	}
}

func (s *Session) processEventPacket(packet Packet) {
	if packet.EventName() == "" {
		log.Debug("the packet name is empty")
		return
	}

	route := s.sc.router.getRoute(packet.EventName())

	if route == nil {
		log.Debugf("route %s not found", packet.EventName())
		return
	}

	if err := route.handle(s, packet.Payload()); err != nil {
		log.Debugf("an error has occured %s", err.Error())

		s.sc.errorHandler(s, err)
	}
}

// Coming soon
func (s *Session) processAckPacket(packet Packet) {

}
