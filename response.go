//+build js

package gonreli

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"syscall/js"
	"time"
)

//Response implements http.ResponseWriter
type Response struct {
	hijacked bool

	reqReader io.Reader
	headers   http.Header
	js.Value
	body bytes.Buffer

	wroteHeader bool
}

func newResponse(v js.Value) *Response {

	pr, bodyWritter := io.Pipe()

	br := &bodyReader{r: v, pr: pr, writer: bodyWritter}

	return &Response{Value: v, headers: http.Header{}, reqReader: br}
}

func (w *Response) Header() http.Header {
	return w.headers
}

func (w *Response) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(200)
	}
	if len(data) == 0 {
		return 0, nil
	}

	ta := js.TypedArrayOf(data)
	defer ta.Release()
	//todo check error
	w.Call("write", jsBufferType.New(ta))
	return (len(data)), nil
}

func (w *Response) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.Set("statusCode", statusCode)
		w.writeHeader()
	}
}

//todo improve this is fragile
func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	conn := newNodeConn(w.Value, w.reqReader)
	brw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return conn, brw, nil

}

func (w *Response) writeHeader() {
	for i := range w.headers {
		w.Call("setHeader", i, w.headers[i][0])
	}
	w.wroteHeader = true
}

type nodeConn struct {
	js.Value
	r io.Reader
}

func newNodeConn(v js.Value, reader io.Reader) *nodeConn {
	if !v.InstanceOf(jsSocketType) {
		v = v.Get("socket")
	}
	var nc nodeConn
	nc.r = reader

	nc.Value = v
	nc.Call("on", "timeout", js.NewCallback(func(args []js.Value) {
		log.Println("timeout received")
		nc.Call("end")
	}))
	return &nc
}

func (nc *nodeConn) Read(b []byte) (n int, err error) {
	return nc.r.Read(b)
}

func (nc *nodeConn) Write(b []byte) (n int, err error) {
	ta := js.TypedArrayOf(b)
	nc.Call("write", jsBufferType.Call("from", ta))
	ta.Release()
	return len(b), nil
}

func (nc *nodeConn) Close() error {
	nc.Call("end")
	return nil
}

func (nc *nodeConn) LocalAddr() net.Addr {
	addr := nc.Get("localAddress").String()
	port := nc.Get("localPort").Int()
	ip := net.ParseIP(addr)
	return &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
}

func (nc *nodeConn) RemoteAddr() net.Addr {
	addr := nc.Get("remoteAddress").String()
	port := nc.Get("remotePort").Int()
	ip := net.ParseIP(addr)
	return &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
}

func (nc *nodeConn) SetDeadline(t time.Time) error {
	//socket.setTimeout(3000);
	var timeout time.Duration
	if t.IsZero() {
		timeout = 0
	} else {
		now := time.Now()
		timeout = t.Sub(now)
	}

	nc.Call("setTimeout", timeout/time.Millisecond)
	return nil
}

func (nc *nodeConn) SetReadDeadline(t time.Time) error {
	var timeout time.Duration
	if t.IsZero() {
		timeout = 0
	} else {
		now := time.Now()
		timeout = t.Sub(now)
	}

	nc.Call("setTimeout", timeout/time.Millisecond)
	return nil
}

func (nc *nodeConn) SetWriteDeadline(t time.Time) error {
	var timeout time.Duration
	if t.IsZero() {
		timeout = 0
	} else {
		now := time.Now()
		timeout = t.Sub(now)
	}

	nc.Call("setTimeout", timeout/time.Millisecond)
	return nil
}
