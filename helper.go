package mrz

import (
	"image"
	"image/draw"
	"math"

	"gocv.io/x/gocv"
)

func resize(mat gocv.Mat, width, height float64, interp gocv.InterpolationFlags) gocv.Mat {
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

func absolute(mat gocv.Mat) gocv.Mat {
	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			v := mat.GetFloatAt(i, j)
			abs := math.Float32frombits(math.Float32bits(v) &^ (1 << 31))
			mat.SetFloatAt(i, j, abs)
		}
	}

	return mat
}

func min(mat gocv.Mat) float32 {
	min := float32(math.MaxFloat32)

	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			v := mat.GetFloatAt(i, j)
			if v < min {
				min = v
			}
		}
	}

	return min
}

func max(mat gocv.Mat) float32 {
	max := float32(math.SmallestNonzeroFloat32)

	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			v := mat.GetFloatAt(i, j)
			if v > max {
				max = v
			}
		}
	}

	return max
}

func computeScharrGradient(mat gocv.Mat, minVal, maxVal float32) gocv.Mat {
	dst := gocv.NewMat()

	mat.SubtractFloat(minVal)
	mat.DivideFloat(maxVal - minVal)
	mat.MultiplyFloat(float32(255))
	mat.ConvertTo(&dst, gocv.MatTypeCV8U)

	return dst
}

func erode(mat gocv.Mat, iteration int) gocv.Mat {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 3, Y: 3})

	for i := 0; i < 4; i++ {
		gocv.Erode(mat, &mat, kernel)
	}

	return mat
}

func cvtColor(mat gocv.Mat, c gocv.ColorConversionCode) gocv.Mat {
	dst := gocv.NewMat()
	gocv.CvtColor(mat, &dst, c)
	return dst
}

func gaussianBlur(mat gocv.Mat, ksize image.Point, fX, fY float64, borderType gocv.BorderType) gocv.Mat {
	gocv.GaussianBlur(mat, &mat, ksize, fX, fY, borderType)
	return mat
}

func morphologyEx(mat gocv.Mat, morphType gocv.MorphType, kernel gocv.Mat) gocv.Mat {
	gocv.MorphologyEx(mat, &mat, morphType, kernel)
	return mat
}

func sobel(mat gocv.Mat, ddepth, dx, dy, ksize int, scale, delta float64, borderType gocv.BorderType) gocv.Mat {
	dst := gocv.NewMat()
	gocv.Sobel(mat, &dst, ddepth, dx, dy, ksize, scale, delta, borderType)
	return dst
}

func threshold(mat gocv.Mat, thresh float32, maxvalue float32, threshType gocv.ThresholdType) gocv.Mat {
	gocv.Threshold(mat, &mat, thresh, maxvalue, threshType)
	return mat
}

func crop(src image.Image, rect image.Rectangle) image.Image {
	dst := image.NewRGBA(rect)
	draw.Draw(dst, rect, src, rect.Min, draw.Src)
	return dst
}
