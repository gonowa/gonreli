package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type Test struct {
	Name     string
	Request  *http.Request
	Handler  http.HandlerFunc
	Response http.Response
}

var echoHandler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	for name, value := range request.Header {
		writer.Header().Set(name, value[0])
	}
	if body != nil {
		_, err = io.Copy(writer, bytes.NewReader(body))
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
	}
}
var test = []Test{
	{
		Name: "simple Get",
		Request: &http.Request{
			Header: http.Header{
				"a": []string{"a"},
			},
		},
	},
}

func TestEnd2End(t *testing.T) {

}
