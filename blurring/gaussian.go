package blurring

import (
	"errors"
	"image"
	"math"

	"github.com/alidadar7676/ComputerVision/convolution"
)

func GaussianBlurGray(img *image.Gray, radius float64, sigma float64) (*image.Gray, error) {
	if radius <= 0 {
		return nil, errors.New("radius must be bigger then 0")
	}
	return convolution.ConvolveGray(img, generateGaussianKernel(radius, sigma).Normalize())
}

func generateGaussianKernel(radius float64, sigma float64) *convolution.Kernel {
	length := int(math.Ceil(2*radius + 1))
	kernel, _ := convolution.NewKernel(length, length)
	for x := 0; x < length; x++ {
		for y := 0; y < length; y++ {
			kernel.Set(x, y, gaussianFunc(float64(x)-radius, float64(y)-radius, sigma))
		}
	}
	return kernel
}

func gaussianFunc(x, y, sigma float64) float64 {
	sigSqr := sigma * sigma
	return (1.0 / (2 * math.Pi * sigSqr)) * math.Exp(-(x*x + y*y)/(2*sigSqr))
}
