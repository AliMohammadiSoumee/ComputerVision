package utils

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/coraldane/resize"
)

func ForEachPixel(size image.Point, f func(x int, y int)) {
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			f(i, j)
		}
	}
}
func Clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func AddGrayWeighted(img1 *image.Gray, w1 float64, img2 *image.Gray, w2 float64) (*image.Gray, error) {
	size1 := img1.Bounds().Size()
	size2 := img2.Bounds().Size()
	if size1.X != size2.X || size1.Y != size2.Y {
		return nil, errors.New("The size of the two image does not match")
	}
	res := image.NewGray(img1.Bounds())
	ForEachPixel(size1, func(x int, y int) {
		p1 := img1.GrayAt(x, y)
		p2 := img2.GrayAt(x, y)
		sum := Clamp(float64(p1.Y)*w1+float64(p2.Y)*w2, MinUint8, float64(MaxUint8))
		res.SetGray(x, y, color.Gray{uint8(sum)})
	})
	return res, nil
}

func GrayScale(img image.Image) *image.Gray {
	gray := image.NewGray(img.Bounds())
	size := img.Bounds().Size()
	ForEachPixel(size, func(x, y int) {
		gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
	})
	return gray
}

func SubtractGrayImages(firstImg, secondImg image.Image) *image.Gray {
	if firstImg.Bounds() != secondImg.Bounds() {
		return image.NewGray(image.Rectangle{})
	}
	resImg := image.NewGray(firstImg.Bounds())

	ForEachPixel(firstImg.Bounds().Size(), func(x, y int) {
		resImg.Set(x, y, SubtractGrayColor(firstImg.At(x, y), secondImg.At(x, y)))
	})

	return resImg
}

func SubtractGrayColor(firstCol, secondCol color.Color) color.Color {
	_, firstGray, _, _ := firstCol.RGBA()
	_, secondGray, _, _ := secondCol.RGBA()

	sub := firstGray - secondGray
	return color.Gray{Y: uint8(sub)}
}

func Create3DDirection() (dirx, diry, dirz []int) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				dirx = append(dirx, x)
				diry = append(diry, y)
				dirz = append(dirz, z)
			}
		}
	}
	return
}

func HalveImage(srcImg *image.Gray) *image.Gray {
	bounds := srcImg.Bounds()
	scaledImg := resize.Resize(bounds.Dx()/2, bounds.Dy()/2, srcImg, resize.Bilinear)
	resImg := GrayScale(scaledImg)
	return resImg
}

func DiscreteOrientation(x float64) (float64, error) {
	angle := 180 * x / math.Pi
	if IsBetween(angle, 0, 22.5) || IsBetween(angle, -180, -157.5) {
		return 0, nil
	}
	if IsBetween(angle, 157.5, 180) || IsBetween(angle, -22.5, 0) {
		return 0, nil
	}
	if IsBetween(angle, 22.5, 67.5) || IsBetween(angle, -157.5, -112.5) {
		return 45, nil
	}
	if IsBetween(angle, 67.5, 112.5) || IsBetween(angle, -112.5, -67.5) {
		return 90, nil
	}
	if IsBetween(angle, 112.5, 157.5) || IsBetween(angle, -67.5, -22.5) {
		return 135, nil
	}
	return 0, errors.New("Invalid angle")
}

func IsBetween(val float64, lowerBound float64, upperBound float64) bool {
	return val >= lowerBound && val < upperBound
}

func CreateGrayImage(mat [][]float64, rect image.Rectangle) *image.Gray {
	image := image.NewGray(rect)
	ForEachPixel(rect.Size(), func(x, y int) {
		pix := Clamp(mat[x][y], MinUint8, float64(MaxUint8))
		image.Set(x, y, color.Gray{uint8(pix)})

	})
	return image
}
