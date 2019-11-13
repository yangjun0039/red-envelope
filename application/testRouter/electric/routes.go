package electric

import (
	"net/http"
	"red-envelope/network"
	"red-envelope/application/testRouter/electric/activity"
)

var Routes = network.Routes{
	// 获取账户信息
	network.Route{http.MethodGet, "/user-info",network.TestAdapter.End(activity.UserInfo)},

	network.Route{http.MethodGet, "/acc-info",network.TestAdapter.End(activity.AccInfo)},
}
