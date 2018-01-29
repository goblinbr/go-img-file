package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ColorUsed - Colors and bounds
type ColorUsed struct {
	color      color.Color
	count      int
	rectangles []image.Rectangle
}

// OrderByCount - order by count desc
type OrderByCount []ColorUsed

func (a OrderByCount) Len() int           { return len(a) }
func (a OrderByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderByCount) Less(i, j int) bool { return a[i].count > a[j].count }

func getColor255(c uint32) int {
	return int(c / 0x101)
}

func convertColorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	if a == 65535 {
		return strings.Join([]string{strconv.Itoa(getColor255(r)), ",", strconv.Itoa(getColor255(g)), ",", strconv.Itoa(getColor255(b))}, "")
	}
	return strings.Join([]string{strconv.Itoa(getColor255(r)), ",", strconv.Itoa(getColor255(g)), ",", strconv.Itoa(getColor255(b)), ",", strconv.Itoa(getColor255(a))}, "")
}

func convertRectangleToString(r image.Rectangle) string {
	x1 := r.Min.X
	y1 := r.Min.Y
	x2 := r.Max.X
	y2 := r.Max.Y
	if y1 == y2 {
		if x1 == x2 {
			return strconv.Itoa(x1) + "," + strconv.Itoa(y1)
		}
		width := x2 - x1
		return strconv.Itoa(x1) + "," + strconv.Itoa(y1) + "," + strconv.Itoa(width)
	}
	return strconv.Itoa(x1) + "," + strconv.Itoa(y1) + "," + strconv.Itoa(x2) + "," + strconv.Itoa(y2)
}

func inRectangles(rcs []image.Rectangle, x int, y int) bool {
	for _, rect := range rcs {
		if rect.Min.X <= x && rect.Max.X >= x && rect.Min.Y <= y && rect.Max.Y >= y {
			return true
		}
	}
	return false
}

func createColorRect(imageData image.Image, colorUsed *ColorUsed, x0 int, y0 int) image.Rectangle {
	bounds := imageData.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	col := colorUsed.color

	x1 := width - 1
	for x := x0; x < width; x++ {
		c := imageData.At(x, y0)
		if c != col || inRectangles(colorUsed.rectangles, x, y0) {
			x1 = x - 1
			break
		}
	}

	y1 := height - 1
	dif := false
	for y := y0 + 1; y < height && !dif; y++ {
		for x := x0; x <= x1; x++ {
			c := imageData.At(x, y)
			if c != col || inRectangles(colorUsed.rectangles, x, y) {
				y1 = y - 1
				dif = true
				break
			}
		}
	}
	return image.Rect(x0, y0, x1, y1)
}

func readFile(arq string) ([]ColorUsed, error) {
	var colors []ColorUsed
	file, err1 := os.Open(arq)
	fmt.Println(file)
	fmt.Println(err1)
	if err1 == nil {
		defer file.Close()
		imageData, err2 := png.Decode(file)
		if err2 == nil {
			bounds := imageData.Bounds()
			width := bounds.Max.X
			height := bounds.Max.Y

			var colorCountMap = make(map[color.Color]*ColorUsed)
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					col := imageData.At(x, y)
					colorUsed := colorCountMap[col]
					if colorUsed == nil {
						colorUsed = new(ColorUsed)
						colorUsed.color = col
						colorUsed.count = 1
						colorUsed.rectangles = []image.Rectangle{}

						colorCountMap[col] = colorUsed
					} else {
						colorUsed.count++
					}
				}
			}
			fmt.Println()

			colors = make([]ColorUsed, len(colorCountMap))
			i := 0
			for _, value := range colorCountMap {
				colors[i] = *value
				i++
			}

			sort.Sort(OrderByCount(colors))

			for i := range colors {
				colorUsed := &colors[i]
				if i > 0 {
					for y := 0; y < height; y++ {
						for x := 0; x < width; x++ {
							col := imageData.At(x, y)
							if colorUsed.color == col {
								if !inRectangles(colorUsed.rectangles, x, y) {
									colorUsed.rectangles = append(colorUsed.rectangles, createColorRect(imageData, colorUsed, x, y))
								}
							}
						}
					}
				}
			}
		} else {
			return nil, err2
		}
	} else {
		return nil, err1
	}
	return colors, nil
}

func writeFile(arq string, colors []ColorUsed) error {
	f, err := os.Create(arq)
	if err != nil {
		return err
	}
	defer f.Close()

	totalBytes := 0
	for i, colorUsed := range colors {
		if i > 0 {
			f.WriteString("|")
			f.WriteString(convertColorToString(colorUsed.color) + "=")
			for _, rect := range colorUsed.rectangles {
				f.WriteString(convertRectangleToString(rect) + ";")
			}
		} else {
			f.WriteString(convertColorToString(colorUsed.color))
		}
	}
	fmt.Println(totalBytes)
	return nil
}

func main() {
	colors, err := readFile("teste0.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(colors)

	writeFile("teste0.rbs", colors)
}
