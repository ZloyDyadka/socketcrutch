// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

type (
	HandlerFunc    func(*Session, []byte) error
	MiddlewareFunc func(*Session, []byte)

	Router struct {
		separator string
		store     RouteStore

		RouteGroup
	}

	RouteStore interface {
		Get(string) *Route
		Set(string, *Route)
	}
)

func newRouter(separator string) *Router {
	router := &Router{
		store:     newRouteStore(),
		separator: separator,
	}

	router.RouteGroup = newRouteGroup("", router, []HandlerFunc{})

	return router
}

// SetRouteSeparator sets the character as a separator between the route parts.
func (r *Router) SetRouteSeparator(separator string) {
	r.separator = separator
}

func (r *Router) addRoute(path string, route *Route) {
	r.store.Set(path, route)
}

func (r *Router) getRoute(path string) *Route {
	return r.store.Get(path)
}
