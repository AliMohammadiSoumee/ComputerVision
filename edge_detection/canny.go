package edge_detection

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/alidadar7676/ComputerVision/blurring"
	"github.com/alidadar7676/ComputerVision/utils"
)

func CannyGray(img *image.Gray, kernelSize uint) (*image.Gray, error) {
	blurred, err := blurring.GaussianBlurGray(img, float64(kernelSize), 1)
	if err != nil {
		return nil, err
	}

	vertical, err := VerticalSobelGray(blurred)
	if err != nil {
		return nil, err
	}
	horizontal, err := HorizontalSobelGray(blurred)
	if err != nil {
		return nil, err
	}

	g, theta, err := gradientAndOrientation(vertical, horizontal)
	if err != nil {
		return nil, err
	}

	image := image.NewGray(blurred.Rect)
	for i := 0; i < blurred.Bounds().Size().X; i++ {
		for j := 0; j < blurred.Bounds().Size().Y; j++ {
			image.SetGray(i, j, color.Gray{Y: uint8(g[i][j])})
		}
	}

	thinEdges := nonMaxSuppression(blurred, g, theta)

	hist := doubleThreshold(thinEdges, g, 100)
	return hist, nil

}

func gradientAndOrientation(vertical *image.Gray, horizontal *image.Gray) ([][]float64, [][]float64, error) {
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

func isBiggerThenNeighbours(val float64, neighbour1 float64, neighbour2 float64) bool {
	return val > neighbour1 && val > neighbour2
}

func nonMaxSuppression(img *image.Gray, g [][]float64, theta [][]float64) *image.Gray {
	size := img.Bounds().Size()
	thinEdges := image.NewGray(image.Rect(0, 0, size.X, size.Y))
	utils.ForEachPixel(size, func(x, y int) {
		isLocalMax := false
		if x > 0 && x < size.X-1 && y > 0 && y < size.Y-1 {
			switch theta[x][y] {
			case 45:
				if isBiggerThenNeighbours(g[x][y], g[x+1][y-1], g[x-1][y+1]) {
					isLocalMax = true
				}
			case 90:
				if isBiggerThenNeighbours(g[x][y], g[x+1][y], g[x-1][y]) {
					isLocalMax = true
				}
			case 135:
				if isBiggerThenNeighbours(g[x][y], g[x-1][y-1], g[x+1][y+1]) {
					isLocalMax = true
				}
			case 0:
				if isBiggerThenNeighbours(g[x][y], g[x][y+1], g[x][y-1]) {
					isLocalMax = true
				}
			}
		}
		if isLocalMax {
			thinEdges.SetGray(x, y, color.Gray{Y: utils.MaxUint8})
		}
	})
	return thinEdges
}

func checkNeighbours(x, y int, img *image.Gray) bool {
	return img.GrayAt(x-1, y-1).Y == utils.MaxUint8 || img.GrayAt(x-1, y).Y == utils.MaxUint8 ||
		img.GrayAt(x-1, y+1).Y == utils.MaxUint8 || img.GrayAt(x, y-1).Y == utils.MaxUint8 ||
		img.GrayAt(x, y+1).Y == utils.MaxUint8 || img.GrayAt(x+1, y-1).Y == utils.MaxUint8 ||
		img.GrayAt(x+1, y).Y == utils.MaxUint8 || img.GrayAt(x+1, y+1).Y == utils.MaxUint8
}

func doubleThreshold(img *image.Gray, g [][]float64, upperBound float64) *image.Gray {
	size := img.Bounds().Size()
	res := image.NewGray(img.Rect)
	utils.ForEachPixel(size, func(x int, y int) {
		p := img.GrayAt(x, y)
		if p.Y == utils.MaxUint8 {
			if g[x][y] > upperBound {
				res.SetGray(x, y, color.Gray{Y: utils.MaxUint8})
			} else {
				res.SetGray(x, y, color.Gray{Y: utils.MinUint8})
			}
		}
	})
	return res
}
