package handlers

import (
	"bytes"
	"net/http"

	"github.com/sh2nk/9stream/view"

	uuid "github.com/satori/go.uuid"
)

//Stream - обработка страницы с видеотрансляцией
func Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=informs")

	name := uuid.Must(uuid.NewV4(), nil).String()
	stream := make(chan *bytes.Buffer)

	func() {
		view.Pool.Lock()
		defer view.Pool.Unlock()
		view.Pool.Streams[name] = stream
	}()
	defer func() {
		view.Pool.Lock()
		defer view.Pool.Unlock()
		delete(view.Pool.Streams, name)
	}()

	for buf := range stream {
		w.Write(buf.Bytes())
	}
}
