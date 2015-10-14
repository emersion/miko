package client

import (
	"encoding/base64"
	"github.com/gopherjs/gopherjs/js"
)

type ExportWriter struct {
	output []byte
}

func (w *ExportWriter) Write(data []byte) (n int, err error) {
	w.output = append(w.output, data...)
	return len(data), nil
}

func (w *ExportWriter) flush() []byte {
	o := w.output
	w.output = make([]byte, 0)
	return o
}

func (w *ExportWriter) Export() {
	bin := w.flush()
	str := base64.StdEncoding.EncodeToString(bin)
	js.Global.Get("window").Call("open", "data:application/octet-stream;base64,"+str)
}

type ExportReader struct {
	File *js.Object
	input []byte
}

func (r *ExportReader) Read(data []byte) (n int, err error) {
	if len(r.input) == 0 {
		readCh := make(chan []byte, 1)

		fileReader := js.Global.Get("FileReader").New()
		fileReader.Call("addEventListener", "load", func() {
			readCh <- []byte(fileReader.Get("result").String())
		})
		fileReader.Call("readAsBinaryString", r.File)

		r.input = <-readCh
	}

	i := 0
	for i = range data {
		data[i] = r.input[i]
	}

	r.input = r.input[i+1:]

	return i+1, nil
}
