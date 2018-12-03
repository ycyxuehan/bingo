package bingo

import (
	"github.com/ycyxuehan/bingo/bingdb/mysql"
	"github.com/ycyxuehan/bingo/router"
	"github.com/ycyxuehan/bingo/logger"
	"github.com/ycyxuehan/bingo/config"
	"fmt"
	"net/http"
)

//Bingo bingo
type Bingo struct {
	Config config.IniConfig
	Logger *logger.BingLog
}

//New new a bingo
func New(config string)*Bingo{
	b := Bingo{}
	b.Config.Load(config)
	b.Logger = logger.New(&b.Config)
	return &b
}


//Run run
func (b *Bingo)Run(r *router.Router){
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
	r.Logger = b.Logger
	b.Logger.Start()
	if dburi := b.Config.Get("dburi"); dburi != "" {
		dbtype := b.Config.Get("dbtype")
		switch dbtype {
		case "mysql":
			dbi := mysql.New(dburi)
			err := dbi.Connect()
			if err == nil {
				b.Logger.Info("database: %s connected", dburi)
				r.DBI = dbi
			}else {
				b.Logger.Error("connect to %s error: %s", dburi, err)
				return
			}
		}
	}
	b.Logger.Info("running http server %s:%s", host, port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r.Router())
}