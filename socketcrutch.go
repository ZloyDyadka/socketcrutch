// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.
package socketcrutch

import (
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"net/http"
	"sync/atomic"
)

type (
	// SocketCrutch is server of the websocket
	SocketCrutch struct {
		upgrader websocket.Upgrader

		handshakeHandler  HandshakeHandler
		connectHandler    SessionHandler
		disconnectHandler SessionHandler
		errorHandler      ErrorHandler

		connectCounter int32

		codec Codec

		router Router

		config
	}

	SessionHandler   func(*Session)
	ErrorHandler     func(*Session, error)
	HandshakeHandler func(http.ResponseWriter, *http.Request) bool
)

// New returns a SocketCrutch instance with default Upgrader, Config, Handlers.
func New() *SocketCrutch {
	defaultConfig := *getDefaultConfig()

	s := &SocketCrutch{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  defaultConfig.readBufferSize,
			WriteBufferSize: defaultConfig.writeBufferSize,
		},
		config: defaultConfig,
	}

	s.router = newRouter(defaultConfig.routeSeparator)

	s.SetErrorHandler(func(*Session, error) {})

	return s
}

func (sc *SocketCrutch) GetRouter() *Router {
	return sc.router
}

// SetCodec sets the codec of the connection
// to the supplied implementation of the Codec interface.
func (sc *SocketCrutch) SetCodec(p Codec) *SocketCrutch {
	sc.codec = p

	return sc
}

// SetConnectHandler sets the callback,
// that is called when a websocket connection was successfully established.
func (sc *SocketCrutch) SetConnectHandler(fn SessionHandler) *SocketCrutch {
	sc.connectHandler = fn

	return sc
}

// SetDisconnectHandler sets the callback, that is called when the connection is closed.
func (sc *SocketCrutch) SetDisconnectHandler(fn SessionHandler) *SocketCrutch {
	sc.disconnectHandler = fn

	return sc
}

// SetErrorHandler sets the callback, that is called when in connection error has occurred.
func (sc *SocketCrutch) SetErrorHandler(fn ErrorHandler) *SocketCrutch {
	sc.errorHandler = fn

	return sc
}

// SetHandshakeHandler sets the callback for handshake verification.
func (sc *SocketCrutch) SetHandshakeHandler(fn HandshakeHandler) *SocketCrutch {
	sc.handshakeHandler = fn

	return sc
}

// SetCheckOrigin sets the callback for the request Origin header validation.
func (sc *SocketCrutch) SetOriginHandler(f func(*http.Request) bool) *SocketCrutch {
	sc.upgrader.CheckOrigin = f

	return sc
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (sc *SocketCrutch) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Debugf("handling %s http request %s", req.Method, req.URL)

	if sc.handshakeHandler != nil {
		if ok := sc.handshakeHandler(rw, req); !ok {
			log.Debug("custom handshake failed")

			http.Error(rw, "handshake failed", http.StatusForbidden)
			return
		}
	}

	if atomic.LoadInt32(&sc.connectCounter)+1 > sc.config.maxConnections {
		log.Debug("too many connections")

		http.Error(rw, "too many connections", http.StatusServiceUnavailable)

		return
	}

	// Upgrade websocket connection.
	protocols := websocket.Subprotocols(req)
	var responseHeader http.Header = nil

	if len(protocols) > 0 {
		responseHeader = http.Header{"Sec-Websocket-Protocol": {protocols[0]}}
	}

	ws, err := sc.upgrader.Upgrade(rw, req, responseHeader)

	if err != nil {
		log.Debugf("failed to upgrade websocket: %v", err)

		if _, ok := err.(websocket.HandshakeError); !ok {
			sc.errorHandler(nil, err)

			log.Errorf("unknown error while upgrade websocket: %v", err)
		}

		return
	}

	sc.createClient(ws, req)
}

func (sc *SocketCrutch) createClient(ws *websocket.Conn, req *http.Request) {
	id := xid.New().String()

	conn := newConnection(ws, sc)
	session := newSession(conn, id, req, sc)

	atomic.AddInt32(&sc.connectCounter, 1)
	defer atomic.AddInt32(&sc.connectCounter, -1)

	log.Debugf("the client session #%s created, the current number of connections %d", id, sc.connectCounter)

	if sc.connectHandler != nil {
		sc.connectHandler(session)
	}

	// Start reading and writing the goroutines.
	conn.run()
	// Wait disconnect
	conn.wait()

	if sc.disconnectHandler != nil {
		sc.disconnectHandler(session)
	}

	log.Debugf("the client session #%s disconnected, the current number of connections %d", id, sc.connectCounter)
}
