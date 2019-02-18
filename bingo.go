package bingo

import (
	"syscall"
	"os/signal"
	"os"
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

func Daemon(){
	pidfile := BingConf.Get("pid")
	if pidfile == "" {
		pidfile = fmt.Sprintf("/var/lib/%s", BingConf.Get("appname"))
	}
	File, err := os.OpenFile(pidfile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Logger.Error("write pid file error: ", err)
		return
	}
	info, _ := File.Stat()
	if info.Size() != 0{
		Logger.Error("pid file is exist.")
		return
	}
	if os.Getppid() != 1 {
		args := append([]string{os.Args[0]}, os.Args[1:]...)
		os.StartProcess(os.Args[0], args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return
	}
	File.WriteString(fmt.Sprint(os.Getpid()))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGUSR2)
	go App.Run()
	for {
		sign := <- signalChan
		switch sign {
		case syscall.SIGUSR2:
			//user customer signal 2
		case os.Interrupt:
			//safety exit
			Exit(File)
			break
		}
	}
}

func Exit(F *os.File){
	err := F.Close()
	if err != nil {
		fmt.Println("close pid file error: %s", err)
	}
	os.Remove(F.Name())

}