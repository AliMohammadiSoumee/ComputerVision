package convolution

import (
	"fmt"
	"image"
	"image/color"

	"github.com/alidadar7676/ComputerVision/padding"
	"github.com/alidadar7676/ComputerVision/utils"
)

func ConvolveGray(img *image.Gray, kernel *Kernel) (*image.Gray, error) {
	originalSize := img.Bounds().Size()
	resultImage := image.NewGray(img.Bounds())
	kernelSize := kernel.Size()

	padded, error := padding.Padding(img, kernelSize)
	if error != nil {
		return nil, error
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
		//fmt.Println(x, y, sum)
		sum = utils.Clamp(sum, utils.MinUint8, float64(utils.MaxUint8))
		resultImage.Set(x, y, color.Gray{uint8(sum)})
	})
	fmt.Println(resultImage.Bounds().Size())
	return resultImage, nil
}
