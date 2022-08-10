package main

import (
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face1 := truetype.NewFace(font, &truetype.Options{Size: 50})
	face2 := truetype.NewFace(font, &truetype.Options{Size: 16})
	dc := gg.NewContext(1024, 1024)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	marginH := 250.0
	marginV := 490.0
	x := marginH
	y := marginV
	w := float64(dc.Width()) - (2.0 * marginH)
	h := float64(dc.Height()) - (2.0 * marginV)
	dc.SetColor(color.RGBA{0, 191, 255, 0})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.SetFontFace(face1)
	dc.DrawStringAnchored("Hello, ", 348, 512, 0.5, 0.5)
	dc.SetFontFace(face2)
	dc.DrawStringAnchored("the best programmer in the", 512, 524, 0.5, 0.5)
	dc.SetFontFace(face1)
	dc.DrawStringAnchored(" world!", 676, 512, 0.5, 0.5)

	dc.SavePNG("out.png")
}
