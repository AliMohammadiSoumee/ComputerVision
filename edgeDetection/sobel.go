package edgeDetection

import (
	"image"

	"github.com/alidadar7676/ComputerVision/utils"
	"github.com/alidadar7676/ComputerVision/gradient"
)

func SobelGray(img *image.Gray) (*image.Gray, error) {
	horizontal, err := gradient.HorizontalSobelGray(img)
	if err != nil {
		return nil, err
	}

	vertical, err := gradient.VerticalSobelGray(img)
	if err != nil {
		return nil, err
	}

	res, err := utils.AddGrayWeighted(horizontal, 0.5, vertical, 0.5)
	if err != nil {
		return nil, err
	}
	return res, nil
}
