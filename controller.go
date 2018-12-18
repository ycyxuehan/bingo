package bingo

import (
	// "reflect"
	"net/http"
)
//Result a result template
type Result struct{
	ResultCode int `json:"ResultCode"`
	ResultString string `json:"ResultString"`
	ResultData interface{} `json:"ResultData"`
}

//CtrlInterface controller interface
type CtrlInterface interface {
	Init()
	Release()
	SetContext(*Context)
	Context()*Context
	Config(string)string
	SetThis(CtrlInterface)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Get()
	Post()
	Put()
	Delete()
}

//Controller controller
type Controller struct {
	Ctx *Context
	this CtrlInterface
}

//Init initail the controller
func (c *Controller)Init(){
	c.Ctx = NewContext()
}

//SetThis set this
func (c *Controller)SetThis(ci CtrlInterface){
	c.this = ci
}

//Release release the controller
func (c *Controller)Release(){

}

//SetContext set controller context
func (c *Controller)SetContext(ctx *Context){
	c.Ctx = ctx
}

//Context get context
func (c *Controller)Context()*Context{
	return c.Ctx
}

//Config get a config
func (c *Controller)Config(key string)string{
	return BingConf.Get(key)
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
	Logger.Info("[%s] %s",r.Method, c.Ctx.URL())
	switch r.Method {
	case "GET":
		c.this.Get()
		break;
	case "PUT":
		c.this.Put()
		break;
	case "POST":
		c.this.Post()
		break
	case "DELETE":
		c.this.Delete()
		break
	default:
		c.NotFound()
	}
}

