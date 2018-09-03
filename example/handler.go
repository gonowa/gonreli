//+build js,wasm

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gonowa/gonobridge"
	"github.com/gonowa/gonreli"
	"github.com/rs/zerolog/log"
)

func main() {

	r := gin.New()
	r.GET("/go/hello", helloHandler)
	r.GET("/go/spacetravel", func(context *gin.Context) {
		context.Writer.Header().Set("Content-Type", "image/gif")
		err := spaceTravel(context.Writer)
		if err != nil {
			context.Error(err)
		}
	})
	r.GET("/go/hijacked", hijack)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, c.Request.URL)
	})
	gonobridge.Emit("handler", gonreli.Wrap(r))
	gonobridge.Wait()
}

var helloHandler gin.HandlerFunc = func(context *gin.Context) {
	//do something otherwise call overhead will dominate
	var s = 0
	for i := 0; i < 1000000; i++ {
		var x = i
		for x != 0 {
			x &= x - 1
			s++
		}
	}
	context.Writer.Write([]byte("Hello World from go!"))
}

var hijack gin.HandlerFunc = func(context *gin.Context) {
	conn, brw, err := context.Writer.Hijack()
	if err != nil {
		context.Error(err)
	}
	_ = conn

	var httpreply = `HTTP/1.1 200 OK
Date: Wed, 26 Nov 2018 03:37:57 GMT
Content-Length: %d
Content-Type: text/plain; charset=utf-8

%s`
	var text = "hello hijacked"

	fmt.Fprintf(brw, httpreply, len(text), text)
	brw.Flush()
	log.Print(conn.LocalAddr())
	log.Print(conn.RemoteAddr())
}
