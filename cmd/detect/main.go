package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/danikarik/mrz"
)

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	input := flag.String("input", "", "Input File")
	output := flag.String("output", "", "Output File")

	flag.Parse()

	roi, err := mrz.Detect(*input)
	if err != nil {
		errorExit(err)
	}

	out, err := os.Create(*output)
	if err != nil {
		errorExit(err)
	}

	err = png.Encode(out, roi)
	if err != nil {
		errorExit(err)
	}
}
