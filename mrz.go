package mrz

import (
	"errors"
	"image"
	"io/ioutil"
	"sort"

	"gocv.io/x/gocv"
)

// ErrNotFound returned if mrz of id/passport not found.
var ErrNotFound = errors.New("mrz: not found")

// Version prints gocv versions.
func Version() string { return gocv.Version() }

// OpenCVVersion prints opencv versions.
func OpenCVVersion() string { return gocv.OpenCVVersion() }

// Detect tries to detect mrz from file.
func Detect(filename string) (image.Image, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return DetectFromBytes(src)
}

// DetectFromBytes tries to detect mrz from bytes source.
func DetectFromBytes(source []byte) (image.Image, error) {
	rectKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 13, Y: 5})
	sqKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 21, Y: 21})

	original, err := gocv.IMDecode(source, gocv.IMReadUnchanged)
	if err != nil {
		return nil, err
	}
	original = resize(original, 0, 600, gocv.InterpolationArea)

	gray := cvtColor(original, gocv.ColorBGRToGray)
	gray = gaussianBlur(gray, image.Point{X: 3, Y: 3}, 0, 0, gocv.BorderDefault)
	blackhat := morphologyEx(gray, gocv.MorphBlackhat, rectKernel)

	gradX := sobel(blackhat, gocv.MatTypeCV32F, 1, 0, -1, 1, 0, gocv.BorderDefault)
	gradX = absolute(gradX)
	minVal, maxVal := min(gradX), max(gradX)
	gradX = computeScharrGradient(gradX, minVal, maxVal)
	gradX = morphologyEx(gradX, gocv.MorphClose, rectKernel)

	thresh := threshold(gradX, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	thresh = morphologyEx(thresh, gocv.MorphClose, sqKernel)
	thresh = erode(thresh, 4)

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

	cnts := gocv.FindContours(thresh.Clone(), gocv.RetrievalExternal, gocv.ChainApproxSimple)
	sort.SliceStable(cnts, func(i, j int) bool {
		return gocv.ContourArea(cnts[i]) > gocv.ContourArea(cnts[j])
	})

	for _, c := range cnts {
		rect := gocv.BoundingRect(c)
		x, y, w, h := rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy()
		ar := float64(w) / float64(h)
		crWidth := float64(w) / float64(gray.Size()[1])

		if ar > 5 && crWidth > 0.75 {
			pX := int(float64(x+w) * 0.03)
			pY := int(float64(y+h) * 0.03)
			x, y = x-pX, y-pY
			w, h = w+(pX*2), h+(pY*2)

			box := image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + w, Y: y + h},
			}

			img, err := original.ToImage()
			if err != nil {
				return nil, err
			}

			roi := crop(img, box)
			return roi, nil
		}
	}

	return nil, ErrNotFound
}
