package bingo

import (
	"net/http"
)

//CtrlInterface controller interface
type CtrlInterface interface {
	Init()
	Release()
	SetContext(Context)
	Context()*Context
	Config(string)string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	SetLogger(*BingLog)
}

//Controller controller
type Controller struct {
	Ctx Context
	conf IniConfig
	Logger *BingLog
}

//Init initail the controller
func (c *Controller)Init(){

}

//Release release the controller
func (c *Controller)Release(){

}

//SetContext set controller context
func (c *Controller)SetContext(ctx Context){
	c.Ctx = ctx
}

//Context get context
func (c *Controller)Context()*Context{
	return &c.Ctx
}

//Config get a config
func (c *Controller)Config(key string)string{
	return c.conf.Get(key)
}

//Get response http get request
func (c *Controller)Get(){

}
//Put response http put request
func (c *Controller)Put(){

}
//Post response http post request
func (c *Controller)Post(){

}
//Delete response http delete request
func (c *Controller)Delete(){

}
//NotFound not found page
func (c *Controller)NotFound(){
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Ctx.ResponseWriter.Write([]byte("not found"))
}
//ServeHTTP serve http
func (c *Controller)ServeHTTP(w http.ResponseWriter, r *http.Request){
	c.Init()
	c.Ctx.Init(r, w)
	switch r.Method {
	case "GET":
		c.Get()
		break;
	case "PUT":
		c.Put()
		break;
	case "POST":
		c.Post()
		break
	case "DELETE":
		c.Delete()
		break
	default:
		c.NotFound()
	}
}

//SetLogger set logger
func (c *Controller)SetLogger(bl *BingLog){
	c.Logger = bl
}