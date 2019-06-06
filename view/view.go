package view

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/blackjack/webcam"
)

//Camera - глобальный объект камеры.
var Camera *webcam.Webcam

//Device - устройство v4l.
var Device string

//WhiteBalance - используем ли автоматический баланс белого.
var WhiteBalance bool

//Timeout - таймаут камеры
type Timeout *webcam.Timeout

//FrameSizes - обертка для []webcam.FrameSize
type FrameSizes []webcam.FrameSize

//CameraInit - инициализация камеры
func CameraInit(d string, wb bool) (err error) {
	if Camera, err = webcam.Open(fmt.Sprintf("/dev/%s", d)); err != nil {
		return
	}
	if wb {
		if err = Camera.SetAutoWhiteBalance(true); err != nil {
			return
		}
	}

	return
}

//StreamPool - структура пула стрима
type StreamPool struct {
	sync.RWMutex
	Streams map[string]chan *bytes.Buffer
}

//Pool - пул клиентов стрима
var Pool = &StreamPool{
	Streams: make(map[string]chan *bytes.Buffer, 12),
}

//StopStream - отановить ли стрим
var StopStream bool

// ReadStream - читаем поток
func ReadStream(camera *webcam.Webcam, pool *StreamPool, jpeg bool) {
	var err error

	// Универсальная обработка ошибок
	defer func() {
		if rec := recover(); rec != nil {
			// паника!
			var ok bool
			if err, ok = rec.(error); !ok {
				log.Printf("ReadStream panic: %#v", rec)
			}
		}

		if err != nil {
			log.Printf("ReadStream error: %#v", err)
		}
	}()

	err = camera.StartStreaming()
	if err != nil {
		return
	}
	defer camera.StopStreaming()

	for {
		if StopStream {
			return
		}

		//Таймаут 5 секунд
		if err = camera.WaitForFrame(5); err != nil {
			if _, ok := err.(*webcam.Timeout); ok {
				continue
			}
			return
		}

		frame, err := camera.ReadFrame()
		if err != nil {
			return
		}

		if len(frame) == 0 {
			continue
		}

		buf := new(bytes.Buffer)
		buf.Grow(len(frame) + 100)
		if jpeg {
			buf.Write([]byte(fmt.Sprintf("Content-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(frame))))
			buf.Write(frame)
			buf.Write([]byte("\r\n--informs\r\n"))
		} else {
			buf.Write([]byte(fmt.Sprintf("Content-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(frame))))
			buf.Write(frame)
			buf.Write([]byte("\r\n--informs\r\n"))
		}

		func() {
			pool.RLock()
			defer pool.RUnlock()

			for name := range pool.Streams {
				pool.Streams[name] <- buf
			}
		}()
	}
}
