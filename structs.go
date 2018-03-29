package ps

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Group represents a Document or LayerSet.
type Group interface {
	Name() string
	Parent() Group
	SetParent(Group)
	Path() string
	ArtLayers() []*ArtLayer
	LayerSets() []*LayerSet
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

// Document represents a Photoshop document (PSD file).
type Document struct {
	name      string
	height    int
	width     int
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

type DocumentJSON struct {
	Name      string
	Height    int
	Width     int
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DocumentJSON{Name: d.name, Height: d.height,
		Width: d.width, ArtLayers: d.artLayers, LayerSets: d.layerSets})
}

func (d *Document) UnmarshalJSON(b []byte) error {
	tmp := &DocumentJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	d.name = tmp.Name
	d.height = tmp.Height
	d.width = tmp.Width
	d.artLayers = tmp.ArtLayers
	for _, lyr := range d.artLayers {
		lyr.SetParent(d)
	}
	d.layerSets = tmp.LayerSets
	for _, set := range d.layerSets {
		set.SetParent(d)
	}
	return nil
}

// Name returns the document's title.
// This fufills the Group interface.
func (d *Document) Name() string {
	return d.name
}

// The height of the document, in pixels.
func (d *Document) Height() int {
	return d.height
}

func (d *Document) ArtLayers() []*ArtLayer {
	return d.artLayers
}

// LayerSets returns all the document's top level LayerSets.
func (d *Document) LayerSets() []*LayerSet {
	return d.layerSets
}

// LayerSet returns the first top level LayerSet matching
// the given name.
func (d *Document) LayerSet(name string) *LayerSet {
	for _, set := range d.layerSets {
		if set.name == name {
			return set
		}
	}
	return nil
}

func (d *Document) Parent() Group {
	return nil
}
func (d *Document) SetParent(g Group) {}

func (d *Document) Path() string {
	return ""
}

// Filename returns the path to the json file for this document.
func (d *Document) Filename() string {
	_, dir, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("No caller information")
	}
	return filepath.Join(filepath.Dir(dir), "data",
		strings.TrimRight(string(d.name), "\r\n")+".txt")
}

func ActiveDocument() (*Document, error) {
	log.Println("Loading ActiveDoucment")
	d := &Document{}

	byt, err := DoJs("activeDocName.jsx")
	if err != nil {
		return nil, err
	}
	d.name = string(byt)
	if Mode != 1 {
		byt, err = ioutil.ReadFile(d.Filename())
		if err == nil {
			log.Println("Previous version found, loading")
			err = json.Unmarshal(byt, &d)
			if err == nil {
				return d, err
			}
		}
	}
	log.Println("Loading manually (This could take awhile)")
	byt, err = DoJs("getActiveDoc.jsx")
	err = json.Unmarshal(byt, &d)
	for _, lyr := range d.artLayers {
		lyr.SetParent(d)
	}
	for i, set := range d.layerSets {
		s, err := NewLayerSet(set.Path()+"/", d)
		if err != nil {
			log.Fatal(err)
		}
		d.layerSets[i] = s
		s.SetParent(d)
	}
	d.Dump()
	return d, err
}

func (d *Document) Dump() {
	f, err := os.Create(d.Filename())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	byt, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(byt)
}

// ArtLayer reflects certain values from an Art Layer in a Photoshop document.
type ArtLayer struct {
	name    string    // The layer's name.
	Text    *string   // The contents of a text layer.
	bounds  [2][2]int // The layers' corners.
	parent  Group     // The LayerSet/Document this layer is in.
	visible bool      // Whether or not the layer is visible.
	current bool      // Whether we've checked this layer since we loaded from disk.
	Color             // The layer's color overlay.
	*Stroke           // The layer's stroke.
}

// Bounds returns the furthest corners of the ArtLayer.
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
	Text      *string
}

// MarshalJSON fulfills the json.Marshaler interface, allowing the ArtLayer to be
// saved to disk in JSON format.
func (a *ArtLayer) MarshalJSON() ([]byte, error) {
	// txt := strings.Replace(*a.Text, "\r", "\\r", -1)
	return json.Marshal(&ArtLayerJSON{
		Name:      a.name,
		Bounds:    a.bounds,
		Visible:   a.visible,
		Color:     a.Color.RGB(),
		Stroke:    a.Stroke.RGB(),
		StrokeAmt: a.Stroke.Size,
		Text:      a.Text,
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
	if tmp.Text != nil {
		// s := strings.Replace(*tmp.Text, "\\r", "\r", -1)
		a.Text = tmp.Text
	}
	a.current = false
	return nil
}

func (a *ArtLayer) Name() string {
	return a.name
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
	if Mode == 2 && a.Color.RGB() == c.RGB() {
		return
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
	if Mode == 2 {
		if stk.Size == a.Stroke.Size && stk.Color.RGB() == a.Color.RGB() {
			if a.Color.RGB() == fill.RGB() {
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

func (a *ArtLayer) Parent() Group {
	return a.parent
}

func (a *ArtLayer) Path() string {
	return fmt.Sprintf("%s%s", a.parent.Path(), a.name)
}

// Layer returns an ArtLayer from the active document given a specified
// path string.
func Layer(path string) (ArtLayer, error) {
	byt, err := DoJs("getLayer.jsx", JSLayer(path))
	var out ArtLayer
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return ArtLayer{}, err
	}
	return out, err
}

// SetVisible makes the layer visible.
func (a *ArtLayer) SetVisible(b bool) {
	if a.Visible() == b {
		return
	}
	if b {
		log.Printf("Showing %s", a.name)
	} else {
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
//
// TODO: Test TR and BR
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
	var bounds [2][2]int
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byt, bounds)
	a.bounds = bounds
}

type LayerSet struct {
	name      string
	parent    Group
	current   bool // Whether we've checked this layer since we loaded from disk.
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

type LayerSetJSON struct {
	Name      string
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (l *LayerSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&LayerSetJSON{
		Name:      l.name,
		ArtLayers: l.artLayers,
		LayerSets: l.layerSets,
	})
}

func (l *LayerSet) UnmarshalJSON(b []byte) error {
	tmp := &LayerSetJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	l.name = tmp.Name
	l.artLayers = tmp.ArtLayers
	for _, lyr := range l.artLayers {
		lyr.SetParent(l)
	}
	l.layerSets = tmp.LayerSets
	for _, set := range l.layerSets {
		set.SetParent(l)
	}
	l.current = false
	return nil
}

func (l *LayerSet) Name() string {
	return l.name
}

func (l *LayerSet) ArtLayers() []*ArtLayer {
	for i := 0; i < len(l.artLayers); i++ {
		if l.artLayers[i] != nil && !l.artLayers[i].current {
			l.artLayers[i] = l.ArtLayer(l.artLayers[i].name)
		}
	}
	return l.artLayers
}

// ArtLayer returns the first top level ArtLayer matching
// the given name.
func (l *LayerSet) ArtLayer(name string) *ArtLayer {
	for _, lyr := range l.artLayers {
		if lyr.name == name {
			if Mode == 0 && !lyr.current {
				byt, err := DoJs("getLayer.jsx", JSLayer(lyr.Path()))
				if err != nil {
					log.Panic(err)
				}
				var lyr2 *ArtLayer
				err = json.Unmarshal(byt, &lyr2)
				if err != nil {
					log.Panic(err)
				}
				lyr.name = lyr2.name
				lyr.bounds = lyr2.bounds
				lyr.visible = lyr2.visible
				lyr.current = true
				lyr.Text = lyr2.Text
			}
			return lyr
		}
	}
	return nil
}

func (l *LayerSet) LayerSets() []*LayerSet {
	return l.layerSets
}

// LayerSet returns the first top level LayerSet matching
// the given name.
func (l *LayerSet) LayerSet(name string) *LayerSet {
	for _, set := range l.layerSets {
		if set.name == name {
			return set
		}
	}
	return nil
}

func (l *LayerSet) SetParent(c Group) {
	l.parent = c
}

func (l *LayerSet) Parent() Group {
	return l.parent
}

func (l *LayerSet) Path() string {
	if l.parent == nil {
		return l.name
	}
	return fmt.Sprintf("%s%s/", l.parent.Path(), l.name)
}

func NewLayerSet(path string, g Group) (*LayerSet, error) {
	path = strings.Replace(path, "//", "/", -1)
	byt, err := DoJs("getLayerSet.jsx", JSLayer(path))
	var out *LayerSet
	err = json.Unmarshal(byt, &out)
	if err != nil {
		log.Println(string(byt))
		log.Panic(err)
	}
	if flag.Lookup("test.v") != nil {
		// log.Println(path)
		// log.Println(out)
	}
	out.SetParent(g)
	log.Printf("Loading ActiveDocument/%s\n", out.Path())
	if err != nil {
		return &LayerSet{}, err
	}
	for _, lyr := range out.artLayers {
		lyr.SetParent(out)
	}
	for i, set := range out.layerSets {
		s, err := NewLayerSet(fmt.Sprintf("%s%s/", path, set.Name()), out)
		if err != nil {
			log.Fatal(err)
		}
		out.layerSets[i] = s
		s.SetParent(out)
	}
	return out, err
}

// SetVisible makes the LayerSet visible.
func (l *LayerSet) SetVisible(b bool) {
	js := fmt.Sprintf("%s%v", JSLayer(strings.TrimRight(l.Path(), ";")), b)
	DoJs("compilejs.jsx", js)
}
