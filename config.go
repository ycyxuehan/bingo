package bingo

import (
	"strings"
	"io"
	"bufio"
	"os"

)

//IniConfig config used ini format
type IniConfig map[string]string

func NewConfig()*IniConfig{
	conf := IniConfig{}
	conf = make(map[string]string)
	return &conf
}

var BingConf *IniConfig

func init(){
	BingConf = NewConfig()
	BingConf.Load("conf/app.conf")
}

//Load load config from file
func (i *IniConfig)Load(path string)error{
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF{
			break
		}
		words := strings.Split(line, "=")
		if len(words) >= 2 && words[0] != "" {
			i.Set(words[0], strings.Trim(words[1], "\n"))
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

//Set set a config
func (i *IniConfig)Set(key, val string){
	(*i)[key] = val
}

//Get get a config
func (i *IniConfig)Get(key string)string{
	for k, v := range *i {
		if k == key {
			return v
		}
	}
	return ""
}