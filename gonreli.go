//+build js,wasm

package gonreli

import (
	"errors"
	"net/http"
	"syscall/js"
)

var (
	jsRequire    = js.Global().Get("require")
	jsSocketType = jsRequire.Invoke("net").Get("Socket")
	jsBufferType = js.Global().Get("Buffer")
	jsProcess    = js.Global().Get("process")
	undefined    = js.Undefined()
)

func Wrap(h http.Handler) js.Callback {
	return js.NewCallback(func(args []js.Value) {
		req, resp := newRequest(args[0]), newResponse(args[1])
		go func() {
			var err error
			defer func() {
				r := recover()
				if r != nil {
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("Unknown error")
					}
					http.Error(resp, err.Error(), http.StatusInternalServerError)
				}
			}()

			h.ServeHTTP(resp, req)
			if resp.hijacked {
				//resp.Get("connection").Call("end")
			} else {
				resp.Call("end")
			}
		}()
	})
}
