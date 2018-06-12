package ps

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// ArtLayer reflects some values from an Art Layer in a Photoshop document.
//
// TODO: (2) Make TextLayer a subclass of ArtLayer.
type ArtLayer struct {
	name      string    // The layer's name.
	bounds    [2][2]int // The corners of the layer's bounding box.
	parent    Group     // The LayerSet/Document this layer is in.
	visible   bool      // Whether or not the layer is visible.
	current   bool      // Whether we've checked this layer since we loaded from disk.
	Color               // The layer's color overlay effect (if any).
	*Stroke             // The layer's stroke effect (if any).
	*TextItem           // The layer's text, if it's a text layer.
}

// Bounds returns the coordinates of the corners of the ArtLayer's bounding box.
func (a *ArtLayer) Bounds() [2][2]int {
	return a.bounds
}

// ArtLayerJSON is a bridge between the ArtLayer struct and
// the encoding/json package, allowing ArtLayer's unexported fields
// to ber written to and read from by the json package.
type ArtLayerJSON struct {
	Name      string
	Bounds    [2][2]int
	Visible   bool
	Color     [3]int
	Stroke    [3]int
	StrokeAmt float32
	TextItem  *TextItem
}

// MarshalJSON implements the json.Marshaler interface, allowing the ArtLayer to be
// saved to disk in JSON format.
func (a *ArtLayer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ArtLayerJSON{
		Name:      a.name,
		Bounds:    a.bounds,
		Visible:   a.visible,
		Color:     a.Color.RGB(),
		Stroke:    a.Stroke.RGB(),
		StrokeAmt: a.Stroke.Size,
		TextItem:  a.TextItem,
	})
}

func (a *ArtLayer) UnmarshalJSON(b []byte) error {
	tmp := &ArtLayerJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	a.name = tmp.Name
	a.bounds = tmp.Bounds
	a.Color = RGB{tmp.Color[0], tmp.Color[1], tmp.Color[2]}
	a.Stroke = &Stroke{tmp.StrokeAmt, RGB{tmp.Stroke[0], tmp.Stroke[1], tmp.Stroke[2]}}
	a.visible = tmp.Visible
	a.current = false
	a.TextItem = tmp.TextItem
	if a.TextItem != nil {
		a.TextItem.parent = a
	}
	return nil
}

func (a *ArtLayer) Name() string {
	return a.name
}

// Parent returns the Document or LayerSet this layer is contained in.
func (a *ArtLayer) Parent() Group {
	return a.parent
}

// X1 returns the layer's leftmost x value.
func (a *ArtLayer) X1() int {
	return a.bounds[0][0]
}

// X2 returns the layer's rightmost x value.
func (a *ArtLayer) X2() int {
	return a.bounds[1][0]
}

// Y1 returns the layer's topmost y value.
func (a *ArtLayer) Y1() int {
	return a.bounds[0][1]
}

// Y2 returns the layer's bottommost y value.
func (a *ArtLayer) Y2() int {
	return a.bounds[1][1]
}

func (a *ArtLayer) SetParent(c Group) {
	a.parent = c
}

// SetActive makes this layer active in Photoshop.
// Layers need to be active to perform certain operations
func (a *ArtLayer) SetActive() ([]byte, error) {
	js := fmt.Sprintf("app.activeDocument.activeLayer=%s", JSLayer(a.Path()))
	return DoJs("compilejs.jsx", js)
}

// SetColor creates a color overlay for the layer
func (a *ArtLayer) SetColor(c Color) {
	if a.Color.RGB() == c.RGB() {
		if Mode == 2 || (Mode == 0 && a.current) {
			// log.Println("Skipping color: already set.")
			return
		}
	}
	if a.Stroke.Size != 0 {
		a.SetStroke(*a.Stroke, c)
		return
	}
	a.Color = c
	cols := a.Color.RGB()
	log.Printf(`Setting layer "%s" to color %v`, a.name, cols)
	r := cols[0]
	g := cols[1]
	b := cols[2]
	byt, err := a.SetActive()
	if len(byt) != 0 {
		log.Println(string(byt), "err")
	}
	if err != nil {
		log.Println(a.Path())
		log.Panic(err)
	}
	byt, err = run("colorLayer", fmt.Sprint(r), fmt.Sprint(g), fmt.Sprint(b))
	if len(byt) != 0 {
		log.Println(string(byt), "err")
	}
	if err != nil {
		log.Panic(err)
	}
}

func (a *ArtLayer) SetStroke(stk Stroke, fill Color) {
	if stk.Size == 0 {
		a.Stroke = &stk
		a.SetColor(fill)
		return
	}
	if fill == nil {
		fill = a.Color
	}
	if stk.Size == a.Stroke.Size && stk.Color.RGB() == a.Stroke.Color.RGB() {
		if a.Color.RGB() == fill.RGB() {
			if Mode == 2 || (Mode == 0 && a.current) {
				// log.Println("Skipping stroke: already set.")
				return
			}
		}
	}
	byt, err := a.SetActive()
	if len(byt) != 0 {
		log.Println(string(byt))
	}
	if err != nil {
		log.Panic(err)
	}
	a.Stroke = &stk
	a.Color = fill
	stkCol := stk.Color.RGB()
	col := fill.RGB()
	log.Printf("Setting layer %s stroke to %.2fpt %v and color to %v\n", a.name, a.Stroke.Size,
		a.Stroke.Color.RGB(), a.Color.RGB())
	byt, err = run("colorStroke", fmt.Sprint(col[0]), fmt.Sprint(col[1]), fmt.Sprint(col[2]),
		fmt.Sprintf("%.2f", stk.Size), fmt.Sprint(stkCol[0]), fmt.Sprint(stkCol[1]), fmt.Sprint(stkCol[2]))
	if len(byt) != 0 {
		log.Println(string(byt))
	}
	if err != nil {
		log.Panic(err)
	}
}

func (a *ArtLayer) Path() string {
	return fmt.Sprintf("%s%s", a.parent.Path(), a.name)
}

// Layer returns an ArtLayer from the active document given a specified
// path string.
func layer(path string) (ArtLayer, error) {
	byt, err := DoJs("getLayer.jsx", JSLayer(path))
	if err != nil {
		return ArtLayer{}, err
	}
	var out ArtLayer
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return ArtLayer{}, err
	}
	return out, err
}

// SetVisible makes the layer visible.
func (a *ArtLayer) SetVisible(b bool) {
	if a.visible == b {
		return
	}
	a.visible = b
	switch b {
	case true:
		log.Printf("Showing %s", a.name)
	case false:
		log.Printf("Hiding %s", a.name)
	}
	js := fmt.Sprintf("%s.visible=%v;",
		strings.TrimRight(JSLayer(a.Path()), ";"), b)
	DoJs("compilejs.jsx", js)
}

// Visible returns whether or not the layer is currently hidden.
func (a *ArtLayer) Visible() bool {
	return a.visible
}

// SetPos snaps the given layer boundry to the given point.
// Valid options for bound are: "TL", "TR", "BL", "BR"
func (a *ArtLayer) SetPos(x, y int, bound string) {
	if !a.visible || (x == 0 && y == 0) {
		return
	}
	var lyrX, lyrY int
	switch bound[:1] {
	case "B":
		lyrY = a.Y2()
	case "T":
		fallthrough
	default:
		lyrY = a.Y1()
	}
	switch bound[1:] {
	case "R":
		lyrX = a.X2()
	case "L":
		fallthrough
	default:
		lyrX = a.X1()
	}
	byt, err := DoJs("moveLayer.jsx", JSLayer(a.Path()), fmt.Sprint(x-lyrX), fmt.Sprint(y-lyrY))
	if err != nil {
		panic(err)
	}
	var lyr ArtLayer
	err = json.Unmarshal(byt, &lyr)
	if err != nil {
		log.Panic(err)
	}
	a.bounds = lyr.bounds
}

func (a *ArtLayer) Refresh() error {
	tmp, err := layer(a.Path())
	if err != nil {
		return err
	}
	tmp.SetParent(a.Parent())
	a.name = tmp.name
	a.bounds = tmp.bounds
	a.TextItem = tmp.TextItem
	if a.TextItem != nil {
		a.TextItem.parent = a
	}
	a.parent = tmp.Parent()
	a.visible = tmp.visible
	a.current = true
	return nil
}
