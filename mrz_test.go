package mrz_test

import (
	"testing"

	"github.com/danikarik/mrz"
	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	version := mrz.Version()

	if assert.NotEmpty(version) {
		t.Log("gocv version:", version)
	}
}

func TestOpenCVVersion(t *testing.T) {
	assert := assert.New(t)
	version := mrz.OpenCVVersion()

	if assert.NotEmpty(version) {
		t.Log("opencv version:", version)
	}
}

func TestResize(t *testing.T) {
	testCases := []struct {
		Name           string
		Path           string
		Width          float64
		Height         float64
		ExpectedWidth  float64
		ExpectedHeight float64
	}{
		{
			Name:           "Passport01",
			Path:           "testdata/passport_01.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  413,
			ExpectedHeight: 600,
		},
		{
			Name:           "Passport02",
			Path:           "testdata/passport_02.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  404,
			ExpectedHeight: 600,
		},
		{
			Name:           "Passport03",
			Path:           "testdata/passport_03.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  415,
			ExpectedHeight: 600,
		},
		{
			Name:           "Passport04",
			Path:           "testdata/passport_04.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  418,
			ExpectedHeight: 600,
		},
		{
			Name:           "Passport05",
			Path:           "testdata/passport_05.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  432,
			ExpectedHeight: 600,
		},
		{
			Name:           "Passport06",
			Path:           "testdata/passport_06.jpg",
			Width:          0,
			Height:         600,
			ExpectedWidth:  434,
			ExpectedHeight: 600,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			mat := gocv.IMRead(tc.Path, gocv.IMReadUnchanged)
			mat = mrz.Resize(mat, tc.Width, tc.Height, gocv.InterpolationArea)

			assert.Equal(int(tc.ExpectedWidth), mat.Cols())
			assert.Equal(int(tc.ExpectedHeight), mat.Rows())
		})
	}
}
