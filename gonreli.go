//+build js

package gonreli

import (
	"net/http"
	"syscall/js"
)

var (
	jsRequire    = js.Global().Get("require")
	jsSocketType = jsRequire.Invoke("net").Get("Socket")
	jsBufferType = js.Global().Get("Buffer")
	jsProcess    = js.Global().Get("process")
)

func Wrap(h http.Handler) js.Callback {
	return js.NewCallback(func(args []js.Value) {
		go func() {
			req := newRequest(args[0])
			req, resp := req, newResponse(args[1])
			h.ServeHTTP(resp, req)
			if resp.hijacked {
				//resp.Get("connection").Call("end")
			} else {
				resp.Call("end")
			}
		}()
	})
}
