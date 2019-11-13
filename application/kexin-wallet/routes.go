package kexin_wallet

import (
	"red-envelope/network"
	"red-envelope/application/kexin-wallet/red-packet"
	"net/http"
)

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/wallet", allRoute())
}

func allRoute() network.Routes {
	routes := network.Routes{
		network.Route{Method: http.MethodPost, Path: "/send-red-envelop", HandlerFunc: network.DefaultAdapter.End(red_packet.SendRedEnvelop)},
		network.Route{Method: http.MethodPost, Path: "/jrmf-send", HandlerFunc: network.DefaultAdapter.End(red_packet.SendJRMFRedEnvelop)},
		network.Route{Method: http.MethodPost, Path: "/jrmf-recive", HandlerFunc: network.DefaultAdapter.End(red_packet.ReciveJRMFRedEnvelip)},
		network.Route{Method: http.MethodGet, Path: "/jrmf-show", HandlerFunc: network.DefaultAdapter.End(red_packet.ShowJRMFRedEnvelop)},
		network.Route{Method: http.MethodGet, Path: "/jrmf-controller", HandlerFunc: network.DefaultAdapter.End(red_packet.ShowJRMF)},
	}
	return routes
}
