package sift

import (
	"image"
	"math"
	"github.com/alidadar7676/ComputerVision/utils"
	"github.com/coraldane/resize"
	"github.com/alidadar7676/ComputerVision/blurring"
	"fmt"
)

type Sift struct {
	img        image.Gray
	scaleSpace [][]image.Gray
	octave     int
	scaleLvl   int
	dog        [][]image.Gray
}

func (s *Sift) SiftFeatures(img *image.Gray, oct int, scale int) ([]int, []int) {
	s.img = *img
	s.octave = oct
	s.scaleLvl = scale

	s.createScaleSpace()
	s.createDoG()
	return s.extractKeyPoints()
}

func (s *Sift) createScaleSpace() {
	k := math.Pow(2.0, 1.0/float64(s.scaleLvl))
	sig := make([]float64, s.scaleLvl)

	sig[0] = 1
	for i := 1; i < s.scaleLvl; i++ {
		sig[i] = sig[i-1] * k
	}

	s.scaleSpace = make([][]image.Gray, s.octave)
	for row := 0; row < s.octave; row++ {
		s.scaleSpace[row] = make([]image.Gray, s.scaleLvl)

		for col := 0; col < s.scaleLvl; col++ {
			if row == 0 && col == 0 {
				s.scaleSpace[row][col] = s.img
			} else if col == 0 {
				s.scaleSpace[row][col] = halveImage(&s.scaleSpace[row-1][col])
			} else {
				tmp, err := blurring.GaussianBlurGray(&s.scaleSpace[row][col-1], 5, sig[col])
				if err != nil {
					panic("Can not blur the image")
				}
				s.scaleSpace[row][col] = *tmp
			}
		}
	}
}

func halveImage(srcImg image.Image) image.Gray {
	bounds := srcImg.Bounds()
	scaledImg := resize.Resize(bounds.Dx()/2, bounds.Dy()/2, srcImg, resize.Bilinear)
	resImg := utils.GrayScale(scaledImg)
	return *resImg
}

func (s *Sift) createDoG() {
	s.dog = make([][]image.Gray, s.octave)
	for row := 0; row < s.octave; row++ {
		s.dog[row] = make([]image.Gray, s.scaleLvl-1)

		for col := 0; col < s.scaleLvl-1; col++ {
			s.dog[row][col] = *utils.SubtractGrayImages(&s.scaleSpace[row][col], &s.scaleSpace[row][col+1])
		}
	}
}

func (s *Sift) extractKeyPoints() ([]int, []int) {
	keysX := []int{}
	keysY := []int{}
	dirx, diry, dirz := utils.Create3DDirection()

	localMinimum := func(row, col, x, y int) bool {
		pix := s.dog[row][col].GrayAt(x, y).Y
		for d := 0; d < 27; d++ {
			if s.dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y < pix {
				return false
			}
		}
		return true
	}

	localMaximum := func(row, col, x, y int) bool {
		pix := s.dog[row][col].GrayAt(x, y).Y
		for d := 0; d < 27; d++ {
			if s.dog[row][col+dirz[d]].GrayAt(x+dirx[d], y+diry[d]).Y > pix {
				return false
			}
		}
		return true
	}

	for row := 0; row < s.octave; row++ {
		for col := 1; col < s.scaleLvl-2; col++ {
			counter := 0

			bound := s.dog[row][col].Bounds().Size()
			utils.ForEachPixel(bound, func(x, y int) {
				if x == 0 || y == 0 || x+1 >= s.dog[row][col].Bounds().Size().X || y+1 >= s.dog[row][col].Bounds().Size().Y {
					return
				}
				if localMinimum(row, col, x, y) || localMaximum(row, col, x, y) {
					keysX = append(keysX, x)
					keysY = append(keysY, y)
				} else {
					counter++
				}
			})
			fmt.Println(len(keysX), counter, "---------> ", s.dog[row][col].Bounds().Size().X * s.dog[row][col].Bounds().Size().Y)
		}
	}

	return keysX, keysY
}
