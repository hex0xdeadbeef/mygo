package threedsurface

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	width, height        = 600, 320                     // canvas size in pixels
	cells                = 100                          // number of grid
	xyrange              = 30.0                         // axis ranges (-xyrange...+xyrange)
	xyscale              = float64(width) / 2 / xyrange // pixels per x or y unit
	zscale               = float64(height) * 0.4        // pixels per z unit
	angle                = math.Pi / 6                  // angle of x, y axes (=30*)
	sin30, cos30         = math.Sin(angle), math.Cos(angle)
	fuction       string = "sin"
)

func getSurface(writer http.ResponseWriter, request *http.Request) {
	parameters := request.URL.Query()

	// Initialize request parameters
	for parameter := range parameters {
		switch parameter {
		case "height":
			parseNumberParameter(parameters[parameter][0], &height)
		case "width":
			parseNumberParameter(parameters[parameter][0], &width)
		case "function":
			parseStringParameter(parameters[parameter][0])
		default:
			break
		}
	}

	// Write type of content we'll write
	writer.Header().Set("Content-Type", "image/svg+xml")

	// Set the scale factor (adjust as needed)
	scale := 1.0

	// Write the svg/xml top info with centering
	fmt.Fprintf(writer, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: red; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d' viewBox='0 0 %d %d'>", width, height, width, height)

	// Find and reflect points on 2D canvas
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)

			// Apply scaling
			ax *= scale
			bx *= scale
			cx *= scale
			dx *= scale
			ay *= scale
			by *= scale
			cy *= scale
			dy *= scale

			// Calculate average height for the polygon
			h1, _ := findZCoordinate(fuction, ax, ay)
			h2, _ := findZCoordinate(fuction, bx, by)
			h3, _ := findZCoordinate(fuction, cx, cy)
			h4, _ := findZCoordinate(fuction, dx, dy)
			avgHeight := (h1 + h2 + h3 + h4) / 4

			// Map the height to a color from blue to red
			color := colorMap(avgHeight)

			// Output polygon with fill color
			fmt.Fprintf(writer, "<polygon points='%g, %g, %g, %g, %g, %g, %g, %g' fill='%s' />\n", ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Fprintf(writer, "</svg>\n")
}

// Find points to be reflected on 2D canvas
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i)/float64(cells) - 0.5)
	y := xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z
	z, err := findZCoordinate(fuction, x, y)
	if err != nil {
		return 0, 0 // handle the error, for now, just return (0, 0)
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx, sy)
	sx := float64(width/2) + (x-y)*cos30*xyscale
	sy := float64(height/2) + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

// Find the z respective coorditate
func findZCoordinate(f string, x, y float64) (float64, error) {
	r := math.Hypot(x, y) // distance from (0, 0)
	var result float64

	switch f {
	case "sin":
		result = math.Sin(r) / r
	case "cos":
		result = math.Cos(r) / r
	case "tan":
		result = math.Tan(r) / r
	case "atan":
		result = math.Atan(r) / r
	case "sin^2":
		result = math.Pow(math.Sin(r), 2) / r
	case "cos^2":
		result = math.Pow(math.Cos(r), 2) / r
	case "sin*cos":
		result = math.Cos(r) * math.Sin(r) / r
	}

	if math.IsInf(result, 0) {
		return 0, fmt.Errorf("invalid result for z coordinate")
	}
	return result, nil
}

// Returns the color that depends on the height
func colorMap(height float64) string {
	// Map height to a color from blue to red
	normalizedHeight := (height + 1) / 2   // Normalize height to the range [0, 1]
	r := int(255 * normalizedHeight)       // Red component
	b := int(255 * (1 - normalizedHeight)) // Blue component
	g := 0                                 // Green component (you can customize this)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// Parse and set the reques parameters if it's possible
func parseNumberParameter(parameter string, ptr *int) {
	if parameter != "" {
		if value, err := strconv.ParseFloat(parameter, 64); err != nil {
			fmt.Fprintf(os.Stdout, "The error occured: ", err)
			os.Exit(1)
		} else {
			*ptr = int(value)
		}
	}
}

// Parse and set the reques parameters if it's possible
func parseStringParameter(parameter string) {
	if parameter != "" &&
		(strings.Contains(parameter, "sin")) || (strings.Contains(parameter, "cos")) ||
		(strings.Contains(parameter, "tan")) || (strings.Contains(parameter, "atan")) ||
		(strings.Contains(parameter, "sin^2")) || (strings.Contains(parameter, "cos^2")) {
		fuction = parameter
	} else {
		fmt.Fprintf(os.Stdout, "Server does not have this parameter")
		os.Exit(1)
	}
}

// Start the local server and catches errors if they occur
func LocalServer() {
	http.HandleFunc("/getsurface", surfaceHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Handle the /getsurface request
func surfaceHandler(writer http.ResponseWriter, request *http.Request) {
	getSurface(writer, request)
}
