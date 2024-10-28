package sprites

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
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
	referenceBorderDimensions *referenceBorderDimensions
	referenceBorderBits       *borderBits
}

type referenceBorderDimensions struct {
	topLeft, bottomLeft, topRight, bottomRight                 image.Rectangle
	horizCopyLeft, horizCopyRight, vertCopyLeft, vertCopyRight image.Rectangle
}

type borderBits struct {
	topLeft        *ebiten.Image
	bottomRight    *ebiten.Image
	topRight       *ebiten.Image
	bottomLeft     *ebiten.Image
	vertCopyLeft   *ebiten.Image
	horizCopyLeft  *ebiten.Image
	vertCopyRight  *ebiten.Image
	horizCopyRight *ebiten.Image
}

func NewBorderSprites() *BorderSprites {
	borderSprites := &BorderSprites{}
	borderSprites.VeryPixelatedRoundedBlueWhite.border = NewPngSprite(veryPixelatedRoundedBlueWhite)

	return borderSprites
}

func (sb *ReferenceBorder) createBorderImage(idealWidthForScaling, width, height int) *ebiten.Image {
	windowWidth, _ := ebiten.WindowSize()
	idealScaleX := float64(windowWidth) / float64(idealWidthForScaling)
	scaledDimensions := sb.referenceBorderDimensions.createScaledDimensions(idealScaleX, idealScaleX)

	img := ebiten.NewImage(width, height)

	//cornerHeight := scaledDimensions.topLeft.Max.Y
	bothCornerHeight := scaledDimensions.topLeft.Bounds().Dy() + scaledDimensions.bottomRight.Bounds().Dy()
	bothCornerWidth := scaledDimensions.topLeft.Bounds().Dx() + scaledDimensions.bottomRight.Bounds().Dx()

	/// START WITH SIDE AND TOP BORDERS, THEY WILL GET OVERWRITTEN IF NASTY
	// now fill in the blanks with the horizontal and vertical bits
	// sprite is horizontal
	// goes from LEFT SIDE, top to bottom
	horizScaledOpLeft := ebiten.GeoM{}
	horizScaledOpLeft.Scale(idealScaleX, 1)
	horizScaledOpLeft.Translate(
		0,
		float64(scaledDimensions.topRight.Bounds().Dy()))
	for i := 0; i < height-(bothCornerHeight); i++ {
		img.DrawImage(sb.referenceBorderBits.horizCopyLeft, &ebiten.DrawImageOptions{GeoM: horizScaledOpLeft})
		horizScaledOpLeft.Translate(0, 1)
	}

	// RIGHT HORIZONTAL
	// sprite is horizontal
	// goes from RIGHT, top to bottom
	horizScaledOpRight := ebiten.GeoM{}
	horizScaledOpRight.Scale(idealScaleX, 1)
	horizScaledOpRight.Translate(
		float64(width-scaledDimensions.horizCopyRight.Dx()),
		float64(scaledDimensions.topRight.Dy()))
	for i := 0; i < height-(bothCornerHeight); i++ {
		img.DrawImage(sb.referenceBorderBits.horizCopyRight, &ebiten.DrawImageOptions{GeoM: horizScaledOpRight})
		horizScaledOpRight.Translate(0, 1)
	}

	// LEFT VERTICAL
	// sprite is vertical
	// goes from TOP, left to right
	vertScaledOpLeft := ebiten.GeoM{}
	vertScaledOpLeft.Scale(1, idealScaleX)
	vertScaledOpLeft.Translate(float64(scaledDimensions.topLeft.Dx()), 0)
	for i := 0; i <= width-(bothCornerWidth); i++ {
		img.DrawImage(sb.referenceBorderBits.vertCopyLeft, &ebiten.DrawImageOptions{GeoM: vertScaledOpLeft})
		vertScaledOpLeft.Translate(1, 0)
	}

	// RIGHT VERTICAL
	// sprite is vertical
	// goes from BOTTOM, left to right
	vertScaledOpRight := ebiten.GeoM{}
	vertScaledOpRight.Scale(1, idealScaleX)
	vertScaledOpRight.Translate(
		float64(scaledDimensions.bottomLeft.Dx()),
		float64(height-scaledDimensions.vertCopyRight.Dy()),
	)
	for i := 0; i < width-(bothCornerWidth); i++ {
		img.DrawImage(sb.referenceBorderBits.vertCopyRight, &ebiten.DrawImageOptions{GeoM: vertScaledOpRight})
		vertScaledOpRight.Translate(1, 0)
	}

	/// START CORNERS
	topLeftOp := ebiten.GeoM{}
	topLeftOp.Scale(idealScaleX, idealScaleX)
	img.DrawImage(sb.referenceBorderBits.topLeft, &ebiten.DrawImageOptions{
		GeoM: topLeftOp,
	})

	bottomLeftOp := ebiten.GeoM{}
	bottomLeftOp.Scale(idealScaleX, idealScaleX)
	bottomLeftOp.Translate(0,
		float64(height-scaledDimensions.bottomLeft.Dy()),
	)
	img.DrawImage(sb.referenceBorderBits.bottomLeft, &ebiten.DrawImageOptions{GeoM: bottomLeftOp})

	topRightOp := ebiten.GeoM{}
	topRightOp.Scale(idealScaleX, idealScaleX)
	topRightOp.Translate(
		float64(width-scaledDimensions.topRight.Dx()),
		0,
	)
	img.DrawImage(sb.referenceBorderBits.topRight, &ebiten.DrawImageOptions{GeoM: topRightOp})

	bottomRightOp := ebiten.GeoM{}
	bottomRightOp.Scale(idealScaleX, idealScaleX)
	bottomRightOp.Translate(
		float64(width-scaledDimensions.bottomRight.Dx()),
		float64(height-scaledDimensions.bottomRight.Dy()),
	)
	img.DrawImage(sb.referenceBorderBits.bottomRight, &ebiten.DrawImageOptions{GeoM: bottomRightOp})

	return img

}

func (sb *ReferenceBorder) createBorderBits() *borderBits {
	bits := borderBits{}

	// make our corner bits
	bits.topLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.topLeft))
	bits.bottomLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.bottomLeft))
	bits.topRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.topRight))
	bits.bottomRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.bottomRight))
	bits.horizCopyLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.horizCopyLeft))
	bits.vertCopyLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.vertCopyLeft))
	bits.horizCopyRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.horizCopyRight))
	bits.vertCopyRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.vertCopyRight))

	return &bits
}

func (sb *ReferenceBorder) CreateSizedAndScaledBorderSprite(idealWidthForScaling int, placement PercentBasedPlacement) (*ebiten.Image, *ebiten.DrawImageOptions) {

	// get the corners and copies before we scale it
	sb.referenceBorderDimensions = getCornersOfReferenceBorder(sb.border.Bounds().Dx(), sb.border.Bounds().Dy())

	op := &ebiten.DrawImageOptions{}

	// THIS STUFF IS FOR PLACEMENT AFTER WE HAVE THE IMAGE
	// get screen sizes
	screenWidth, screenHeight := ebiten.WindowSize()

	// get the x start and end values based on the percent
	var xLeft = float64(screenWidth) * placement.StartPercentX
	var xRight = float64(screenWidth) * placement.EndPercentX
	var yTop = float64(screenHeight) * placement.StartPercentY
	var yBottom = float64(screenHeight) * placement.EndPercentY

	targetWidth := xRight - xLeft
	targetHeight := yBottom - yTop

	op.GeoM.Translate(xLeft, yTop)

	sb.referenceBorderBits = sb.createBorderBits()
	goodBorder := sb.createBorderImage(idealWidthForScaling, int(targetWidth), int(targetHeight))

	return goodBorder, op
}

func ScalePoint(point *image.Point, scaleX, scaleY float64) image.Point {
	return image.Point{
		X: int(math.Round(float64(point.X) * scaleX)),
		Y: int(math.Round(float64(point.Y) * scaleY)),
		//X: int(float64(point.X) * scaleX),
		//Y: int(float64(point.Y) * scaleY),
	}
}

func (r *referenceBorderDimensions) createScaledDimensions(scaleX, scaleY float64) *referenceBorderDimensions {
	referenceBorderDimensions := referenceBorderDimensions{
		topLeft: image.Rectangle{
			Min: ScalePoint(&r.topLeft.Min, scaleX, scaleY),
			Max: ScalePoint(&r.topLeft.Max, scaleX, scaleY),
		},
		bottomLeft: image.Rectangle{
			Min: ScalePoint(&r.bottomLeft.Min, scaleX, scaleY),
			Max: ScalePoint(&r.bottomLeft.Max, scaleX, scaleY),
		},
		topRight: image.Rectangle{
			Min: ScalePoint(&r.topRight.Min, scaleX, scaleY),
			Max: ScalePoint(&r.topRight.Max, scaleX, scaleY),
		},
		bottomRight: image.Rectangle{
			Min: ScalePoint(&r.bottomRight.Min, scaleX, scaleY),
			Max: ScalePoint(&r.bottomRight.Max, scaleX, scaleY),
		},
		horizCopyLeft: image.Rectangle{
			Min: ScalePoint(&r.horizCopyLeft.Min, scaleX, 1),
			Max: ScalePoint(&r.horizCopyLeft.Max, scaleX, 1),
		},
		vertCopyLeft: image.Rectangle{
			Min: ScalePoint(&r.vertCopyLeft.Min, 1, scaleY),
			Max: ScalePoint(&r.vertCopyLeft.Max, 1, scaleY),
		},
		horizCopyRight: image.Rectangle{
			Min: ScalePoint(&r.horizCopyRight.Min, scaleX, 1),
			Max: ScalePoint(&r.horizCopyRight.Max, scaleX, 1),
		},
		vertCopyRight: image.Rectangle{
			Min: ScalePoint(&r.vertCopyRight.Min, 1, scaleY),
			Max: ScalePoint(&r.vertCopyRight.Max, 1, scaleY),
		},
	}
	return &referenceBorderDimensions
}

// GetCornersOfReferenceBorder
// Determines the corners for the border, plus the horizontal and vertical slices it will fill in the blanks with
// NOTE: you must do this on the reference because you will muck up the aspect ratio if you do it AFTER a scale
func getCornersOfReferenceBorder(width, height int) *referenceBorderDimensions {
	referenceBorderDimensions := referenceBorderDimensions{}
	//xMiddle := int(float64(width) / 2)
	//yMiddle := int(float64(height) / 2)
	xMiddle := int(math.Round(float64(width) / 2))
	yMiddle := int(math.Round(float64(height) / 2))

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
	referenceBorderDimensions.horizCopyLeft = image.Rectangle{
		Min: image.Point{X: 0, Y: yMiddle},
		Max: image.Point{X: width / 2, Y: yMiddle + 1},
	}
	referenceBorderDimensions.vertCopyLeft = image.Rectangle{
		Min: image.Point{X: xMiddle, Y: 0},
		Max: image.Point{X: xMiddle + 1, Y: height / 2},
	}
	referenceBorderDimensions.horizCopyRight = image.Rectangle{
		Min: image.Point{X: width / 2, Y: yMiddle},
		Max: image.Point{X: width, Y: yMiddle + 1},
	}
	referenceBorderDimensions.vertCopyRight = image.Rectangle{
		Min: image.Point{X: xMiddle, Y: height / 2},
		Max: image.Point{X: xMiddle + 1, Y: height},
	}

	return &referenceBorderDimensions
}
