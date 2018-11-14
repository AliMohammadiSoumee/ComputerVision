package sift

import (
	"image"
	"math"
	"github.com/alidadar7676/ComputerVision/utils"
	"github.com/alidadar7676/ComputerVision/blurring"
	"fmt"
	"github.com/alidadar7676/ComputerVision/gradient"
)

type KeyPoint struct {
	octave int
	scale  int
	x      int
	y      int
}

func SiftFeatures(img *image.Gray, oct int, scale int, tresh float64) []KeyPoint {
	scaleSpace := createScaleSpace(img, oct, scale)
	dog := createDoG(oct, scale, scaleSpace)
	candidate := extractKeyPoints(oct, scale, dog)
	gradSpec := createGradientSpec(oct, scale, dog)
	fmt.Println("Len of Candidate keyPoints: ", len(candidate))
	result := filterKeyPoints(candidate, gradSpec, tresh)

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
			if dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y < pix {
				return false
			}
		}
		return true
	}

	localMaximum := func(row, col, x, y int) bool {
		pix := dog[row][col].GrayAt(x, y).Y
		for d := 0; d < 27; d++ {
			if dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y > pix {
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

type gradianSpec struct {
	g     [][]float64
	theta [][]float64
}

func (grad *gradianSpec) normalize() {
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

func createGradientSpec(octave, scale int, dog [][]*image.Gray) [][]gradianSpec {
	gX := make([][]*image.Gray, octave)
	gY := make([][]*image.Gray, octave)
	gradSpec := make([][]gradianSpec, octave)

	for i := 0; i < octave; i++ {
		gX[i] = make([]*image.Gray, scale)
		gY[i] = make([]*image.Gray, scale)
		gradSpec[i] = make([]gradianSpec, scale)
	}

	var err error
	for row := 0; row < octave; row++ {
		for col := 1; col < scale-2; col++ {

			if gX[row][col], err = gradient.HorizontalSobelGray(dog[row][col]); err != nil {
				panic("Cannot create gX")
			}
			if gY[row][col], err = gradient.VerticalSobelGray(dog[row][col]); err != nil {
				panic("Cannot create gX")
			}
			if g, th, err := gradient.GradientAndOrientation(gX[row][col], gY[row][col]); err != nil {
				panic("Cannot create g and theta")
			} else {
				gradSpec[row][col].g = g
				gradSpec[row][col].theta = th
				gradSpec[row][col].normalize()
			}
		}
	}

	return gradSpec
}

func filterKeyPoints(candidate []KeyPoint, gradSpec [][]gradianSpec, tresh float64) []KeyPoint {
	result := make([]KeyPoint, 0)

	for _, can := range candidate {
		if gradSpec[can.octave][can.scale].g[can.x][can.y] > tresh {
			result = append(result, can)
		}
	}
	return result
}
