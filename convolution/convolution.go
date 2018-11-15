package convolution

import (
	"image"

	"github.com/alidadar7676/ComputerVision/padding"
	"github.com/alidadar7676/ComputerVision/utils"
)

func ConvolveGray(img *image.Gray, kernel *Kernel) ([][]float64, error) {
	originalSize := img.Bounds().Size()
	kernelSize := kernel.Size()

	result := make([][]float64, originalSize.X)
	for i := 0; i < originalSize.X; i++ {
		result[i] = make([]float64, originalSize.Y)
	}

	padded, err := padding.Padding(img, kernelSize)
	if err != nil {
		return nil, err
	}

	utils.ForEachPixel(originalSize, func(x int, y int) {
		sum := float64(0)
		for ky := 0; ky < kernelSize.Y; ky++ {
			for kx := 0; kx < kernelSize.X; kx++ {
				pixel := padded.GrayAt(x+kx, y+ky)
				kE := kernel.At(kx, ky)
				sum += float64(pixel.Y) * kE
			}
		}
		result[x][y] = sum
	})
	return result, nil
}
