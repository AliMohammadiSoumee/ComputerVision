package edgeDetection

import (
	"image"

	"github.com/alidadar7676/ComputerVision/gradient"
	"github.com/alidadar7676/ComputerVision/utils"
)

func SobelGray(img *image.Gray) (*image.Gray, error) {
	hor, err := gradient.Horizontal(img)
	if err != nil {
		return nil, err
	}

	ver, err := gradient.Vertical(img)
	if err != nil {
		return nil, err
	}
	vertical := utils.CreateGrayImage(ver, img.Rect)
	horizontal := utils.CreateGrayImage(hor, img.Rect)

	res, err := utils.AddGrayWeighted(horizontal, 0.5, vertical, 0.5)
	if err != nil {
		return nil, err
	}
	return res, nil
}
