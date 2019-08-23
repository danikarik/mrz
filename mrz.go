package mrz

import (
	"errors"
	"image"

	"gocv.io/x/gocv"
)

// Version prints gocv versions.
func Version() string { return gocv.Version() }

// OpenCVVersion prints opencv versions.
func OpenCVVersion() string { return gocv.OpenCVVersion() }

// Resize image with ratio.
func Resize(mat gocv.Mat, width, height float64, interp gocv.InterpolationFlags) gocv.Mat {
	if width == 0 && height == 0 {
		return mat
	}

	w, h := mat.Cols(), mat.Rows()

	if height > 0 {
		r := height / float64(h)
		w, h = int(float64(w)*r), int(height)
	} else {
		r := width / float64(w)
		w, h = int(width), int(float64(h)*r)
	}

	gocv.Resize(mat, &mat, image.Point{X: w, Y: h}, 0, 0, interp)
	return mat
}

// GrabContours returns points depends on OpenCV.
func GrabContours(cnts [][]image.Point) ([]image.Point, error) {
	if len(cnts) == 2 {
		return cnts[0], nil
	}
	if len(cnts) == 3 {
		return cnts[1], nil
	}
	return nil, errors.New("contours must have length 2 or 3")
}
