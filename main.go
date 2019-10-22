package imgedit

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// Region 区域
type Region struct {
	XStart int
	XEnd   int
	YStart int
	YEnd   int
}

// Weight 图宽
func Weight(beforeImg string) int {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	return bounds.Dx()
}

// Height 图高
func Height(beforeImg string) int {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	return bounds.Dy()
}

// openImg 打开图片
func openImg(beforeImg string) image.Image {
	// 打开文件
	f, err := os.Open(beforeImg)
	if err != nil {
		panic("打开图片失败")
	}
	defer f.Close()

	// 图片解码
	m, _, err := image.Decode(f)
	if err != nil {
		panic("图片解码失败")
	}

	return m
}

// createNew 生成新图片
func createNew(afterImg string, afterRgba *image.RGBA) {
	c, err := os.Create(afterImg)
	if err != nil {
		panic("生成新图片失败")
	}
	defer c.Close()
	png.Encode(c, afterRgba)
}

// AntiColor 反色
func AntiColor(beforeImg, afterImg string) {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	afterRgba := image.NewRGBA(bounds)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			rUint8 := 255 - uint8(r>>8)
			gUint8 := 255 - uint8(g>>8)
			bUint8 := 255 - uint8(b>>8)
			aUint8 := uint8(a >> 8)

			//设置像素点
			afterRgba.SetRGBA(x, y, color.RGBA{rUint8, gUint8, bUint8, aUint8})
		}
	}

	createNew(afterImg, afterRgba)
}

// Grayscale 灰度
func Grayscale(beforeImg, afterImg string) {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	afterRgba := image.NewRGBA(bounds)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()

			//加权平均并转换为 255 值
			rUint8 := uint8(int(float64(r)*0.299) >> 8)
			gUint8 := uint8(int(float64(g)*0.587) >> 8)
			bUint8 := uint8(int(float64(b)*0.114) >> 8)
			gray := rUint8 + gUint8 + bUint8
			aUint8 := uint8(a >> 8)

			//设置像素点
			afterRgba.SetRGBA(x, y, color.RGBA{gray, gray, gray, aUint8})
		}
	}

	createNew(afterImg, afterRgba)
}

// FrostedGlass 毛玻璃
func FrostedGlass(beforeImg, afterImg string, pixel int) {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	afterRgba := image.NewRGBA(bounds)

	// 随机种子数
	rand.Seed(time.Now().UnixNano())

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			mblX := rand.Intn(pixel) + x
			mblY := rand.Intn(pixel) + y

			// 处理边缘不满取值范围大小问题
			if (x + pixel) > width {
				mblX = rand.Intn(width-x) + x
			}
			if (y + pixel) > height {
				mblY = rand.Intn(height-y) + y
			}

			colorRgb := m.At(mblX, mblY)
			r, g, b, a := colorRgb.RGBA()

			rUint8 := uint8(r >> 8)
			gUint8 := uint8(g >> 8)
			bUint8 := uint8(b >> 8)
			aUint8 := uint8(a >> 8)

			//设置像素点
			afterRgba.SetRGBA(x, y, color.RGBA{rUint8, gUint8, bUint8, aUint8})
		}
	}

	createNew(afterImg, afterRgba)
}

// Mosaic 马赛克
func Mosaic(beforeImg, afterImg string, pixel int, region Region) {
	m := openImg(beforeImg)
	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 判断马赛克区域是否合法
	if (region.XStart > width) || (region.XEnd > width) {
		panic("x轴超出范围")
	}
	if (region.YStart > height) || (region.YEnd > height) {
		panic("y轴超出范围")
	}

	// 先生成新图片
	afterRgba := image.NewRGBA(bounds)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			rUint8 := uint8(r >> 8)
			gUint8 := uint8(g >> 8)
			bUint8 := uint8(b >> 8)
			aUint8 := uint8(a >> 8)
			afterRgba.SetRGBA(i, j, color.RGBA{rUint8, gUint8, bUint8, aUint8})
		}
	}

	// 随机种子数
	rand.Seed(time.Now().UnixNano())

	// 处理马赛克区域
	for x := region.XStart; x < region.XEnd; x += pixel {
		for y := region.YStart; y < region.YEnd; y += pixel {

			mskX := rand.Intn(pixel) + x
			mskY := rand.Intn(pixel) + y

			// 处理边缘不满取值范围大小问题
			if (x + pixel) > width {
				mskX = rand.Intn(width-x) + x
			}
			if (y + pixel) > height {
				mskY = rand.Intn(height-y) + y
			}

			colorRgb := m.At(mskX, mskY)
			r, g, b, a := colorRgb.RGBA()
			rUint8 := uint8(r >> 8)
			gUint8 := uint8(g >> 8)
			bUint8 := uint8(b >> 8)
			aUint8 := uint8(a >> 8)

			//设置马赛克方格像素点
			for k := 0; k < pixel; k++ {
				for p := 0; p < pixel; p++ {
					afterRgba.SetRGBA(x+k, y+p, color.RGBA{rUint8, gUint8, bUint8, aUint8})
				}
			}
		}
	}

	createNew(afterImg, afterRgba)
}
