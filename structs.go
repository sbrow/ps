// TODO: Count skipped steps.
package ps

import (
	"encoding/json"
	"errors"
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

func (d *Document) Parent() Group {
	return nil
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
			if Mode != Fast && !set.current {
				set.Refresh()
			}
			return set
		}
	}
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
	d.name = strings.TrimRight(string(byt), "\r\n")
	if Mode != Safe {
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
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(byt, &d)
	if err != nil {
		d.Dump()
		fmt.Println(string(byt))
		log.Panic(err)
	}
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
	log.Println("Dumping to disk")
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

type LayerSet struct {
	name      string
	bounds    [2][2]int
	parent    Group
	current   bool // Whether we've checked this layer since we loaded from disk.
	visible   bool
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

type LayerSetJSON struct {
	Name      string
	Bounds    [2][2]int
	Visible   bool
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (l *LayerSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&LayerSetJSON{
		Name:      l.name,
		Bounds:    l.bounds,
		Visible:   l.visible,
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
	l.bounds = tmp.Bounds
	l.visible = tmp.Visible
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
	if Mode != 2 {
		for _, lyr := range l.artLayers {
			if !lyr.current {
				lyr.Refresh()
			}
		}
	}
	return l.artLayers
}

// ArtLayer returns the first top level ArtLayer matching
// the given name.
// TODO: Does funky things when passed invalid layername.
func (l *LayerSet) ArtLayer(name string) *ArtLayer {
	for _, lyr := range l.artLayers {
		if lyr.name == name {
			if Mode == 0 && !lyr.current {
				err := lyr.Refresh()
				if err != nil {
					l.Refresh()
					err := lyr.Refresh()
					if err != nil {
						log.Panic(err)
					}
				}
			}
			return lyr
		}
	}
	// l.Refresh()
	// for _, lyr := range l.artLayers {
	// fmt.Println(lyr)
	// }
	lyr := l.ArtLayer(name)
	fmt.Println(lyr)
	if lyr == nil {
		log.Panic(errors.New("Layer not found!"))
	}
	return lyr
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

// Bounds returns the furthest corners of the LayerSet.
func (l *LayerSet) Bounds() [2][2]int {
	return l.bounds
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
	if err != nil {
		log.Panic(err)
	}
	var out *LayerSet
	err = json.Unmarshal(byt, &out)
	if err != nil {
		log.Println(JSLayer(path))
		log.Println(string(byt))
		log.Panic(err)
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
	out.current = true
	return out, err
}

func (l *LayerSet) Visible() bool {
	return l.visible
}

// SetVisible makes the LayerSet visible.
func (l *LayerSet) SetVisible(b bool) {
	if l.visible == b {
		return
	}
	js := fmt.Sprintf("%s.visible=%v;", strings.TrimRight(
		JSLayer(l.Path()), ";"), b)
	DoJs("compilejs.jsx", js)
	l.visible = b
}

// SetPos snaps the given layerset boundry to the given point.
// Valid options for bound are: "TL", "TR", "BL", "BR"
func (l *LayerSet) SetPos(x, y int, bound string) {
	if !l.visible || (x == 0 && y == 0) {
		return
	}
	byt, err := DoJs("LayerSetBounds.jsx", JSLayer(l.Path()),
		JSLayer(l.Path(), true))
	if err != nil {
		log.Println(string(byt))
		log.Panic(err)
	}
	var bnds *[2][2]int
	err = json.Unmarshal(byt, &bnds)
	if err != nil {
		fmt.Println(string(byt))
		log.Panic(err)
	}
	l.bounds = *bnds
	var lyrX, lyrY int
	switch bound[:1] {
	case "B":
		lyrY = l.bounds[1][1]
	case "T":
		fallthrough
	default:
		lyrY = l.bounds[0][1]
	}
	switch bound[1:] {
	case "R":
		lyrX = l.bounds[1][0]
	case "L":
		fallthrough
	default:
		lyrX = l.bounds[0][0]
	}
	byt, err = DoJs("moveLayer.jsx", JSLayer(l.Path()), fmt.Sprint(x-lyrX),
		fmt.Sprint(y-lyrY), JSLayer(l.Path(), true))
	if err != nil {
		fmt.Println("byte:", string(byt))
		panic(err)
	}
	var lyr LayerSet
	err = json.Unmarshal(byt, &lyr)
	if err != nil {
		fmt.Println("byte:", string(byt))
		log.Panic(err)
	}
	l.bounds = lyr.bounds
}

func (l *LayerSet) Refresh() {
	var tmp *LayerSet
	byt, err := DoJs("getLayerSet.jsx", JSLayer(l.Path()), JSLayer(l.Path(), true))
	err = json.Unmarshal(byt, &tmp)
	if err != nil {
		log.Println("Error in LayerSet.Refresh() \"", string(byt), "\"", "for", l.Path())
		log.Panic(err)
	}
	tmp.SetParent(l.Parent())
	for _, lyr := range l.artLayers {
		err := lyr.Refresh()
		if err != nil {
			l.artLayers = tmp.artLayers
			break
		}
	}
	for _, set := range l.layerSets {
		set.Refresh()
	}
	l.name = tmp.name
	l.bounds = tmp.bounds
	l.visible = tmp.visible
	l.current = true
}

type TextItem struct {
	contents string
	size     float64
	// color    Color
	font   string
	parent *ArtLayer
}

type TextItemJSON struct {
	Contents string
	Size     float64
	// Color    [3]int
	Font string
}

func (t *TextItem) Contents() string {
	return t.contents
}

func (t *TextItem) Size() float64 {
	return t.size
}

// MarshalJSON implements the json.Marshaler interface, allowing the TextItem to be
// saved to disk in JSON format.
func (t *TextItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&TextItemJSON{
		Contents: t.contents,
		Size:     t.size,
		// Color:    t.color.RGB(),
		Font: t.font,
	})
}

func (t *TextItem) UnmarshalJSON(b []byte) error {
	tmp := &TextItemJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	t.contents = tmp.Contents
	t.size = tmp.Size
	// t.color = RGB{tmp.Color[0], tmp.Color[1], tmp.Color[2]}
	t.font = tmp.Font
	return nil
}

func (t *TextItem) SetText(txt string) {
	if txt == t.contents {
		return
	}
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	bndtext := "[[' + lyr.bounds[0] + ',' + lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + lyr.bounds[3] + ']]"
	js := fmt.Sprintf(`%s.textItem.contents='%s';var lyr = %[1]s;stdout.writeln(('%[3]s').replace(/ px/g, ''));`,
		lyr, txt, bndtext)
	byt, err := DoJs("compilejs.jsx", js)
	var bnds *[2][2]int
	json.Unmarshal(byt, &bnds)
	if err != nil || bnds == nil {
		log.Println("text:", txt)
		log.Println("js:", js)
		fmt.Printf("byt: '%s'\n", string(byt))
		log.Panic(err)
	}
	t.contents = txt
	t.parent.bounds = *bnds
}

func (t *TextItem) SetSize(s float64) {
	if t.size == s {
		return
	}
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	js := fmt.Sprintf("%s.textItem.size=%f;", lyr, s)
	_, err := DoJs("compilejs.jsx", js)
	if err != nil {
		t.size = s
	}
}

// TODO: Documentation for Format(), make to textItem
func (t *TextItem) Fmt(start, end int, font, style string) {
	var err error
	if !t.parent.Visible() {
		return
	}
	_, err = DoJs("fmtText.jsx", fmt.Sprint(start), fmt.Sprint(end),
		font, style)
	if err != nil {
		log.Panic(err)
	}
}