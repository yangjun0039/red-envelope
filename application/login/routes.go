package login

import (
	"red-envelope/network"
	"net/http"
)


var loginRputer = network.Routes{
	{http.MethodPost, "", network.LoginHandler.End(Login)},
}

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/login", loginRputer)
}


