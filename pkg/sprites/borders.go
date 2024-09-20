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
	referenceBorderBits       *BorderBits
}

type ReferenceBorderDimensions struct {
	topLeft, bottomLeft, topRight, bottomRight image.Rectangle
	horizCopy, vertCopy                        image.Rectangle
}

type BorderBits struct {
	topLeft     *ebiten.Image
	bottomRight *ebiten.Image
	topRight    *ebiten.Image
	bottomLeft  *ebiten.Image
	vertCopy    *ebiten.Image
	horizCopy   *ebiten.Image
}

func NewBorderSprites() *BorderSprites {
	borderSprites := &BorderSprites{}
	borderSprites.VeryPixelatedRoundedBlueWhite.border = NewPngSprite(veryPixelatedRoundedBlueWhite)

	return borderSprites
}

func (sb *ReferenceBorder) CreateBorderImage(width, height int) *ebiten.Image {
	windowWidth, _ := ebiten.WindowSize()
	idealScaleX := float64(windowWidth) / float64(475)
	scaledDimensions := sb.referenceBorderDimensions.CreateScaledDimensions(idealScaleX, idealScaleX)

	img := ebiten.NewImage(width, height)

	topLeftOp := ebiten.GeoM{}
	topLeftOp.Scale(idealScaleX, idealScaleX)
	img.DrawImage(sb.referenceBorderBits.topLeft, &ebiten.DrawImageOptions{
		GeoM: topLeftOp,
	})

	bottomLeftOp := ebiten.GeoM{}
	bottomLeftOp.Scale(idealScaleX, idealScaleX)
	bottomLeftOp.Translate(0,
		float64(height-scaledDimensions.bottomLeft.Min.Y),
	)
	img.DrawImage(sb.referenceBorderBits.bottomLeft, &ebiten.DrawImageOptions{GeoM: bottomLeftOp})

	topRightOp := ebiten.GeoM{}
	topRightOp.Scale(idealScaleX, idealScaleX)
	topRightOp.Translate(
		float64(width-scaledDimensions.topRight.Min.X),
		0,
	)
	img.DrawImage(sb.referenceBorderBits.topRight, &ebiten.DrawImageOptions{GeoM: topRightOp})

	bottomRightOp := ebiten.GeoM{}
	bottomRightOp.Scale(idealScaleX, idealScaleX)
	bottomRightOp.Translate(
		float64(width-scaledDimensions.bottomRight.Min.X),
		float64(height-scaledDimensions.bottomRight.Min.Y),
	)
	img.DrawImage(sb.referenceBorderBits.bottomRight, &ebiten.DrawImageOptions{GeoM: bottomRightOp})

	return img

}

func (sb *ReferenceBorder) CreateBorderBits() *BorderBits {
	bits := BorderBits{}

	// make our corner bits
	bits.topLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.topLeft))
	bits.bottomLeft = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.bottomLeft))
	bits.topRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.topRight))
	bits.bottomRight = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.bottomRight))
	bits.horizCopy = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.horizCopy))
	bits.vertCopy = ebiten.NewImageFromImage(sb.border.SubImage(sb.referenceBorderDimensions.vertCopy))

	return &bits
}

// CreateSizedAndScaledBorderSprite
func (sb *ReferenceBorder) CreateSizedAndScaledBorderSprite(placement PercentBasedPlacement) (*ebiten.Image, *ebiten.DrawImageOptions) {

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

	//originalImgWidth := sb.border.Bounds().Dx()
	//originalImgHeight := sb.border.Bounds().Dy()
	//
	//// we scale to X value, and match the ratio on Y so the image doesn't stretch and distort
	//scaleX := targetWidth / float64(originalImgWidth)
	//scaleY := targetHeight / float64(originalImgHeight)
	//op.GeoM.Scale(scaleX, scaleY)
	//
	//scaledWidth := scaleX * float64(originalImgWidth)
	////scaledHeight := scaleY * float64(originalImgHeight)
	//
	//topLeftX := (float64(screenWidth) - scaledWidth) / 2
	//op.GeoM.Translate(topLeftX, yTop)
	op.GeoM.Translate(50, 400)

	//scaledDimensions := sb.referenceBorderDimensions.CreateScaledDimensions(scaleX, scaleY)

	// border bits
	sb.referenceBorderBits = sb.CreateBorderBits()

	// scale the border to screen scale

	// find the corners
	// scale the corners

	// fill in horiz and vert
	goodBorder := sb.CreateBorderImage(int(targetWidth), int(targetHeight))

	return goodBorder, op
}

//func (sb *ReferenceBorder) CreateSizedAndScaledBorderSprite(placement PercentBasedPlacement) (*ebiten.Image, *ebiten.DrawImageOptions) {
//
//	// get the corners and copies before we scale it
//	sb.referenceBorderDimensions = getCornersOfReferenceBorder(sb.border.Bounds().Dx(), sb.border.Bounds().Dy())
//
//	op := &ebiten.DrawImageOptions{}
//
//	// THIS STUFF IS FOR PLACEMENT AFTER WE HAVE THE IMAGE
//	// get screen sizes
//	screenWidth, screenHeight := ebiten.WindowSize()
//
//	// get the x start and end values based on the percent
//	var xLeft = float64(screenWidth) * placement.StartPercentX
//	var xRight = float64(screenWidth) * placement.EndPercentX
//	var yTop = float64(screenHeight) * placement.StartPercentY
//	var yBottom = float64(screenHeight) * placement.EndPercentY
//
//	targetWidth := xRight - xLeft
//	targetHeight := yBottom - yTop
//
//	originalImgWidth := sb.border.Bounds().Dx()
//	originalImgHeight := sb.border.Bounds().Dy()
//
//	// we scale to X value, and match the ratio on Y so the image doesn't stretch and distort
//	scaleX := targetWidth / float64(originalImgWidth)
//	scaleY := targetHeight / float64(originalImgHeight)
//	op.GeoM.Scale(scaleX, scaleY)
//
//	scaledWidth := scaleX * float64(originalImgWidth)
//	//scaledHeight := scaleY * float64(originalImgHeight)
//
//	topLeftX := (float64(screenWidth) - scaledWidth) / 2
//	op.GeoM.Translate(topLeftX, yTop)
//
//	//scaledDimensions := sb.referenceBorderDimensions.CreateScaledDimensions(scaleX, scaleY)
//
//	// border bits
//	sb.referenceBorderBits = sb.CreateBorderBits()
//
//	// scale the border to screen scale
//
//	// find the corners
//	// scale the corners
//
//	// fill in horiz and vert
//	goodBorder := sb.CreateBorderImage(int(targetWidth), int(targetHeight))
//
//	return goodBorder, op
//}

func ScalePoint(point *image.Point, scaleX, scaleY float64) image.Point {
	return image.Point{
		X: int(float64(point.X) * scaleX),
		Y: int(float64(point.Y) * scaleY),
	}
}

func (r *ReferenceBorderDimensions) CreateScaledDimensions(scaleX, scaleY float64) *ReferenceBorderDimensions {
	referenceBorderDimensions := ReferenceBorderDimensions{
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
		horizCopy: image.Rectangle{
			Min: ScalePoint(&r.horizCopy.Min, scaleX, 1),
			Max: ScalePoint(&r.horizCopy.Max, scaleX, 1),
		},
		vertCopy: image.Rectangle{
			Min: ScalePoint(&r.vertCopy.Min, 1, scaleY),
			Max: ScalePoint(&r.vertCopy.Max, 1, scaleY),
		},
	}
	return &referenceBorderDimensions
}

// GetCornersOfReferenceBorder
// Determines the corners for the border, plus the horizontal and vertical slices it will fill in the blanks with
// NOTE: you must do this on the reference because you will muck up the aspect ratio if you do it AFTER a scale
func getCornersOfReferenceBorder(width, height int) *ReferenceBorderDimensions {
	referenceBorderDimensions := ReferenceBorderDimensions{}
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
	referenceBorderDimensions.horizCopy = image.Rectangle{
		Min: image.Point{X: 0, Y: yMiddle},
		Max: image.Point{X: width - 1, Y: yMiddle + 1},
		//Max: image.Point{X: width - 1, Y: yMiddle},
	}
	referenceBorderDimensions.vertCopy = image.Rectangle{
		Min: image.Point{X: xMiddle, Y: 0},
		Max: image.Point{X: xMiddle + 1, Y: height - 1},
		//Max: image.Point{X: xMiddle, Y: height - 1},
	}
	return &referenceBorderDimensions
}
