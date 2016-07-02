// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

// RouteGroup represents a group of routes that share the same path prefix.
type RouteGroup struct {
	prefix   string
	router   *Router
	handlers []HandlerFunc
}

// newRouteGroup creates a new RouteGroup with the given path prefix, router, and handlers.
func newRouteGroup(prefix string, router *Router, handlers []HandlerFunc) RouteGroup {
	if prefix != "" {
		prefix = prefix + router.separator
	}

	return RouteGroup{
		prefix:   prefix,
		router:   router,
		handlers: handlers,
	}
}

// Route adds a route to the router with the given route path and handlers.
func (g *RouteGroup) Route(path string, handlers ...HandlerFunc) {
	g.addRoute(g.prefix+path, handlers...)
}

// Use registers one or multiple handlers to the current route group.
// These handlers will be shared by all routes belong to this group and its subgroups.
func (g *RouteGroup) Use(handlers ...HandlerFunc) {
	g.handlers = append(g.handlers, handlers...)
}

// addRoute registers the route, the handlers to the router.
// The handlers will be combined with the handlers of the route group.
func (g *RouteGroup) addRoute(path string, handlers ...HandlerFunc) {
	g.router.addRoute(path, newRoute(mergeHandlers(g.handlers, handlers)))
}

// Group creates a RouteGroup with the given route path prefix and handlers.
// The new group will combine the existing path prefix with the new one.
// If no handler is provided, the new group will inherit the handlers registered
// with the current group.
func (g *RouteGroup) Group(prefix string, handlers ...HandlerFunc) *RouteGroup {
	if len(handlers) == 0 {
		handlers = make([]HandlerFunc, len(g.handlers))
		copy(handlers, g.handlers)
	}

	return newRouteGroup(g.prefix+prefix, g.router, handlers)
}

// mergeHandlers merges two lists of handlers into a new list.
func mergeHandlers(h1 []HandlerFunc, h2 []HandlerFunc) []HandlerFunc {
	hh := make([]HandlerFunc, len(h1)+len(h2))

	copy(hh, h1)
	copy(hh[len(h1):], h2)

	return hh
}
