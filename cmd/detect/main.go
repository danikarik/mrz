package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/danikarik/mrz"
	"gocv.io/x/gocv"
)

func main() {
	rectKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 13, Y: 5})
	sqKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 21, Y: 21})

	original := gocv.IMRead("testdata/passport_01.jpg", gocv.IMReadUnchanged)
	dst := original.Clone()

	dst = mrz.Resize(dst, 0, 600, gocv.InterpolationArea)
	dst = mrz.Color(dst, gocv.ColorBGRToGray)

	gocv.GaussianBlur(dst, &dst, image.Point{X: 3, Y: 3}, 0, 0, gocv.BorderConstant)
	gocv.MorphologyEx(dst, &dst, gocv.MorphBlackhat, rectKernel)

	gocv.Sobel(dst, &dst, 0, 1, 0, 1, 1, 0, gocv.BorderConstant)
	gocv.ConvertScaleAbs(dst, &dst, 1, 0)

	gocv.MorphologyEx(dst, &dst, gocv.MorphClose, rectKernel)
	gocv.Threshold(dst, &dst, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	gocv.MorphologyEx(dst, &dst, gocv.MorphClose, sqKernel)
	m := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV8U)
	gocv.Erode(dst, &dst, m)

	img, err := dst.ToImage()
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("./output.png")
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}

	window := gocv.NewWindow("Hello")
	window.IMShow(dst)
	window.WaitKey(0)
}
