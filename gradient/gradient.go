package gradient

import (
	"image"
	"errors"
	"math"
	"github.com/alidadar7676/ComputerVision/convolution"
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

func GradientAndOrientation(vertical *image.Gray, horizontal *image.Gray) ([][]float64, [][]float64, error) {
	size := vertical.Bounds().Size()
	theta := make([][]float64, size.X)
	g := make([][]float64, size.X)
	for x := 0; x < size.X; x++ {
		theta[x] = make([]float64, size.Y)
		g[x] = make([]float64, size.Y)
		err := errors.New("none")
		for y := 0; y < size.Y; y++ {
			px := float64(vertical.GrayAt(x, y).Y)
			py := float64(horizontal.GrayAt(x, y).Y)
			g[x][y] = math.Hypot(px, py)
			theta[x][y], err = orientation(math.Atan2(float64(vertical.GrayAt(x, y).Y), float64(horizontal.GrayAt(x, y).Y)))
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return g, theta, nil
}

func isBetween(val float64, lowerBound float64, upperBound float64) bool {
	return val >= lowerBound && val < upperBound
}

func orientation(x float64) (float64, error) {
	angle := 180 * x / math.Pi
	if isBetween(angle, 0, 22.5) || isBetween(angle, -180, -157.5) {
		return 0, nil
	}
	if isBetween(angle, 157.5, 180) || isBetween(angle, -22.5, 0) {
		return 0, nil
	}
	if isBetween(angle, 22.5, 67.5) || isBetween(angle, -157.5, -112.5) {
		return 45, nil
	}
	if isBetween(angle, 67.5, 112.5) || isBetween(angle, -112.5, -67.5) {
		return 90, nil
	}
	if isBetween(angle, 112.5, 157.5) || isBetween(angle, -67.5, -22.5) {
		return 135, nil
	}
	return 0, errors.New("Invalid angle")
}
