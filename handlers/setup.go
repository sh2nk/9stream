package handlers

import (
	"9stream/view"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blackjack/webcam"
)

//Ready - проверка готовности
var Ready bool

//Setup - обработчик страницы настройки
func Setup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	formatDesc := view.Camera.GetSupportedFormats()
	var formats []webcam.PixelFormat
	for f := range formatDesc {
		formats = append(formats, f)
	}

	formatKey := r.URL.Query().Get("format")
	sizeKey := r.URL.Query().Get("size")

	if formatKey == "" && sizeKey == "" {
		fmt.Fprintf(w, "<div><h1>Available formats:</h1></div>\n")
		for i, value := range formats {
			fmt.Fprintf(w, "<div>[%d] <a href='/setup?format=%d'>%s</a></div>\n", i+1, i, formatDesc[value])
		}
	}

	if formatKey != "" && sizeKey == "" {
		sel, err := strconv.ParseInt(formatKey, 10, 32)
		if err != nil {
			log.Println("Falied on size")
			log.Fatal(err)
		}
		format := formats[sel]

		fmt.Fprintf(w, "<div><h1>Supported frame sizes for format %s</h></div>\n", formatDesc[format])
		frames := view.FrameSizes(view.Camera.GetSupportedFrameSizes(format))
		sort.Slice(frames, func(i, j int) bool {
			ls := frames[i].MaxWidth * frames[i].MaxHeight
			rs := frames[j].MaxWidth * frames[j].MaxHeight
			return ls < rs
		})

		for i, value := range frames {
			fmt.Fprintf(w, "<div>[%d] <a href='/setup?format=%d&size=%d'>%s</a></div>\n", i+1, sel, i, value.GetString())
		}
	}

	if formatKey != "" && sizeKey != "" {
		selFormat, err := strconv.ParseInt(formatKey, 10, 32)
		if err != nil {
			log.Println("Falied on final stage parsing fkey")
			return
		}
		selSize, err := strconv.ParseInt(sizeKey, 10, 32)
		if err != nil {
			log.Println("Falied on final stage parsing skey")
			return
		}

		format := formats[selFormat]
		frames := view.FrameSizes(view.Camera.GetSupportedFrameSizes(format))
		size := frames[selSize]

		view.StopStream = true
		if Ready {
			time.Sleep(time.Second)
		}
		view.StopStream = false

		view.Camera.Close()
		if err = view.CameraInit(view.Device, view.WhiteBalance); err != nil {
			log.Println("Falied on final stage setting cam format")
			fmt.Fprint(w, err)
			Ready = false
			return
		}

		_, _, _, err = view.Camera.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
		if err != nil {
			log.Println("Falied on final stage setting cam format")
			fmt.Fprint(w, err)
			Ready = false
			return
		}

		go view.ReadStream(view.Camera, view.Pool, strings.Contains(strings.ToLower(formatDesc[formats[selFormat]]), "jpeg"))
		Ready = true
		http.Redirect(w, r, "/stream", http.StatusSeeOther)
	}
}
