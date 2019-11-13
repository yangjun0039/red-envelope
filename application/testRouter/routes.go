package testRouter

import (
	"red-envelope/network"
	"red-envelope/application/testRouter/electric"
)

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/test", allRoute(electric.Routes))
}

func allRoute(route ...[]network.Route) network.Routes {
	routes := network.Routes{}
	for _, v := range route {
		routes = append(routes, v...)
	}
	return routes
}


