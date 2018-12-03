package bingo

import (
	"net/http"
)

//Context context
type Context struct{
	Input *Input
	Output *Output
	Request *http.Request
	ResponseWriter http.ResponseWriter
	Encoding string
}

func NewContext()*Context{
	c := Context{
		Input: NewInput(),
		Output:NewOutput(),
		Request: nil,
		ResponseWriter: nil,
	}
	return &c
}

//Init init context
func (c *Context)Init(r *http.Request, w http.ResponseWriter){
	c.Reset()
	c.Request = r
	c.ResponseWriter = w
	c.Input.Init(r)
	c.Output.Init(w)
}

//URI request uri
func (c *Context)URI()string {
	return c.Input.URI()
}

//URL request path
func (c *Context)URL()string{
	return c.Input.URL()
}

//Method request method
func (c *Context)Method()string{
	return c.Input.Method()
}

//Param request param
func (c *Context)Param(key string)string{
	if key == "" {
		return ""
	}
	return c.Input.Param.Get(key)
}

//Header get request header
func (c *Context)Header(key string)string{
	return c.Input.Header(key)
}

//Cookie get request cookie
func (c *Context)Cookie(key string)string{
	return c.Input.Cookie(key)
}

//ResponseEncoding get response encoding
func (c *Context)ResponseEncoding()string{
	return c.Header("Accept-Encoding")
}

//Serve response
func (c *Context)Serve(content []byte, contentType string)error{
	return c.Output.Serve(content, contentType, c.Encoding)
}

//ServeJSON response json
func (c *Context)ServeJSON(obj interface{})error{
	return c.Output.ServeJSON(obj, c.Encoding)
}

//ServeString response string
func (c *Context)ServeString(body string)error{
	return c.Output.ServeString(body, c.Encoding)
}

//Reset reset
func (c *Context)Reset(){
	c.Input.Reset()
	c.Output.Reset()
}