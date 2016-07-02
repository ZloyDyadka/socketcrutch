// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

import (
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pingTimeout = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingInterval = 30 * time.Second

	// Maximum message size allowed from peer.
	readLimit = 1024

	// Maximum read buffer size.
	readBufferSize = 4096

	// Maximum write buffer size.
	writeBufferSize = 4096

	// Maximum count connections.
	maxConnections = 1024

	// The max amount of messages that can be in a sessions buffer before it starts dropping them.
	messageBufferSize = 256

	//
	routeSeparator = "."
)

type config struct {
	pingTimeout       time.Duration
	pingInterval      time.Duration
	readLimit         int64
	maxConnections    int32
	messageBufferSize int
	readBufferSize    int
	writeBufferSize   int
	routeSeparator    string
}

func getDefaultConfig() *config {
	return &config{
		pingInterval:      pingInterval,
		pingTimeout:       pingTimeout,
		readLimit:         readLimit,
		messageBufferSize: messageBufferSize,
		maxConnections:    maxConnections,
		routeSeparator:    routeSeparator,
		readBufferSize:    readBufferSize,
		writeBufferSize:   writeBufferSize,
	}
}

// SetPingTimeout sets ping timeout.
// By default ping timeout is 30 seconds.
func (c *config) SetPingTimeout(t time.Duration) *config {
	c.pingTimeout = t

	return c
}

// SetPingInterval sets ping interval.
// By default ping interval is 30 seconds.
func (c *config) SetPingInterval(t time.Duration) *config {
	c.pingInterval = t

	return c
}

// SetReadBufferSize set specify I/O buffer sizes.
// The I/O buffer sizes do not limit the size of the messages that can be received.
// By default read buffer size is 4096
func (c *config) SetReadBufferSize(size int) *config {
	c.readBufferSize = size

	return c
}

// SetWriteBufferSize set specify I/O buffer sizes.
// The I/O buffer sizes do not limit the size of the messages that can be sent.
// By default write buffer size is 4096
func (c *config) SetWriteBufferSize(size int) *config {
	c.writeBufferSize = size

	return c
}

// SetReadLimit sets the maximum size for a message read from the peer.
// By default read limit size is 1024.
func (c *config) SetReadLimit(limit int64) *config {
	c.readLimit = limit

	return c
}

// SetMessageBufferSize sets the maximum amount of messages in the buffer.
// The max amount of messages that can be in a sessions buffer before it starts dropping them.
// By default buffer size is 256.
func (c *config) SetMessageBufferSize(size int) *config {
	c.messageBufferSize = size

	return c
}

// SetMaxConnections sets maximum numbers of connections.
// By default amount connections is 1024.
func (c *config) SetMaxConnections(max int32) *config {
	c.maxConnections = max

	return c
}
