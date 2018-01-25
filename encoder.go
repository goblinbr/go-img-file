package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	file, err1 := os.Open("teste0.png")
	fmt.Println(file)
	fmt.Println(err1)
	if err1 == nil {
		defer file.Close()
		//imageData, imageType, err2 := image.Decode(file)
		imageData, err2 := png.Decode(file)
		if err2 == nil {
			bounds := imageData.Bounds()
			width := bounds.Max.X
			height := bounds.Max.Y

			var maxUsedColor color.Color
			maxUsedCount := 0
			var colorCountMap = make(map[color.Color]int)
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					col := imageData.At(x, y)
					var count = colorCountMap[col]
					count++
					colorCountMap[col] = count
					if count > maxUsedCount {
						maxUsedCount = count
						maxUsedColor = col
					}
					r, g, b, a := col.RGBA()
					fmt.Printf("%d,%d,%d,%d = %d\n", r, g, b, a, count)
				}
			}

			fmt.Printf("máx. used color = %s\n", convertColorToString(maxUsedColor))
		} else {
			fmt.Println(err2)
		}
	} else {
		fmt.Println(err1)
	}

}
