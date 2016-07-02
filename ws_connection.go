// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

import (
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

var (
	ErrSend = errors.New("the message buffer full or the connection is closed")
)

type Conn interface {
	// Send sends the given raw data to the websocket connection.
	Send([]byte) error

	// Close closes the underlying websocket connection.
	Close() error

	// NextMessage returns the next data message received from the connection.
	NextMessage() ([]byte, bool)

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr

	// Codec returns the established codec.
	Codec() Codec
}

// connection represents a WebSocket connection.
type wsConnection struct {
	ws *websocket.Conn
	sc *SocketCrutch

	inboxMessages  chan []byte
	outboxMessages chan []byte

	done chan struct{}
}

func newConnection(ws *websocket.Conn, sc *SocketCrutch) *wsConnection {
	conn := &wsConnection{
		ws: ws,
		sc: sc,

		inboxMessages:  make(chan []byte, sc.messageBufferSize),
		outboxMessages: make(chan []byte, sc.messageBufferSize),
	}

	conn.configure()

	return conn
}

func (c *wsConnection) configure() {
	c.ws.SetReadLimit(c.sc.readLimit)
	c.ws.SetPongHandler(c.pongHandler)
}

// run starts reading and writing the goroutines.
func (c *wsConnection) run() {
	go c.writeMessagesLoop(c.Codec().WriteMode())
	go c.readMessagesLoop(c.Codec().ReadMode())
}

func (c *wsConnection) wait() {
	<-c.done
}

func (c *wsConnection) Codec() Codec {
	return c.sc.codec
}

func (c *wsConnection) RemoteAddr() net.Addr {
	return c.ws.RemoteAddr()
}

func (c *wsConnection) Send(packet []byte) error {
	select {
	case c.outboxMessages <- packet:
		return nil
	default:
		return ErrSend
	}
}

func (c *wsConnection) Close() error {
	return c.ws.Close()
}

func (c *wsConnection) NextMessage() ([]byte, bool) {
	data, ok := <-c.inboxMessages

	return data, ok
}

func (c *wsConnection) pongHandler(string) error {
	log.Debugf("received pong message")

	c.ws.SetReadDeadline(time.Now().Add(c.sc.pingTimeout))

	return nil
}

// write writes a message with the given message type and payload.
func (c *wsConnection) write(mode int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(pingTimeout))

	return c.ws.WriteMessage(mode, payload)
}

// readMessagesLoop reads messages from the websocket connection.
func (c *wsConnection) readMessagesLoop(mode int) {
	// When disconnect
	defer func() {
		close(c.done)
		close(c.outboxMessages)
		close(c.inboxMessages)
	}()

	for {
		// Reset the read deadline.
		c.ws.SetReadDeadline(time.Now().Add(c.sc.pingTimeout))

		messageType, message, err := c.ws.ReadMessage()

		if err != nil {
			isUnexpectedError := websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure,
				websocket.CloseGoingAway, websocket.CloseNoStatusReceived)

			if err != io.EOF && isUnexpectedError {
				log.Error("failed to read data from websocket: %v", err)
			}

			return
		}

		switch messageType {
		case mode:
			log.Debug("received message with a payload")
			c.inboxMessages <- message
		case websocket.CloseMessage:
			log.Debug("got close message")
			return
		default:
			log.Debug("got wrong message, message dropped")
		}
	}
}

// writeMessagesLoop writes messages from the session to the websocket connection.
func (c *wsConnection) writeMessagesLoop(mode int) {
	var emptyMessage []byte

	pingTicker := time.NewTicker(c.sc.pingInterval)
	defer pingTicker.Stop()

	for {
		select {
		case message, ok := <-c.outboxMessages:
			if !ok {
				log.Debug("write close message")

				c.write(websocket.CloseMessage, emptyMessage)
				return
			}

			log.Debugf("write message with a payload")

			if err := c.write(mode, message); err != nil {
				log.Errorf("an error has occurred in write binary message: %v", err)
				return
			}

		case <-pingTicker.C:
			log.Debug("write ping message")

			if err := c.write(websocket.PingMessage, emptyMessage); err != nil {
				log.Errorf("an error has occurred in write ping message: %v", err)
				return
			}
		}
	}
}
