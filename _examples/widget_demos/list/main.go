package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Struct used as data for combobox. You could just use strings instead of a struct.
type ListEntry struct {
	id   int
	Name string
}

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(20)

	// Create array of list entries
	numEntries := 20
	entries := make([]any, 0, numEntries)
	for i := 1; i <= numEntries; i++ {
		entries = append(entries, ListEntry{i, fmt.Sprintf("Entry %d", i)})
	}

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)

	// Construct a list. This is one of the more complicated widgets to use since
	// it is composed of multiple widget types

	list := widget.NewList(
		//Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchVertical:    true,
			}),
		)),
		//Set the entries in the list
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			//Set the background images/color for the list
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			//Set the background images/color for the background of the slider track
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(5),
			//Set how wide the track should be
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		//Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),
		//Set the font for the list options
		widget.ListOpts.EntryFontFace(face),
		//Set the colors for the list
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{0, 255, 0, 255},                 //Foreground color for the unfocused selected entry
			Unselected:                 color.NRGBA{254, 255, 255, 255},             //Foreground color for the unfocused unselected entry
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, //Background color for the unfocused selected entry
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, //Background color for the focused selected entry
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, //Background color for the focused unselected entry
			DisabledUnselected:         color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled unselected entry
			DisabledSelected:           color.NRGBA{100, 100, 100, 255},             //Foreground color for the disabled selected entry
			DisabledSelectedBackground: color.NRGBA{100, 100, 100, 255},             //Background color for the disabled selected entry
		}),
		//This required function returns the string displayed in the list
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(ListEntry).Name
		}),
		//Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		//This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			fmt.Println("Is: ", args.Entry, " Was: ", args.PreviousEntry)
		}),
	)

	//Add list to the root container
	rootContainer.AddChild(list)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - List")

	game := game{
		ui: &ui,
	}

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements Game.
func (g *game) Update() error {
	// update the UI
	g.ui.Update()
	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}