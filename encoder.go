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
	cor        color.Color
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
						colorUsed.cor = col
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
							if colorUsed.cor == col {
								if x == width-1 {
									colorUsed.rectangles = append(colorUsed.rectangles, image.Rect(x, y, x, y))
								}
								for x2 := x + 1; x2 < width; x2++ {
									col := imageData.At(x2, y)
									if colorUsed.cor != col {
										colorUsed.rectangles = append(colorUsed.rectangles, image.Rect(x, y, x2-1, y))
										x = x2 - 1
										break
									}
									if x2 == width-1 {
										colorUsed.rectangles = append(colorUsed.rectangles, image.Rect(x, y, x2, y))
										x = x2
										break
									}
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
			f.WriteString(convertColorToString(colorUsed.cor) + "=")
			for _, rect := range colorUsed.rectangles {
				f.WriteString(convertRectangleToString(rect) + ";")
			}
		} else {
			f.WriteString(convertColorToString(colorUsed.cor))
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
