package padding

import (
	"errors"
	"image"
	"image/color"
)

type Paddings struct {
	PaddingLeft   int
	PaddingRight  int
	PaddingTop    int
	PaddingBottom int
}

func Padding(img *image.Gray, kernelSize image.Point) (*image.Gray, error) {
	originalSize := img.Bounds().Size()
	p, error := calculatePaddings(kernelSize, image.Point{kernelSize.X / 2, kernelSize.Y / 2})
	if error != nil {
		return nil, error
	}
	rect := getRectangleFromPaddings(p, originalSize)
	padded := image.NewGray(rect)

	for x := p.PaddingLeft; x < originalSize.X+p.PaddingRight; x++ {
		for y := p.PaddingTop; y < originalSize.Y+p.PaddingBottom; y++ {
			padded.Set(x, y, img.GrayAt(x-p.PaddingLeft, y-p.PaddingTop))
		}
	}

	topPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
		padded.Set(x, y, pixel)
	})
	bottomPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
		padded.Set(x, y, pixel)
	})
	leftPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
		padded.Set(x, y, pixel)
	})
	rightPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
		padded.Set(x, y, pixel)
	})

	return padded, nil
}

func calculatePaddings(kernelSize image.Point, anchor image.Point) (Paddings, error) {
	var p Paddings
	if kernelSize.X < 0 || kernelSize.Y < 0 {
		return p, errors.New("Negative size")
	}
	if anchor.X < 0 || anchor.Y < 0 {
		return p, errors.New("Negative anchor value")
	}
	if anchor.X > kernelSize.X || anchor.Y > kernelSize.Y {
		return p, errors.New("Anchor value outside of the kernel")
	}

	p = Paddings{PaddingLeft: anchor.X, PaddingRight: kernelSize.X - anchor.X - 1, PaddingTop: anchor.Y, PaddingBottom: kernelSize.Y - anchor.Y - 1}

	return p, nil
}

func getRectangleFromPaddings(p Paddings, imgSize image.Point) image.Rectangle {
	x := p.PaddingLeft + p.PaddingRight + imgSize.X
	y := p.PaddingTop + p.PaddingBottom + imgSize.Y
	return image.Rect(0, 0, x, y)
}

func topPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		firstPixel := img.At(x-p.PaddingLeft, p.PaddingTop)
		for y := 0; y < p.PaddingTop; y++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func bottomPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		lastPixel := img.At(x-p.PaddingLeft, originalSize.Y-1)
		for y := p.PaddingTop + originalSize.Y; y < originalSize.Y+p.PaddingTop+p.PaddingBottom; y++ {
			setPixel(x, y, lastPixel)
		}
	}
}

func leftPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		firstPixel := padded.At(p.PaddingLeft, y)
		for x := 0; x < p.PaddingLeft; x++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func rightPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		lastPixel := padded.At(originalSize.X+p.PaddingLeft-1, y)
		for x := originalSize.X + p.PaddingLeft; x < originalSize.X+p.PaddingLeft+p.PaddingRight; x++ {
			setPixel(x, y, lastPixel)
		}
	}
}
