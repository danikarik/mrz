package mrz

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"

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

// Color converts color.
func Color(mat gocv.Mat, code gocv.ColorConversionCode) gocv.Mat {
	gocv.CvtColor(mat, &mat, code)
	return mat
}

// ScanImageFromBytes reads MRZ from given content.
func ScanImageFromBytes(data []byte) ([]byte, error) {
	file, err := ioutil.TempFile("", "mrz_")
	if err != nil {
		return nil, err
	}

	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	mat := gocv.IMRead(file.Name(), gocv.IMReadUnchanged)
	mat = Resize(mat, 0, 600, gocv.InterpolationArea)

	gray := Color(mat, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &gray, image.Point{X: 3, Y: 3}, 0, 0, gocv.BorderConstant)

	img, err := mat.ToImage()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	err = os.Remove(file.Name())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ScanImage reads MRZ from given source.
func ScanImage(fname string) error {
	return errors.New("not implemented yet")
}
