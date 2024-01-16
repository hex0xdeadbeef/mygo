package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
)

func Server() {
	http.HandleFunc("/mandelbrot", getMandelbrotHandler)
	http.HandleFunc("/newton", getNewtonHandler)
	log.Fatalf("Error occured: %v", http.ListenAndServe("localhost:8000", nil))
}

func getMandelbrotHandler(writer http.ResponseWriter, request *http.Request) {
	createFractal(writer, "mandelbrot")
}
func getNewtonHandler(writer http.ResponseWriter, request *http.Request) {
	createFractal(writer, "newton")
}

const (
	width, height          = 1024, 1024
	xmin, ymin, xmax, ymax = -1, -1, 1, 1
)

func createFractal(writer io.Writer, method string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			img.Set(px, py, getAverageColor(method, px, py))
		}
	}
	png.Encode(writer, rotate90Right(img))
}

func getAverageColor(method string, px, py int) color.RGBA {
	var pixelEnvironmentColors []color.RGBA
	for i := 0; i < 3; i++ {
		px := px - 1 + i
		// Scaling the current y coordinate so that it corresponds to location on complex plane
		x := (float64(px)/width)*(xmax-xmin) + xmin
		for j := 0; j < 3; j++ {
			py := py - 1 + j
			// Scaling the current x coordinate so that it corresponds to location on complex plane
			y := (float64(py)/height)*(ymax-ymin) + ymin
			z := complex(-x, y)
			switch method {
			case "mandelbrot":
				pixelEnvironmentColors = append(pixelEnvironmentColors, mandelbrot(z))
			case "newton":
				pixelEnvironmentColors = append(pixelEnvironmentColors, newton(z))
			}
		}
	}

	var r, g, b, a int
	for _, color := range pixelEnvironmentColors {
		r += int(color.R)
		g += int(color.G)
		b += int(color.B)
		a += int(color.A)
	}

	return color.RGBA{uint8(r / 9), uint8(g / 9), uint8(b / 9), uint8(a / 9)}
}

func mandelbrot(z complex128) color.RGBA {
	const (
		iterations = 200
		contrast   = 15
	)

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		var r uint8 = 0 + n*contrast
		var b uint8 = 255 - n*contrast
		//  If absolute value of complex number > 2 (the point is allocated out of circe that has r = 2)
		// We'll color the pixel with shade of the gray color
		if cmplx.Abs(v) > 2 {
			// 255 is the absolute value of white color.
			// Subracting values from 255 we get shades of gray color
			return color.RGBA{r, 0, b, 255 - n*contrast}
		}
	}
	// If absolute value of complex number <= 2 (the point is allocated inside the cirlce that has r = 2)
	// We'll color the point with black
	return color.RGBA{0, 255, 0, 255}
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

func rotate90Right(original *image.RGBA) *image.RGBA {
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

	return rotated
}
