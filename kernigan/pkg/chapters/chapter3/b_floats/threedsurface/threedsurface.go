package threedsurface

import (
	"fmt"
	"math"
	"net/http"
)

var (
	width, height = 600, 320                     // canvas size in pixels
	cells         = 100                          // number of grid
	xyrange       = 30.0                         // axis ranges (-xyrange...+xyrange)
	xyscale       = float64(width) / 2 / xyrange // pixels per x or y unit
	zscale        = float64(height) * 0.4        // pixels per z unit
	angle         = math.Pi / 6                  // angle of x, y axes (=30*)
	sin30, cos30  = math.Sin(angle), math.Cos(angle)
	function      func(x, y float64) float64
)

func GetSurface(writer http.ResponseWriter, fn func(x, y float64) float64) {

	function = fn

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
			h1 := function(ax, ay)
			h2 := function(bx, by)
			h3 := function(cx, cy)
			h4 := function(dx, dy)
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
func corner(i, j int) (a float64, b float64) {
	// Find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i)/float64(cells) - 0.5)
	y := xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z
	z := function(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx, sy)
	a = float64(width/2) + (x-y)*cos30*xyscale
	b = float64(height/2) + (x+y)*sin30*xyscale - z*zscale
	return
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
