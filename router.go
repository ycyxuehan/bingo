package bingo

import (
	"github.com/gorilla/mux"


)
//HTTPEntry httpentry
type HTTPEntry struct{
	PATH string
	C CtrlInterface
}
//NewHTTPEntry new HTTPEntry
func NewHTTPEntry(path string, c CtrlInterface)*HTTPEntry{
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
	mux.Router
}

//NewRouter new a router
func NewRouter()*Router{
	r := Router{}
	
	// r.router = mux.NewRouter()
	return &r
}

//Add add a route
func (r *Router)Add(path string, c CtrlInterface){
	c.SetThis(c)
	r.Handle(path, c)
}
//Register register
func (r *Router)Register(entries ...*HTTPEntry){
	for _, entry := range entries{
		if entry != nil {
			entry.C.SetThis(entry.C)
			r.Handle(entry.PATH, entry.C)
		}
	}
}

//ServeHTTP http response
// func (r *Router)ServeHTTP(response http.ResponseWriter, request *http.Request){
// 	r.ServeHTTP(response, request)
// }

//Config get a config
func (r *Router)Config(key string)string{
	return BingConf.Get(key)
}

//Router return mux router
// func (r *Router)Router()*mux.Router{
// 	return r.router
// }