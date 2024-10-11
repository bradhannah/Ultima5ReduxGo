package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	c_color "image/color"
)

type DebugConsole struct {
	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	ui ebitenui.UI

	gameScene *GameScene
}

func NewDebugConsole(gameScene *GameScene) *DebugConsole {
	debugConsole := DebugConsole{}
	debugConsole.gameScene = gameScene
	debugConsole.initializeDebugBorders()

	font := text.NewUltimaFont(14)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(c_color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x77})),

		//widget.ContainerOpts.Layout(widget.NewAnchorLayout()
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	rect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
		StartPercentX: .015,
		EndPercentX:   .745,
		StartPercentY: .73,
		EndPercentY:   0.98})

	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionEnd,
					MaxWidth:  rect.Dx(),
					MaxHeight: rect.Dy(),
				}),
				widget.WidgetOpts.MinSize(rect.Dx(), rect.Dy()),

				//widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				//	HorizontalPosition: 0,
				//	VerticalPosition:   500,
				//	StretchHorizontal:  true,
				//	StretchVertical:    false,
				//	Padding:            widget.Insets{},
				//}),
			),
		),
		//Set gap between scrollbar and text
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		//Tell the textarea to display bbcodes
		widget.TextAreaOpts.ProcessBBCode(true),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.White),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(font.TextFace),
		//Set the initial text for the textarea
		//It will automatically line wrap and process newlines characters
		//If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text("Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4Hello World\nTest1\nTest2\n[color=ff0000]Red[/color]\n[color=00ff00]Green[/color]\n[color=0000ff]Blue[/color]\nTest3\nTest4"),
		//Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(10)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.BlackSemi),
				Mask: image.NewNineSliceColor(c_color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				//Idle: image.NewNineSliceColor(c_color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				//Mask: image.NewNineSliceColor(c_color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{

					Idle:  image.NewNineSliceColor(color.Black),
					Hover: image.NewNineSliceColor(c_color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
					//Idle:  image.NewNineSliceColor(c_color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
					//Hover: image.NewNineSliceColor(c_color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle: image.NewNineSliceColor(color.LighterBlueSemi),
					//Idle:    image.NewNineSliceColor(c_color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
					Hover:   image.NewNineSliceColor(c_color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
					Pressed: image.NewNineSliceColor(c_color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				},
			),
		),
	)

	//rootContainer.AddChild(innerContainer)
	rootContainer.AddChild(textarea)

	debugConsole.ui = ebitenui.UI{
		Container: rootContainer,
	}

	return &debugConsole
}

func (d *DebugConsole) update() {
	d.ui.Update()

	if ebiten.IsKeyPressed(ebiten.KeyBackquote) {
		if !d.gameScene.keyboard.TryToRegisterKeyPress(ebiten.KeyBackquote) {
			return
		}
		d.gameScene.bShowDebugConsole = false
	}

	return
}

func (d *DebugConsole) drawDebugConsole(screen *ebiten.Image) {
	screen.DrawImage(d.background, d.backgroundDrawOptions)
	d.ui.Draw(screen)
	screen.DrawImage(d.border, d.borderDrawOptions)
}

func (d *DebugConsole) initializeDebugBorders() {
	mainBorder := sprites.NewBorderSprites()
	//const percentOffEdge = 0.04
	const percentOffEdge = 0.0
	percentBased := sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentOffEdge,
		EndPercentX:   .75 + .01 - percentOffEdge,
		StartPercentY: .7,
		EndPercentY:   1,
	}
	d.border, d.borderDrawOptions = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, percentBased)

	d.backgroundDrawOptions = &ebiten.DrawImageOptions{}
	*d.backgroundDrawOptions = *d.borderDrawOptions

	backgroundPercents := percentBased

	rect := sprites.GetRectangleFromPercents(backgroundPercents)

	d.background = ebiten.NewImageWithOptions(*rect, &ebiten.NewImageOptions{})
	xDiff := float32(rect.Dx()) * 0.01
	yDiff := float32(rect.Dy()) * 0.01
	vector.DrawFilledRect(d.background,
		float32(rect.Min.X)+xDiff,
		float32(rect.Min.Y)+yDiff,
		float32(rect.Dx())-xDiff*2,
		float32(rect.Dy())-yDiff*2,
		color.Black,
		false)

	d.backgroundDrawOptions.ColorScale.ScaleAlpha(.85)
}
