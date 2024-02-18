package fractals

// http://localhost:8000/create?type=mandelbrot&smoothing=0
import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	width, height = 919, 1903
	iterations    = 200
	contrast      = 15
)

var (
	img                    *image.RGBA
	parametersToBeApplied  []string
	zoom                           = 1.0
	xmin, ymin, xmax, ymax float64 = -2, -2, 2, 2
)

func reInit() {
	img = image.NewRGBA(image.Rect(0, 0, width, height))
	parametersToBeApplied = make([]string, 0)
}

func Server() {

	http.HandleFunc("/create", fractalHandler)
	log.Fatalf("Error occured: %v", http.ListenAndServe("localhost:8000", nil))
}

func fractalHandler(writer http.ResponseWriter, request *http.Request) {
	parameters := request.URL.Query()
	if len(parameters) != 0 {
		reInit()
		for key, arrayOfValues := range parameters {
			switch key {
			case "type":
				value := arrayOfValues[0]
				switch value {
				case "":
					fmt.Fprint(os.Stderr, "Error occured: no type parameter\n")
					os.Exit(1)
				case "mandelbrot":
					parametersToBeApplied = append(parametersToBeApplied, "mandelbrot")
				case "newton":
					parametersToBeApplied = append(parametersToBeApplied, "newton")

				}

			case "smoothing":
				value := arrayOfValues[0]
				switch value {
				case "":
					fmt.Fprint(os.Stderr, "Error occured: no smoothing parameter\n")
					os.Exit(1)
				case "0":
					parametersToBeApplied = append(parametersToBeApplied, "0")
				case "1":
					parametersToBeApplied = append(parametersToBeApplied, "1")
				}

			case "zoom":
				value := arrayOfValues[0]
				switch value {
				case "":
					fmt.Fprint(os.Stderr, "Error occured: no smoothing parameter\n")
					os.Exit(1)
				default:
					if value, err := strconv.ParseFloat(value, 64); err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
					} else {
						zoom = value
					}
				}
			case "x":
				value := arrayOfValues[0]
				switch value {
				case "":
					fmt.Fprint(os.Stderr, "Error occured: no x parameter\n")
					os.Exit(1)
				default:
					if value, err := strconv.ParseFloat(value, 64); err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
					} else {
						xmin, xmax = -value, value
					}
				}
			case "y":
				value := arrayOfValues[0]
				switch value {
				case "":
					fmt.Fprint(os.Stderr, "Error occured: no y parameter\n")
					os.Exit(1)
				default:
					if value, err := strconv.ParseFloat(value, 64); err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
					} else {
						ymin, ymax = -value, value
					}
				}
			}
		}

	} else {
		fmt.Fprintf(os.Stderr, "Error: no parameters passed")
		os.Exit(1)
	}

	parallelMandelbrot4()
	rotate90Right(img)
	switch parametersToBeApplied[1] {
	case "0":
	case "1":
		smooth()
	}

	png.Encode(writer, img)

}

func createFractal() {
	timer := time.Now()
	for px := 0; px < width; px++ {
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin)*zoom + xmin
		for py := 0; py < height; py++ {
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin)*zoom + ymin
			z := complex(-x, y)
			switch parametersToBeApplied[0] {
			case "mandelbrot":
				img.Set(px, py, mandelbrot(z))
			case "newton":
				img.Set(px, py, newton(z))
			}
		}
	}

	log.Println(time.Since(timer))
}

type vector struct {
	px, py int
	colour color.RGBA
}

func parallelMandelbrot1() {
	timer := time.Now()

	colorChan := make(chan vector)
	var wg sync.WaitGroup

	for px := 0; px < width; px++ {
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin)*zoom + xmin
		for py := 0; py < height; py++ {
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin)*zoom + ymin
			z := complex(-x, y)
			wg.Add(1)
			go func(px, py int) {
				defer wg.Done()
				vector := vector{px: px, py: py, colour: mandelbrot(z)}
				colorChan <- vector
			}(px, py)
		}
	}

	go func() {
		wg.Wait()
		close(colorChan)
	}()

	for vector := range colorChan {
		img.Set(vector.px, vector.py, vector.colour)
	}
	log.Println(time.Since(timer))
}

func parallelMandelbrot2() {
	timer := time.Now()

	colorChan := make(chan vector, width*height)
	var wg sync.WaitGroup

	for px := 0; px < width; px++ {
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin)*zoom + xmin
		for py := 0; py < height; py++ {
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin)*zoom + ymin
			z := complex(-x, y)
			wg.Add(1)
			go func(px, py int) {
				defer wg.Done()
				vector := vector{px: px, py: py, colour: mandelbrot(z)}
				colorChan <- vector
			}(px, py)
		}
	}

	go func() {
		wg.Wait()
		close(colorChan)
	}()

	for vector := range colorChan {
		img.Set(vector.px, vector.py, vector.colour)
	}
	log.Println(time.Since(timer))
}

func parallelMandelbrot3() {
	timer := time.Now()

	sendColorChan := make(chan vector, width*height)

	var wg sync.WaitGroup

	for px := 0; px < width; px++ {
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin)*zoom + xmin
		for py := 0; py < height; py++ {
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin)*zoom + ymin
			z := complex(-x, y)
			wg.Add(1)
			go func(px, py int) {
				defer wg.Done()
				vector := vector{px: px, py: py, colour: mandelbrot(z)}
				sendColorChan <- vector
				// why not set img here?
			}(px, py)
		}
	}

	for v := range sendColorChan {
		wg.Add(1)
		go func(v vector) {
			img.Set(v.px, v.py, v.colour)
		}(v)
	}

	// sender
	go func() {
		wg.Wait()
		close(sendColorChan)
	}()

	log.Println(time.Since(timer))
}

func parallelMandelbrot4() {
	timer := time.Now()

	var wg sync.WaitGroup

	for px := 0; px < width; px++ {
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin)*zoom + xmin
		for py := 0; py < height; py++ {
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin)*zoom + ymin
			z := complex(-x, y)
			wg.Add(1)
			go func(px, py int) {
				//defer wg.Done()
				img.Set(px, py, mandelbrot(z))
				wg.Done()
			}(px, py)
		}
	}

	wg.Wait()
	log.Println(time.Since(timer))
}

func smooth() {

	for px := 0; px < height; px++ {
		for py := 0; py < width; py++ {
			pixelEnvironmentColors := getPixelEnvironmen(px, py)
			// fmt.Println(pixelEnvironmentColors)
			setAverageColor(px, py, pixelEnvironmentColors)
			// fmt.Println(px, py)
			pixelEnvironmentColors = make([]color.RGBA, 0)
			// fmt.Println(pixelEnvironmentColors)
		}
	}
}

func getPixelEnvironmen(px, py int) []color.RGBA {
	var pixels []color.RGBA
	for i := -1; i < 2; i++ {
		surroundingXPixel := px - 1*i
		for j := -1; j < 2; j++ {
			surroundingYPixel := py - 1*j
			if (surroundingXPixel >= 0 || surroundingXPixel < width) &&
				(surroundingYPixel >= 0 || surroundingYPixel < width) {
				pixels = append(pixels, img.RGBAAt(surroundingXPixel, surroundingYPixel))
			}
		}
	}
	return pixels
}

func setAverageColor(px, py int, pixelEnvironment []color.RGBA) {
	var r, g, b, a int
	for _, color := range pixelEnvironment {
		r += int(color.R)
		g += int(color.G)
		b += int(color.B)
		a += int(color.A)
	}

	newColor := color.RGBA{
		uint8(r / len(pixelEnvironment)),
		uint8(g / len(pixelEnvironment)),
		uint8(b / len(pixelEnvironment)),
		uint8(a / len(pixelEnvironment)),
	}

	img.Set(px, py, newColor)
}
func mandelbrot(z complex128) color.RGBA {
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		var r uint8 = uint8(rand.Intn(256))
		var g uint8 = uint8(rand.Intn(256))
		var b uint8 = uint8(rand.Intn(256))
		//  If absolute value of complex number > 2 (the point is allocated out of circe that has r = 2)
		// We'll color the pixel with shade of the gray color
		if cmplx.Abs(v) > 2*math.Pow(zoom, 2) {
			// 255 is the absolute value of white color.
			// Subracting values from 255 we get shades of gray color
			return color.RGBA{r, g, b, 255 - uint8(rand.Intn(256))}
		}
	}
	// If absolute value of complex number <= 2 (the point is allocated inside the cirlce that has r = 2)
	// We'll color the point with black
	return color.RGBA{0, 0, 255, uint8(rand.Intn(256))}
}

func newton(initialGuess complex128) color.RGBA {
	const epsilon float64 = 1e-5 // Precision of float64
	const maxIterations = 10000  // Max amount of iterations

	var countOfCycles int64 // We'll count this value

	z := initialGuess // Initial point on the complex plane

	// Depending of the found approximation we'll match the color
	rootColors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{255, 255, 0, 255}, // Yellow
	}

	for countOfCycles = 0; cmplx.Abs(f(z)) > epsilon && countOfCycles < maxIterations; countOfCycles++ {
		z -= f(z) / fpr(z)
	}

	// If the iteration count is below the maximum, find the index of the root that was approached
	if countOfCycles < maxIterations {
		rootIndex := findApproachedRoot(z)
		return rootColors[rootIndex]
	}

	// If the iteration count exceeds the maximum, return a default color (you may adjust this)
	return color.RGBA{0, 0, 0, 255}
}

func findApproachedRoot(z complex128) int {
	roots := []complex128{
		cmplx.Pow(complex(1, 0), 4),
		cmplx.Pow(complex(0, 1), 4),
		cmplx.Pow(complex(-1, 0), 4),
		cmplx.Pow(complex(0, -1), 4),
	}

	// Find the index of the root that z is closest to
	minDistance := cmplx.Abs(z - roots[0])
	minIndex := 0
	for i := 1; i < len(roots); i++ {
		distance := cmplx.Abs(z - roots[i])
		if distance < minDistance {
			minDistance = distance
			minIndex = i
		}
	}

	return minIndex
}

func f(z complex128) complex128 {
	return cmplx.Pow(z, 4) - 1
}

func fpr(z complex128) complex128 {
	return 4 * cmplx.Pow(z, 3)
}

func rotate90Right(original *image.RGBA) {
	width := original.Bounds().Dx()
	height := original.Bounds().Dy()

	// Создаем новое изображение с перевернутыми размерами
	rotated := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Копируем пиксели из оригинального изображения с учетом поворота
			rotated.Set(y, width-x, original.At(x, y))
		}
	}

	img = rotated
}
