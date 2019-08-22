package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"io"
	"strconv"
	"strings"
)
var palette = []color.Color{color.RGBA{0xf3, 0x43, 0x11, 0xff},
	color.RGBA{0x44, 0x08, 0x11, 0xff}, color.RGBA{0x11, 0x55, 0xaa, 0xff}}
var mu sync.Mutex
var count int

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", lissajousHandler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/detail", detailHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func lissajous(out io.Writer, cycles float64) {
	const (
//		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolutions
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // ignore encoding error
}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.RawQuery
	cycle := 5
	var queries []string
	if rawQuery != "" {
		queries = strings.Split(rawQuery, "&")
	}
	for _, q := range queries {
		if strings.HasPrefix(q, "cycles") {
			item := strings.Split(q, "=")
			var err error
			cycle, err = strconv.Atoi(item[1])
			if err != nil {
				cycle = 5
			}
		}
	}
	lissajous(w, float64(cycle))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}
