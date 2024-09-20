package sprites

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

var (
	//go:embed assets/VeryPixelatedRoundedBlueWhite.png
	veryPixelatedRoundedBlueWhite []byte
)

type BorderSprites struct {
	// the default border. It is
	VeryPixelatedRoundedBlueWhite ReferenceBorder
}

type ReferenceBorder struct {
	border                    *ebiten.Image
	referenceBorderDimensions *ReferenceBorderDimensions
}

type ReferenceBorderDimensions struct {
	topLeft, bottomLeft, topRight, bottomRight image.Rectangle
	horizCopy, vertCopy                        image.Rectangle
}

func NewBorderSprites() *BorderSprites {
	borderSprites := &BorderSprites{}
	borderSprites.VeryPixelatedRoundedBlueWhite.border = NewPngSprite(veryPixelatedRoundedBlueWhite)

	return borderSprites
}

// CreateSizedAndScaledBorderSprite
func (sb *ReferenceBorder) CreateSizedAndScaledBorderSprite(placement PercentBasedPlacement) (*ebiten.Image, *ebiten.DrawImageOptions) {

	op := &ebiten.DrawImageOptions{}

	// get screen sizes
	screenWidth, screenHeight := ebiten.WindowSize()

	// get the x start and end values based on the percent
	var xLeft = float64(screenWidth) * placement.StartPercentX
	var xRight = float64(screenWidth) * placement.EndPercentX
	var yTop = float64(screenHeight) * placement.StartPercentY
	var yBottom = float64(screenHeight) * placement.EndPercentY

	targetWidth := xRight - xLeft
	targetHeight := yBottom - yTop

	originalImgWidth := sb.border.Bounds().Dx()
	originalImgHeight := sb.border.Bounds().Dy()

	// we scale to X value, and match the ratio on Y so the image doesn't stretch and distort
	scaleX := targetWidth / float64(originalImgWidth)
	scaleY := targetHeight / float64(originalImgHeight)
	op.GeoM.Scale(scaleX, scaleY)

	scaledWidth := scaleX * float64(originalImgWidth)

	topLeftX := (float64(screenWidth) - scaledWidth) / 2
	op.GeoM.Translate(topLeftX, yTop)

	// scale the border to screen scale

	// find the corners
	// scale the corners

	// fill in horiz and vert
	return sb.border, op
}

func GetCornersOfReferenceBorder(width, height int) *ReferenceBorderDimensions {
	referenceBorderDimensions := ReferenceBorderDimensions{}
	//totalWidth := border.Bounds().Dx()
	//totalHeight := border.Bounds().Dy()
	// if even, pick one
	// if odd, pick middle
	xMiddle := width / 2
	yMiddle := height / 2

	referenceBorderDimensions.topLeft = image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{X: xMiddle - 1, Y: yMiddle - 1},
	}
	referenceBorderDimensions.topRight = image.Rectangle{
		Min: image.Point{X: xMiddle + 1, Y: 0},
		Max: image.Point{X: width - 1, Y: yMiddle - 1},
	}
	referenceBorderDimensions.bottomLeft = image.Rectangle{
		Min: image.Point{X: 0, Y: yMiddle + 1},
		Max: image.Point{X: xMiddle - 1, Y: height - 1},
	}
	referenceBorderDimensions.bottomRight = image.Rectangle{
		Min: image.Point{X: xMiddle + 1, Y: yMiddle + 1},
		Max: image.Point{X: width - 1, Y: height - 1},
	}
	return &referenceBorderDimensions
}
