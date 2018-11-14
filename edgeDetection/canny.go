package edgeDetection

import (
	"image"
	"image/color"
	"github.com/alidadar7676/ComputerVision/blurring"
	"github.com/alidadar7676/ComputerVision/utils"
	"github.com/alidadar7676/ComputerVision/gradient"
)

func CannyGray(img *image.Gray, kernelSize uint) (*image.Gray, error) {
	blurred, err := blurring.GaussianBlurGray(img, float64(kernelSize), 1)
	if err != nil {
		return nil, err
	}

	vertical, err := gradient.VerticalSobelGray(blurred)
	if err != nil {
		return nil, err
	}
	horizontal, err := gradient.HorizontalSobelGray(blurred)
	if err != nil {
		return nil, err
	}

	g, theta, err := gradient.GradientAndOrientation(vertical, horizontal)
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
