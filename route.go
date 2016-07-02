// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

type Route struct {
	handlers []HandlerFunc
}

func newRoute(handlers []HandlerFunc) *Route {
	route := &Route{
		handlers: handlers,
	}

	return route
}

// handle runs consistently each handler.
func (r *Route) handle(session *Session, data []byte) error {
	for i, n := 0, len(r.handlers); i < n; i++ {
		if err := r.handlers[i](session, data); err != nil {
			return err
		}
	}

	return nil
}
