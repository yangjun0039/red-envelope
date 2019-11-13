package network

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"log"
	"red-envelope/configer"
)

type NetProtocol int

const (
	HTTP  NetProtocol = iota
	HTTPS
)

var rootRouter *Router

type Router struct {
	Multiplexer *mux.Router
}

func NewRouter(routes Routes) *Router {
	aRouter := &Router{mux.NewRouter().StrictSlash(true)}
	for _, route := range routes {
		aRouter.Multiplexer.Methods(route.Method).Path(route.Path).Name(route.Name()).Handler(route.HandlerFunc)
		//aRouter.Multiplexer.Methods(http.MethodOptions).Path(route.Path).Name("CrossDomain: " + route.Path).Handler(NewHandlerBuilder().
		//	Adapt(RequesterHandler).Then(RecoveryHandler, CrossOriginHanddler).End(optionsHandler))
	}
	return aRouter
}

type Routes []Route

type Route struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

func (r Route) Name() string {
	return fmt.Sprintf("[%v]%v", r.Method, r.Path)
}

func optionsHandler(requester *Requester, w http.ResponseWriter, r *http.Request) {
	if recorder, ok := w.(*ResponseRecorder); ok {
		recorder.Status = http.StatusOK
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (r *Router) timeoutHandler() http.Handler {
	const timeoutDuration = 45
	failure := NewFailure(http.StatusInternalServerError, responseTimeout)
	bytes, _ := json.Marshal(*failure)
	return http.TimeoutHandler(r, timeoutDuration*time.Second, string(bytes))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Multiplexer.ServeHTTP(w, req)
}

func (r *Router) LoadOnRoutes(pathPrefix string, routes Routes) {
	aSubrouter := r.Multiplexer.PathPrefix(pathPrefix).Subrouter()
	for _, route := range routes {
		aSubrouter.Methods(route.Method).Path(route.Path).Name(route.Name()).Handler(route.HandlerFunc)
		//aSubrouter.Methods(http.MethodOptions).Path(route.Path).Name("CrpssDomain" + pathPrefix + route.Path).Handler(NewHandlerBuilder().
		//	Adapt(RequesterHandler).Then(RecoveryHandler, CrossOriginHanddler).End(optionsHandler))
	}
}

func (r *Router) Startup(protocol NetProtocol, port uint64) {
	rootRouter = r
	server := http.Server{Addr: fmt.Sprintf(":%v", port), Handler: r.timeoutHandler()}
	if protocol == HTTPS {
		key, err := configer.Cfg.GetSection("key")
		if err != nil {
			panic(err)
		}
		certFile := key["certificate_file"]
		keyFile := key["private_key_file"]
		fmt.Println("Startup")
		log.Fatal("http server fatal: ", server.ListenAndServeTLS(certFile, keyFile))

	} else if protocol == HTTP {
		// 添加日志
		fmt.Println("Startup")
		log.Fatal("http server fatal: ", server.ListenAndServe())
	}
}
