package router

import (
	"github.com/ycyxuehan/bingo/bingdb"
	"github.com/ycyxuehan/bingo/logger"
	"github.com/ycyxuehan/bingo/config"
	"github.com/ycyxuehan/bingo/controller"
	"github.com/gorilla/mux"
	"net/http"

)
//HTTPEntry httpentry
type HTTPEntry struct{
	PATH string
	C controller.CtrlInterface
}
//NewHTTPEntry new HTTPEntry
func NewHTTPEntry(path string, c controller.CtrlInterface)*HTTPEntry{
	if path == "" || c == nil{
		return nil
	}
	return &HTTPEntry{
		PATH:path,
		C:c,
	}
}
//Router router
type Router struct {
	conf config.IniConfig
	router *mux.Router
	Logger *logger.BingLog
	DBI bingdb.DBInterface
}

//NewRouter new a router
func NewRouter()*Router{
	r := Router{}
	r.router = mux.NewRouter()
	return &r
}

//Add add a route
func (r *Router)Add(path string, c controller.CtrlInterface){
	c.SetLogger(r.Logger)
	c.SetDBI(r.DBI)
	r.router.Handle(path, c)
}
//Register register
func (r *Router)Register(entries ...*HTTPEntry){
	for _, entry := range entries{
		if entry != nil {
			entry.C.SetLogger(r.Logger)
			entry.C.SetDBI(r.DBI)
			r.router.Handle(entry.PATH, entry.C)
		}
	}
}

//ServeHTTP http response
func (r *Router)ServeHTTP(response http.ResponseWriter, request *http.Request){
}

//Config get a config
func (r *Router)Config(key string)string{
	return r.conf.Get(key)
}

//Router return mux router
func (r *Router)Router()*mux.Router{
	return r.router
}