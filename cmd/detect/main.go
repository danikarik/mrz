package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sort"

	"github.com/danikarik/mrz"
	"gocv.io/x/gocv"
)

func main() {
	rectKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 13, Y: 5})
	sqKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 21, Y: 21})

	original := gocv.IMRead("testdata/passport_01.jpg", gocv.IMReadUnchanged)
	original = mrz.Resize(original, 0, 600, gocv.InterpolationArea)

	gray := original.Clone()
	gocv.CvtColor(original, &gray, gocv.ColorBGRToGray)

	gocv.GaussianBlur(gray, &gray, image.Point{X: 3, Y: 3}, 0, 0, gocv.BorderConstant)
	blackhat := gray.Clone()
	gocv.MorphologyEx(gray, &blackhat, gocv.MorphBlackhat, rectKernel)

	gradX := blackhat.Clone()
	gocv.Sobel(blackhat, &gradX, 0, 1, 0, 1, 1, 0, gocv.BorderConstant)
	gocv.ConvertScaleAbs(gradX, &gradX, 1, 0)

	gocv.MorphologyEx(gradX, &gradX, gocv.MorphClose, rectKernel)
	thresh := gradX.Clone()
	gocv.Threshold(gradX, &thresh, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	gocv.MorphologyEx(thresh, &thresh, gocv.MorphClose, sqKernel)
	m := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV8U)
	gocv.Erode(thresh, &thresh, m)

	p := int(float64(thresh.Size()[1]) * 0.05)
	for i := 0; i < thresh.Rows(); i++ {
		for j := 0; j < thresh.Cols(); j++ {
			if j < p {
				thresh.SetUCharAt(i, j, 0)
			}
		}
	}
	for i := 0; i < thresh.Rows(); i++ {
		for j := 0; j < thresh.Cols(); j++ {
			if original.Size()[1]-p <= j {
				thresh.SetUCharAt(i, j, 0)
			}
		}
	}

	// cnts = cv2.findContours(thresh.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
	// cnts = imutils.grab_contours(cnts)
	// cnts = sorted(cnts, key=cv2.contourArea, reverse=True)
	// TODO:
	cnts := gocv.FindContours(thresh.Clone(), gocv.RetrievalExternal, gocv.ChainApproxSimple)
	sort.SliceStable(cnts, func(i, j int) bool {
		return gocv.ContourArea(cnts[i]) > gocv.ContourArea(cnts[j])
	})

	for _, c := range cnts {
		rect := gocv.BoundingRect(c)
		x, y, w, h := rect.Size().X, rect.Size().Y, rect.Dx(), rect.Dy()
		ar := float64(w) / float64(h)
		crWidth := float64(w) / float64(gray.Size()[1])

		if ar > 5 && crWidth > 0.75 {
			pX := int(float64(x+w) * 0.03)
			pY := int(float64(y+h) * 0.03)
			x, y = x-pX, y-pY
			w, h = w+(pX*2), h+(pY*2)

			// roi = image[y:y + h, x:x + w].copy()
			// cv2.rectangle(image, (x, y), (x + w, y + h), (0, 255, 0), 2)
			// TODO:
			greenBox := image.Rectangle{
				Min: image.Point{
					X: x,
					Y: y,
				},
				Max: image.Point{
					X: x + w,
					Y: y + h,
				},
			}
			greenColor := color.RGBA{R: 0, G: 230, B: 64, A: 1}
			gocv.Rectangle(&original, greenBox, greenColor, 2)
		}
	}

	img, err := original.ToImage()
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

	// window := gocv.NewWindow("Image")
	// window.IMShow(original)
	// window.WaitKey(0)
}
