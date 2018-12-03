package bingo

import (
	"strings"
	"github.com/ycyxuehan/bingo/bingdb/mysql"
	"github.com/ycyxuehan/bingo/bingdb"
	"fmt"
	"net/http"
)

//Bingo bingo
type BingApp struct {
	router *Router
}

var App *BingApp
var DBConnection bingdb.DBInterface
func NewApp()*BingApp{
	return &BingApp{
		router: NewRouter(),
	}
}

func init(){
	App = NewApp()

}

func (b *BingApp)Run(){
	host := BingConf.Get("host")
	port := BingConf.Get("port")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "3000"
	}
	dburi := BingConf.Get("dburi")
	if dburi != "" {
		dbtype := strings.Split(dburi, ":")[0]
	
		switch dbtype {
		case "mysql":
			DBConnection = mysql.New(strings.Split(dburi, "://")[1])
			err := DBConnection.Connect()
			if err != nil {
				Logger.Error("connect to %s error: %s", dburi, err)
				return
			}
			Logger.Info("database: %s connected", dburi)
		}
	}
	Logger.Info("running http server %s:%s", host, port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), b.router)
}

func (b *BingApp)RegisterRoute(entries ...*HTTPEntry){
	b.router.Register(entries...)
}

//Run run
func Run(){
	App.Run()
}

func RegisterRoute(entries ...*HTTPEntry){
	App.RegisterRoute(entries...)
}