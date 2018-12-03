package bingo

import (
	"fmt"
	"encoding/json"
	"net/http"
)

//Output output
type Output struct{
	response http.ResponseWriter
}

func NewOutput()*Output{
	return &Output{}
}
//Init init output
func (o *Output)Init(rw http.ResponseWriter){
	o.response = rw
}

//Header set header
func (o *Output)Header(key, val string){
	o.response.Header().Set(key, val)
}

//WriteHeader write status
func (o *Output)WriteHeader(code int){
	o.response.WriteHeader(code)
}

//Write write data
func (o *Output)Write(body []byte)(int, error){
	return o.response.Write(body)
}

//ServeJSON response json
func (o *Output)ServeJSON(obj interface{}, encoding string)error{
	content, err := json.MarshalIndent(obj,"", "\t")
	if err != nil {
		http.Error(o.response, err.Error(), http.StatusInternalServerError)
		return err
	}
	return o.Serve(content, "application/json", encoding)
}

//ServeString response string
func (o *Output)ServeString(body, encoding string)error{
	return o.Serve([]byte(body), "text/plain", encoding)
}

//Serve response
func (o *Output)Serve(content []byte,contentType, encoding string)error{
	o.Header("Content-Type", fmt.Sprintf("%s; %s", contentType, encoding))
	_, err := o.Write(content)
	return err
}

//Reset reset
func (o *Output)Reset(){
	o.response = nil
}