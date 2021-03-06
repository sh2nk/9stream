package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sh2nk/9stream/handlers"
	"github.com/sh2nk/9stream/view"
)

var port int

func init() {
	flag.StringVar(&view.Device, "d", "video0", "Select V4l2 device.")
	flag.IntVar(&port, "p", 8080, "Select server port.")
	flag.BoolVar(&view.WhiteBalance, "w", false, "Enable auto white balance.")
}

func main() {
	var err error

	flag.Parse()
	handlers.Ready = false

	if err = view.CameraInit(view.Device, view.WhiteBalance); err != nil {
		log.Fatalln(err)
	}
	defer view.Camera.Close()

	port := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/stream", handlers.Stream)
	http.HandleFunc("/setup", handlers.Setup)

	log.Printf("Started server on %s port using /dev/%s device", port, view.Device)
	panic(http.ListenAndServe(port, nil))
}
