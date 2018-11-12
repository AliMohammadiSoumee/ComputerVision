package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/alidadar7676/ComputerVision/sift"
	"github.com/alidadar7676/ComputerVision/utils"
)

func main() {
	existingImageFile, err := os.Open(os.Args[1])
	if err != nil {
		panic("Can not find input file")
	}
	defer existingImageFile.Close()

	existingImageFile.Seek(0, 0)

	image, err := jpeg.Decode(existingImageFile)
	if err != nil {
		panic("Can not decode image file. The format of image must be .jpeg")
	}

	grayImage := utils.GrayScale(image)

	s := sift.Sift{}
	sx, sy := s.SiftFeatures(grayImage, 4, 4)
	fmt.Println(sx, sy)

	/*
	outfile, err := os.Create(os.Args[2])
	if err != nil {
		panic("Can not find output file")
	}
	defer outfile.Close()

	jpeg.Encode(outfile, sobelImage, nil)
	*/
}
