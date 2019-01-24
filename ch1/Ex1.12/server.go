package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var palette = []color.Color{color.White, color.RGBA{0x00, 0xff, 0x00, 0xff}} // slice

const (
	whiteIndex = 0 // first color in palette
	greenIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}

		for _, v := range r.Form {
			// fmt.Println(v[0])
			// fmt.Printf("%T", v[0]) // string
			circles, err := strconv.Atoi(v[0])
			if err != nil {
				log.Print(err)
			}
			lissajous(w, circles)
		}
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, cc int) {
	const (
		res     = 0.001
		size    = 100
		nframes = 64 // number of animation frames, 帧数
		delay   = 8
	)

	cycles := float64(cc)

	anim := gif.GIF{LoopCount: nframes} // struct, init LoopCount field with nframe, other fields stay the default zero value

	rand.Seed(time.Now().UTC().UnixNano())
	freq := rand.Float64() * 3.0
	phase := 0.0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1) // 201*201 画板
		// 调色板；所有像素被初始化为画板的 zero value，即第 0 个颜色（白色）
		img := image.NewPaletted(rect, palette)

		// 将某些像素点设为绿色
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			// 将某些像素点设为绿色
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), greenIndex)
		}
		phase += 0.1

		// update Delay and Image fields of anim
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

// http://localhost:8000/?cycles=20
