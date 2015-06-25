// package composite contains helpful functions for compositing information
// onto video in real time
package composite

import (
	"bufio"
	"code.google.com/p/freetype-go/freetype"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

func RenderTextToPNG(text string, filename string, fontPath string) error {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	fg, bg := image.White, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	fsize := 32.0
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(140)
	c.SetFont(font)
	c.SetFontSize(fsize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(freetype.NoHinting)

	pt := freetype.Pt(100, 100+int(c.PointToFix32(fsize)>>8))
	_, err = c.DrawString(text, pt)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, rgba)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}

	return nil
}
