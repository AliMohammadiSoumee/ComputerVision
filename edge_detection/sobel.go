package edge_detection

import (
	"image"

	"github.com/alidadar7676/ComputerVision/convolution"
	"github.com/alidadar7676/ComputerVision/utils"
)

var horizontalKernel = convolution.Kernel{Content: [][]float64{
	{-1, 0, 1},
	{-2, 0, 2},
	{-1, 0, 1},
}, Width: 3, Height: 3}

var verticalKernel = convolution.Kernel{Content: [][]float64{
	{-1, -2, -1},
	{0, 0, 0},
	{1, 2, 1},
}, Width: 3, Height: 3}

func HorizontalSobelGray(gray *image.Gray) (*image.Gray, error) {
	return convolution.ConvolveGray(gray, &horizontalKernel)
}

func VerticalSobelGray(gray *image.Gray) (*image.Gray, error) {
	return convolution.ConvolveGray(gray, &verticalKernel)
}

func SobelGray(img *image.Gray) (*image.Gray, error) {
	horizontal, error := HorizontalSobelGray(img)
	if error != nil {
		return nil, error
	}

	vertical, error := VerticalSobelGray(img)
	if error != nil {
		return nil, error
	}

	res, error := utils.AddGrayWeighted(horizontal, 0.5, vertical, 0.5)
	if error != nil {
		return nil, error
	}
	return res, nil
}
