package chapter1

import (
	"bufio"
	"encoding/json"
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
	"strings"
	"sync"
	"time"
)

// 1.1
func HelloWorld() {
	fmt.Println("Hello world!")
}

// 1.2
func Echo_SimpeFor() {
	var separator, resultString string
	separator = " "
	for i := 0; i < len(os.Args); i++ {
		resultString += os.Args[i] + separator
	}
	fmt.Println(resultString)
}

func For_While() {
	number := 0
	for number < 10 {
		number++
	}
	fmt.Println(number)
}

func For_Infinite() {
	number := 2
	for {
		if number < 10e10 {
			number *= 2
		} else {
			break
		}

	}
	fmt.Println(number)
}

func Echo_RangeFor() {
	str := ""
	for index, cmdArgument := range os.Args {
		fmt.Print(index, " ")
		str += cmdArgument + " "

	}
	fmt.Println(str)
}

func Echo_stingsJoin() {
	fmt.Println(strings.Join(os.Args[0:], " "))
}

func Echo_StraightforwardPrint() {
	fmt.Println(os.Args[0:])
}

// 1.3

func Dup_First() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func Dup_Second(fileNames ...string) {
	counts := make(map[string]int)
	filesNames := fileNames

	if len(filesNames) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, fileName := range filesNames {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Dup: %v\n", err)
				continue
			} else {
				countLines(file, counts)
				file.Close()
			}
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func Dup_Third() {
	counts := make(map[string]int)

	for _, filename := range os.Args[1:] {
		dataFromFile, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprint(os.Stderr, "Dup_Third:\n", err)
			continue
		}

		for _, line := range strings.Split(string(dataFromFile), "\n") {
			counts[line]++
		}

	}

	for line, number := range counts {
		if number > 1 {
			fmt.Printf("%d : %s", number, line)
		}
	}
}

// 1.3 Homework
func Dup_Second_Modified(filesNames ...string) {

	if len(filesNames) == 0 {
		counts := make(map[string]int)

		countLines(os.Stdin, counts)

		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	} else {
		fileCounts := make(map[int]map[string]int) // Create map of maps

		// Fill a map for each file
		for i, fileName := range filesNames {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Dup: %v\n", err)
				continue
			} else {
				counts := make(map[string]int) // Creates an empty map
				countLines(file, counts)       // Pass a map and file into countLines() to be filled
				fileCounts[i] = counts         // Add a filled map into map of maps
				file.Close()                   // Close a file
			}
		}

		unionedMap := mapUnion(fileCounts) // Create an union map

		// Go through the unioned map, print count, string, files names strings are contained in

		for key, value := range unionedMap {
			fmt.Printf("|%d|\t|%s|\t", value, key)
			for i, curMap := range fileCounts {
				if _, ok := curMap[key]; ok {
					fmt.Printf("|%s|\t", filesNames[i])
				}
			}
			fmt.Println()
		}

	}

}

func mapUnion(mapOfMaps map[int]map[string]int) map[string]int {
	unionedMap := make(map[string]int)
	for _, curMap := range mapOfMaps {
		for key, value := range curMap {
			if _, ok := unionedMap[key]; !ok {
				unionedMap[key] = value
			} else {
				unionedMap[key] += value
			}

		}
	}

	fmt.Println()
	return unionedMap
}

func countLines(file *os.File, counts map[string]int) {
	input := bufio.NewScanner(file)
	for input.Scan() {
		counts[input.Text()]++
	}
}

// Определение констант и переменных палитры и индекса цвета

var palette = []color.Color{color.Black, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 255, 0, 255}, color.White}

const (
	whiteIndex = 0 // Первый цвет палитры
	blackIndex = 1 // Второй цвет палитры
)

// Определение функции Lissajous

func Lissajous(out io.Writer, numberOfCycles ...int) {
	start := time.Now() // Start timer
	var cycles float64
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
			var index uint8
			index = uint8(rand.Intn(4))
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	fmt.Fprint(os.Stdout, "%g\n", time.Since(start).Milliseconds())
	gif.EncodeAll(out, &anim)
}

func Fetch(urls ...string) {
	for _, url := range urls {
		if !isHTTPPrefixed(url) {
			addHTTPPrefix(&url)
		}
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		numberOfWrittenBytes, err := io.Copy(os.Stdout, response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bytesArray := make([]byte, 0)
		bytesArray = append(bytesArray, byte(numberOfWrittenBytes))
		bytesArray = append(bytesArray, []byte(" "+response.Status+"\n")...)

		_, err = os.Stdout.Write(bytesArray)
		if err != nil {
			fmt.Println()
			os.Exit(1)
		}
	}
}

func isHTTPPrefixed(url string) bool {
	return strings.HasPrefix(url, "http://")
}

func addHTTPPrefix(url *string) {
	*url = "http://" + *url
}

func Fetch_Concurrently(filename string, urls ...string) {
	start := time.Now() // Start timer
	urlChannel := make(chan string)

	file, err := os.Create(filename)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
		os.Exit(1)
	}
	defer file.Close()

	for _, url := range urls {
		if !isHTTPPrefixed(url) {
			addHTTPPrefix(&url)
		}
		go fetch_Uploading_Channel(url, urlChannel)
	}

	for range urls {
		file.WriteString(fmt.Sprintf("%s\n", <-urlChannel))
	}
	file.WriteString(fmt.Sprintf("%.2fs elapsed\n", time.Since(start).Seconds()))
}

func fetch_Uploading_Channel(url string, channel chan<- string) {
	start := time.Now()

	response, err := http.Get(url)
	if err != nil {
		channel <- fmt.Sprintf(err.Error())
		return
	}

	numberOfWrittenBytes, err := io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		channel <- fmt.Sprintf("While reading %s: %s", url, err.Error())
	}

	secs := time.Since(start).Seconds()
	channel <- fmt.Sprintf("%.2fs %d %s", secs, numberOfWrittenBytes, url)
}

func ParsingWebsitesJson(fileName string) []string {
	type SiteData struct {
		Position int     `json:"position"`
		Domain   string  `json:"domain"`
		Count    int32   `json:"count"`
		Etv      float64 `json:"etv"`
	}

	file, err := os.Open("ranked_domains.json")
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
	}

	var jsonData []SiteData

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		os.Stdout.Write([]byte(err.Error() + " "))
	}

	var urls []string
	for i := 0; i < 100; i++ {
		urls = append(urls, jsonData[i].Domain)
	}

	return urls
}

func First_Server() {
	http.HandleFunc("/", firstHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func firstHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

var mu sync.Mutex
var count int

func Second_Server() {
	//http.HandleFunc("/", secondHandler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/req", thirdHandler)
	http.HandleFunc("/lissajou", lissajouHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func secondHandler(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

func counter(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	fmt.Fprintf(writer, "Count = %d\n", count)
	mu.Unlock()
}

// fmt.Fprintf(writer, "", )
func thirdHandler(writer http.ResponseWriter, request *http.Request) {
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

func lissajouHandler(writer http.ResponseWriter, request *http.Request) {
	cyclesParam := request.URL.Query().Get("cycles")

	if cyclesParam != "" {
		if cyclesValue, err := strconv.ParseFloat(cyclesParam, 64); err != nil {
			fmt.Fprintf(writer, err.Error())
		} else {
			Lissajous(writer, int(cyclesValue))
		}
	} else {
		Lissajous(writer)
	}
}
