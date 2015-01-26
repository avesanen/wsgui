package main

import (
	//"encoding/base64"
	"github.com/googollee/go-socket.io"
	"github.com/zenazn/goji"
	//"image"
	"bytes"
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"os"
)

type eventMouseDown struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type eventDrawCanvas struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	R uint32 `json:"r"`
	G uint32 `json:"g"`
	B uint32 `json:"b"`
	A uint32 `json:"a"`
}

type eventBlitPng struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"w"`
	Height int    `json:"h"`
	Data   []byte `json:"data"`
}

func main() {
	infile, err := os.Open("./test.png")
	if err != nil {
		panic(err.Error())
	}
	defer infile.Close()
	//reader := base64.NewDecoder(base64.StdEncoding, infile)

	m, err := png.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("New Connection!")
		so.Join("chat")
		so.On("mousedown", func(msg string) {
			var e eventMouseDown
			err := json.Unmarshal([]byte(msg), &e)
			if err != nil {
				log.Println("Can't unmarshal:", err.Error())
				return
			}

			buf := new(bytes.Buffer)
			err = png.Encode(buf, m)
			if err != nil {
				log.Println("Can't encode:", err.Error())
				return
			}
			encoded := buf.Bytes()

			blitEvent := &eventBlitPng{
				X:      e.X - (bounds.Max.X / 2),
				Y:      e.Y - (bounds.Max.Y / 2),
				Width:  bounds.Max.X,
				Height: bounds.Max.Y,
				Data:   encoded,
			}

			bytes, err := json.Marshal(blitEvent)
			if err != nil {
				log.Println("Can't marshal:", err.Error())
				return
			}
			so.Emit("blit", string(bytes))
			so.BroadcastTo("chat", "blit", string(bytes))
		})
	})

	goji.Get("/socket.io/", server)
	goji.Get("/*", http.FileServer(http.Dir("./web")))
	goji.Serve()
}
