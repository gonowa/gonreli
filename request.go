//+build js,wasm

package gonreli

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"syscall/js"
)

type request struct {
	js.Value
	headers http.Header
}

func newRequest(v js.Value) *http.Request {
	r := &request{
		Value:   v,
		headers: http.Header{},
	}
	pr, bodyWritter := io.Pipe()

	br := &bodyReader{r: v, pr: pr, writer: bodyWritter}
	reader := io.MultiReader(r.reqHeaderReader(), br)

	r.Call("setEncoding", "utf8")

	r.Call("on", "end", js.NewCallback(func(args []js.Value) {
		//to guard this
		bodyWritter.Close()

	}))
	hr, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		panic(err)
	}
	return hr
}

func (req *request) reqHeaderReader() io.Reader {
	hdrs := req.Get("rawHeaders")
	for i := 0; i < hdrs.Length()-1; i += 2 {
		req.headers.Set(hdrs.Index(i).String(), hdrs.Index(i+1).String())
	}
	m := req.Get("method").String()
	u := req.Get("url").String()
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	httpVersion := req.Get("httpVersion")
	builder.WriteString(fmt.Sprintf("%s %s %s/%s\n", m, url.Path, "HTTP", httpVersion))
	for i := range req.headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", i, req.headers[i][0]))
	}
	builder.WriteString("\r\n")

	return strings.NewReader(builder.String())
}
func (req *request) Headers() http.Header {
	return req.headers
}

type streamWritable struct {
	js.Value
	io.Writer
}

func newStreamWritable(writer io.Writer) *streamWritable {

	sw := streamWritable{Writer: writer}

	stream := jsRequire.Invoke("stream")
	s := stream.Get("Writable").New()
	sw.Value = s
	s.Set("_write", js.NewCallback(func(args []js.Value) {
		chunk := args[0]
		done := args[2]
		var buf = make([]byte, chunk.Length())
		for i := 0; i < chunk.Length(); i++ {
			buf[i] = byte(chunk.Index(i).Int())
		}
		_, err := sw.Write(buf)
		if err != nil {
			done.Invoke(err.Error())
		}
		done.Invoke()
	}))
	return &sw
}

type bodyReader struct {
	r      js.Value
	writer *io.PipeWriter
	once   sync.Once
	pr     io.Reader
}

func (r *bodyReader) Read(p []byte) (n int, err error) {
	r.once.Do(func() {
		streamWriter := newStreamWritable(r.writer).Value
		r.r.Call("pipe", streamWriter)
	})
	return r.pr.Read(p)
}
