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
	// sqKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 21, Y: 21})

	fname := "testdata/passport_01.jpg"

	mat := gocv.IMRead(fname, gocv.IMReadUnchanged)
	mat = mrz.Resize(mat, 0, 600, gocv.InterpolationArea)

	gray := mrz.Color(mat, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &gray, image.Point{X: 3, Y: 3}, 0, 0, gocv.BorderConstant)

	blackhat := gray.Clone()
	gocv.MorphologyEx(gray, &blackhat, gocv.MorphBlackhat, rectKernel)

	// Python:
	// gradX = cv2.Sobel(blackhat, ddepth=cv2.CV_32F, dx=1, dy=0, ksize=-1)
	// gradX = np.absolute(gradX)
	// (minVal, maxVal) = (np.min(gradX), np.max(gradX))
	// gradX = (255 * ((gradX - minVal) / (maxVal - minVal))).astype("uint8")

	// TODO:
	gradX := blackhat.Clone()
	gocv.Sobel(blackhat, &gradX, gocv.MatTypeCV32F, 1, 0, -1, 0, 0, gocv.BorderConstant)

	img, err := gradX.ToImage()
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
	window.IMShow(mat)
	window.WaitKey(0)
}
