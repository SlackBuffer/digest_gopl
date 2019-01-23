// Generates random animations of random Lissajous figures
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var palette = []color.Color{color.White, color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff}, color.RGBA{0xe3, 0xe3, 0xe3, 0xff}} // slice

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
	greenIndex = 2
	redIndex   = 3
	blueIndex  = 4
	greyIndex  = 5
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64 // number of animation frames, 帧数
		delay   = 8
	)

	anim := gif.GIF{LoopCount: nframes} // struct, init LoopCount field with nframe, other fields stay the default zero value

	rand.Seed(time.Now().UTC().UnixNano())
	freq := rand.Float64() * 3.0
	phase := 0.0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) // 201*201 画板
		// 调色板；所有像素被初始化为画板的 zero value，即第 0 个颜色（白色）
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(rand.Float64()*4.0)+1)
		}
		phase += 0.1

		// update Delay and Image fields of anim
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

// go build main.go && ./main >out.gif
