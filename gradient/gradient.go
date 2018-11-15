package gradient

import (
	"image"
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

func Horizontal(gray *image.Gray) ([][]float64, error) {
	if mat, err := convolution.ConvolveGray(gray, &horizontalKernel); err != nil {
		return nil, err
	} else {
		return mat, nil
	}

}

func Vertical(gray *image.Gray) ([][]float64, error) {
	if mat, err := convolution.ConvolveGray(gray, &verticalKernel); err != nil {
		return nil, err
	} else {
		return mat, nil
	}
}

func GradientAndOrientation(bound image.Point, vertical, horizontal [][]float64) ([][]float64, [][]float64) {
	size := bound
	theta := make([][]float64, size.X)
	g := make([][]float64, size.X)
	for x := 0; x < size.X; x++ {
		theta[x] = make([]float64, size.Y)
		g[x] = make([]float64, size.Y)
		for y := 0; y < size.Y; y++ {
			px := vertical[x][y]   //float64(vertical.GrayAt(x, y).y)
			py := horizontal[x][y] //float64(horizontal.GrayAt(x, y).y)
			g[x][y] = math.Hypot(px, py)
			theta[x][y] = math.Atan2(py, px) //math.Atan2(float64(vertical.GrayAt(x, y).y), float64(horizontal.GrayAt(x, y).y))
		}
	}
	return g, theta
}
