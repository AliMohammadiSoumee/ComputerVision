package sift

import (
	"image"
	"math"

	"github.com/alidadar7676/ComputerVision/blurring"
	"github.com/alidadar7676/ComputerVision/gradient"
	"github.com/alidadar7676/ComputerVision/utils"
)

type KeyPoint struct {
	orientation float64
	octave      int
	scale       int
	x           int
	y           int
	Feature     []float64
}

func SiftFeatures(img *image.Gray, oct int, scale int, tresh float64) []KeyPoint {
	scaleSpace := createScaleSpace(img, oct, scale)
	dog := createDoG(oct, scale, scaleSpace)
	candidate := extractKeyPoints(oct, scale, dog)
	gradSpec := createGradientSpec(oct, scale, dog)
	result := filterKeyPoints(candidate, gradSpec, dog, tresh)
	fetchMaxOrientation(result, gradSpec)
	createFeatures(result, gradSpec)

	return result
}

func createScaleSpace(img *image.Gray, octave, scale int) [][]*image.Gray {
	k := math.Pow(2.0, 1.0/float64(scale))
	sig := make([]float64, scale)

	sig[0] = 1
	for i := 1; i < scale; i++ {
		sig[i] = sig[i-1] * k
	}

	scaleSpace := make([][]*image.Gray, octave)
	for row := 0; row < octave; row++ {
		scaleSpace[row] = make([]*image.Gray, scale)

		for col := 0; col < scale; col++ {
			if row == 0 && col == 0 {
				scaleSpace[row][col] = img
			} else if col == 0 {
				scaleSpace[row][col] = utils.HalveImage(scaleSpace[row-1][col])
			} else {
				tmp, err := blurring.GaussianBlurGray(scaleSpace[row][col-1], 5, sig[col])
				if err != nil {
					panic("Can not blur the image")
				}
				scaleSpace[row][col] = tmp
			}
		}
	}
	return scaleSpace
}

func createDoG(octave, scale int, scaleSpace [][]*image.Gray) [][]*image.Gray {
	dog := make([][]*image.Gray, octave)
	for row := 0; row < octave; row++ {
		dog[row] = make([]*image.Gray, scale-1)

		for col := 0; col < scale-1; col++ {
			dog[row][col] = utils.SubtractGrayImages(scaleSpace[row][col], scaleSpace[row][col+1])
		}
	}

	return dog
}

func extractKeyPoints(octave, scale int, dog [][]*image.Gray) []KeyPoint {
	dirx, diry, dirz := utils.Create3DDirection()

	localMinimum := func(row, col, x, y int) bool {
		pix := dog[row][col].GrayAt(x, y).Y
		for d := 0; d < 27; d++ {
			if dirx[d] == 0 && diry[d] == 0 && dirz[d] == 0 {
				continue
			}
			if dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y <= pix {
				return false
			}
		}
		return true
	}

	localMaximum := func(row, col, x, y int) bool {
		pix := dog[row][col].GrayAt(x, y).Y
		for d := 0; d < 27; d++ {
			if dirx[d] == 0 && diry[d] == 0 && dirz[d] == 0 {
				continue
			}
			if dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y >= pix {
				return false
			}
		}
		return true
	}

	candidateKeys := make([]KeyPoint, 0)

	for row := 0; row < octave; row++ {
		for col := 1; col < scale-2; col++ {
			bound := dog[row][col].Bounds().Size()

			utils.ForEachPixel(bound, func(x, y int) {
				if x == 0 || y == 0 || x+1 >= dog[row][col].Bounds().Size().X || y+1 >= dog[row][col].Bounds().Size().Y {
					return
				}
				if localMinimum(row, col, x, y) || localMaximum(row, col, x, y) {
					candidateKeys = append(candidateKeys, KeyPoint{
						x:      x,
						y:      y,
						octave: row,
						scale:  col,
					})
				}
			})
		}
	}
	return candidateKeys
}

type gradientSpec struct {
	g     [][]float64
	theta [][]float64
	gX    [][]float64
	gY    [][]float64
}

func (grad *gradientSpec) normalize() {
	max := 0.0
	for _, row := range grad.g {
		for _, val := range row {
			max = math.Max(max, val)
		}
	}

	for row := range grad.g {
		for col := range grad.g[row] {
			grad.g[row][col] /= max
		}
	}
}

func createGradientSpec(octave, scale int, dog [][]*image.Gray) [][]gradientSpec {
	gradSpec := make([][]gradientSpec, octave)

	for i := 0; i < octave; i++ {
		gradSpec[i] = make([]gradientSpec, scale)
	}

	for row := 0; row < octave; row++ {
		for col := 1; col < scale-2; col++ {
			gs := &gradSpec[row][col]

			if gX, err := gradient.Horizontal(dog[row][col]); err != nil {
				panic("Cannot create gX")
			} else if gY, err := gradient.Vertical(dog[row][col]); err != nil {
				panic("Cannot create gX")
			} else {
				gs.g, gs.theta = gradient.GradientAndOrientation(dog[row][col].Bounds().Size(), gX, gY)
				gs.normalize()
			}
		}
	}

	return gradSpec
}

func filterKeyPoints(candidate []KeyPoint, gradSpec [][]gradientSpec, dog [][]*image.Gray, tresh float64) []KeyPoint {
	result := make([]KeyPoint, 0)

	for _, can := range candidate {
		bound := dog[can.octave][can.scale].Bounds()

		if can.x < 8 || can.x > bound.Size().X-10 || can.y < 8 || can.y > bound.Size().Y-10 {
			continue
		}

		if gradSpec[can.octave][can.scale].g[can.x][can.y] > tresh {
			result = append(result, can)
		}
	}
	return result
}

func fetchMaxOrientation(keys []KeyPoint, gradSpec [][]gradientSpec) {
	for ind, key := range keys {
		keys[ind].orientation = getBiggestOrientation(gradSpec[key.octave][key.scale].theta, key.x, key.y)
	}
}

func getBiggestOrientation(theta [][]float64, x, y int) float64 {
	max := -1000.0
	for row := -8; row < 8; row++ {
		for col := -8; col < 8; col++ {
			max = math.Max(max, theta[x+row][y+col])
		}
	}
	return max
}

func createFeatures(keys []KeyPoint, gradSpec [][]gradientSpec) {
	for ind, key := range keys {
		gs := &gradSpec[key.octave][key.scale]

		for row := -8; row < 8; row += 4 {
			for col := -8; col < 8; col += 4 {
				vector8 := createSub4x4Features(gs, key.x+row, key.y+col, key.orientation)
				keys[ind].Feature = append(keys[ind].Feature, vector8...)
			}
		}
	}
}

func createSub4x4Features(gs *gradientSpec, x, y int, orien float64) []float64 {
	vector := make([]float64, 8)
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			index := convertAngle(gs.theta[x+row][y+col] - orien)
			vector[index] += gs.g[x][y]
		}
	}
	return vector
}

func convertAngle(ang float64) int {
	angle := 180.0 * ang / math.Pi
	if angle < 0 {
		angle += 360
	}

	switch {
	case utils.IsBetween(angle, 0, 45):
		return 0
	case utils.IsBetween(angle, 45, 90):
		return 1
	case utils.IsBetween(angle, 90, 135):
		return 2
	case utils.IsBetween(angle, 135, 180):
		return 3
	case utils.IsBetween(angle, 180, 225):
		return 4
	case utils.IsBetween(angle, 225, 270):
		return 5
	case utils.IsBetween(angle, 270, 315):
		return 6
	case utils.IsBetween(angle, 315, 360):
		return 7
	}
	return 4
}
