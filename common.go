package bingo

import (
	"fmt"
	"net/http"
)

//Bingo bingo
type Bingo struct {
	Config IniConfig
}

//New new a bingo
func New(config string)*Bingo{
	b := Bingo{}
	b.Config.Load(config)
	return &b
}


//Run run
func (b *Bingo)Run(r *Router){
	if r == nil {
		return
	}
	host := b.Config.Get("host")
	port := b.Config.Get("port")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r.Router())
}