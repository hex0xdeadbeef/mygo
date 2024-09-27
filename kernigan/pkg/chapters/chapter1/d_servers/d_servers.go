package d_servers

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func First_Server() {
	http.HandleFunc("/", getRequestPath)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var mu sync.Mutex
var count int

func Second_Server() {
	http.HandleFunc("/", getRequestPath)
	http.HandleFunc("/req", getRequestInformation)
	http.HandleFunc("/lissajou", getLissajouGif)
	http.HandleFunc("/count", getRequestCount)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func getRequestPath(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()

	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

func getRequestCount(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	fmt.Fprintf(writer, "Count = %d\n", count)
	mu.Unlock()
}

func getRequestInformation(writer http.ResponseWriter, request *http.Request) {

	mu.Lock()
	count++
	mu.Unlock()

	fmt.Fprintf(writer, "%s %s %s\n", request.Method, request.URL, request.Proto)
	for k, v := range request.Header {
		fmt.Fprintf(writer, "Header[%s] = %s\n", k, v)
	}
	fmt.Fprintf(writer, "Remote Address = %s\n", request.RemoteAddr)
	if err := request.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range request.Form {
		fmt.Fprintf(writer, "Form[%s] = %s\n", k, v)
	}
}

func getLissajouGif(writer http.ResponseWriter, request *http.Request) {

	mu.Lock()
	count++
	mu.Unlock()

	cyclesParam := request.URL.Query().Get("cycles")

	if cyclesParam != "" {
		if cyclesValue, err := strconv.ParseFloat(cyclesParam, 64); err != nil {
			fmt.Fprintf(writer, err.Error())
		} else {
			lissajous(writer, int(cyclesValue))
		}
	} else {
		lissajous(writer)
	}
}

func lissajous(out io.Writer, numberOfCycles ...int) {

	var (
		cycles  float64
		palette = []color.Color{
			color.Black,
			color.RGBA{0, 255, 0, 255},
			color.RGBA{255, 255, 0, 255},
			color.White,
		}
	)

	start := time.Now() // Start timer

	if len(numberOfCycles) != 0 { //количество полных оборотов осциллятора x
		cycles = float64(numberOfCycles[0])
	} else {
		cycles = 10.0
	}
	const (
		res     = 0.001 // угловое разрешение
		size    = 500   // размер холста [-size...+size]
		nframes = 64    // количество кадров анимации
		delay   = 1     // задержка между кадрами в единицах 10 мс
	)

	freq := rand.Float64() * 3.0 // относительная частота осциллятора y
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // разность фаз
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			index := uint8(rand.Intn(4))
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	fmt.Fprint(os.Stdout, "%g\n", time.Since(start).Milliseconds())
	gif.EncodeAll(out, &anim)
}
