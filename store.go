// Copyright 2016 Ilya Galimyanov  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.SE OR OTHER DEALINGS IN THE SOFTWARE.

package socketcrutch

type routeStore struct {
	routesMap map[string]*Route
}

func newRouteStore() *routeStore {
	return &routeStore{
		routesMap: make(map[string]*Route, 0),
	}
}

func (s *routeStore) Set(path string, route *Route) {
	s.routesMap[path] = route
}

func (s *routeStore) Get(path string) *Route {
	return s.routesMap[path]
}
