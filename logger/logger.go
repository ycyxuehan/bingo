package logger

import (
	"github.com/ycyxuehan/bingo/config"
	"bufio"
	"os"
	"fmt"
	"time"
	"strings"
)

//LogMode log mode
type LogMode int

//Set set mode of log
func (mode LogMode)Set(m string){
	switch strings.ToUpper(m) {
	case "CONSOLE":
		mode = CONSOLE
		break
	case "FILE":
		mode = FILE
		break
	case "SYSLOG":
		mode = SYSLOG
		break
	default:
		mode = CONSOLE
	}
}

//LogMessage log
type LogMessage struct{
	Message string
	Level LogLevel
}

//BingLog log module for bing go
type BingLog struct {
	Level LogLevel
	Mode LogMode
	Path string
	Daily bool
	APPName string
	pool chan LogMessage
}

//LogLevel level of log
type LogLevel int

//Set set level of log
func (level LogLevel)Set(l string){
	switch strings.ToUpper(l) {
	case "EMERGENCY":
		level = EMERGENCY
		break
	case "ALERT":
		level = ALERT
		break
	case "CRITICAL":
		level = CRITICAL
		break
	case "DEBUG":
		level = DEBUG
		break
	case "ERROR":
		level = ERROR
		break
	case "NOTICE":
		level = NOTICE
		break
	case "INFO":
		level = INFO
		break
	default:
		level = INFO
	}
}

//String log level to string
func (level LogLevel)String()string{
	switch level {
	case EMERGENCY:
		return "EMERGENCY"
	case ALERT:
		return "ALERT"
	case CRITICAL:
		return "CRITICAL"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	default:
		return "INFO"
	}
}

const (
	EMERGENCY = iota //
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
	CONSOLE=1
	FILE=2
	SYSLOG=3
	MAXSIZE=1024
)

//New new a binglog
func New(conf *config.IniConfig)*BingLog{
	var level LogLevel = INFO
	level.Set(conf.Get("loglevel"))
	bl := BingLog{}
	bl.Level = level
	var mode LogMode = CONSOLE
	mode.Set(conf.Get("logmode"))
	bl.Mode = mode
	bl.Path = conf.Get("logpath")
	bl.Daily = false
	bl.APPName = conf.Get("appname")
	bl.pool = make(chan LogMessage, MAXSIZE)
	return &bl
}

//WriteLog write log
func (bl *BingLog)WriteLog(level LogLevel, logFormat string, a ...interface{}){
	if len(bl.pool) == MAXSIZE {
		//pool is full, drop the older one
		<- bl.pool
	}
	bl.pool <- LogMessage{
		Message: fmt.Sprintf(logFormat, a),
		Level: level,
	}
}

//write
func (bl *BingLog)write(log LogMessage){
	if log.Level > bl.Level {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	fullLog := fmt.Sprintf("%s [%s] [%s] %s", now, bl.APPName, log.Level.String(), log.Message)
	switch bl.Mode {
	case CONSOLE:
		fmt.Println(fullLog)
		break
	case FILE:
		logFile := bl.APPName
		if bl.Path != "" {
			if _, err := os.Stat(bl.Path); err != nil || os.IsNotExist(err){
				os.MkdirAll(bl.Path, os.ModePerm)
			}
			if []byte(bl.Path)[len(bl.Path) -1] == '/' {
				logFile = fmt.Sprintf("%s%s", bl.Path, logFile)
			}
			logFile = fmt.Sprintf("%s/%s", bl.Path, logFile)
		}
		if bl.Daily {
			logFile = fmt.Sprintf("%s_%s", logFile, time.Now().Format("20060102"))
		}
		logFile = fmt.Sprintf("%s.log",logFile)
		var file *os.File
		if _, err := os.Stat(logFile); err != nil || os.IsNotExist(err){
			file, err = os.Create(logFile)
			if err != nil {
				panic(err)
			}
		} else {
			file, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE, os.ModeAppend)
			if err != nil {
				panic(err)
			}
		}
		writer := bufio.NewWriter(file)
		fmt.Fprintln(writer, fullLog)
		writer.Flush()
		defer file.Close()
		break
	case SYSLOG:
		break
	default:
		fmt.Println(fullLog)
	}
}

//thread
func (bl *BingLog)thread(){
	for {
		select {
		case log := <- bl.pool:
			//recivie a log
			bl.write(log)
		}
	}
}

//Start start the log thread
func (bl *BingLog)Start(){
	go bl.thread()
}
//Info write info log
func (bl *BingLog)Info(logFormat string, a ...interface{}){
	bl.WriteLog(INFO, logFormat, a)
}

//Notice write info log
func (bl *BingLog)Notice(logFormat string, a ...interface{}){
	bl.WriteLog(NOTICE, logFormat, a)
}

//Warning write info log
func (bl *BingLog)Warning(logFormat string, a ...interface{}){
	bl.WriteLog(WARNING, logFormat, a)
}

//Error write info log
func (bl *BingLog)Error(logFormat string, a ...interface{}){
	bl.WriteLog(ERROR, logFormat, a)
}

//Critical write info log
func (bl *BingLog)Critical(logFormat string, a ...interface{}){
	bl.WriteLog(CRITICAL, logFormat, a)
}

//Alert write info log
func (bl *BingLog)Alert(logFormat string, a ...interface{}){
	bl.WriteLog(ALERT, logFormat, a)
}

//Emergency Emergency info log
func (bl *BingLog)Emergency(logFormat string, a ...interface{}){
	bl.WriteLog(EMERGENCY, logFormat, a)
}

//Debug write info log
func (bl *BingLog)Debug(logFormat string, a ...interface{}){
	bl.WriteLog(DEBUG, logFormat, a)
}
