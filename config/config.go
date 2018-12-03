package config

import (
	"strings"
	"io"
	"bufio"
	"os"

)

//IniConfig config used ini format
type IniConfig map[string]string



//Load load config from file
func (i IniConfig)Load(path string){
	i = make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		words := strings.Split(line, "=")
		if len(words) > 2 && words[0] != "" {
			i.Set(words[0], words[1])
		}
	}
}

//Set set a config
func (i IniConfig)Set(key, val string){
	i[key] = val
}

//Get get a config
func (i IniConfig)Get(key string)string{
	for k, v := range i {
		if k == key {
			return v
		}
	}
	return ""
}