package bingo

import (
	"net/url"
	"net/http"
	"strings"
)
//Params params
type Params map[string][]string

//Input input
type Input struct {
	Param url.Values
	RequestBody []byte
	RequestBodyLength int
	request *http.Request
}

//Init init
func (i *Input)Init(r *http.Request){
	i.request = r
	p, e := url.ParseQuery(r.URL.RawQuery)
	if e == nil {
		i.Param = p
	}
	if i.IsPost() {
		reader := r.Body
		i.RequestBodyLength, _ = reader.Read(i.RequestBody)
	}
}

//Is returns boolean of this request is on given method
func (i *Input)Is(method string)bool{
	return strings.ToUpper(i.request.Method) == strings.ToUpper(method)
}

//IsPost is this a POST method request
func (i *Input)IsPost()bool{
	return i.Is("POST")
}

//Cookie get the cookie
func (i *Input)Cookie(key string)string{
	ck, err := i.request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

//Header get the header
func (i *Input)Header(key string)string {
	return i.request.Header.Get(key)
}

//URI the request uri
func (i *Input)URI()string{
	return i.request.RequestURI
}

//URL request url path
func (i *Input)URL()string{
	return i.request.URL.Path
}

//Method request method
func (i *Input)Method()string{
	return i.request.Method
}

//Reset reset
func (i *Input)Reset(){
	i.Param = url.Values{}
	i.RequestBody = []byte{}
	i.RequestBodyLength = 0
	i.request = nil
}