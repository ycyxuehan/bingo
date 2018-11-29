# bingo

bingo is a little web service framework and used gorilla mux router.

- example
```go
package main

import (
	"github.com/ycyxuehan/bingo"
)

//
type TestController struct {
	bingo.Controller
}

func (t *TestController)Get(){
	t.Ctx.ServeString("this is a test")
}

func main(){
	router := bingp.NewRouter()
	router.Add("/", &TestController{})
	b := bingo.New("bingo.conf")
	b.Run(router)
}
```

config file

```conf
    host=0.0.0.0
    port=8080
```