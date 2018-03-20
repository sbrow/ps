package ps

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func (d *Document) UnmarshalJSON(b []byte) error {
	tmp := &DocumentJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	d.name = tmp.Name
	d.height = tmp.Height
	d.width = tmp.Width
	d.artLayers = tmp.ArtLayers
	d.layerSets = tmp.LayerSets
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

func ActiveDocument() (*Document, error) {
	log.Println("Loading ActiveDoucment/")
	byt, err := DoJs("getActiveDoc.jsx")
	var d *Document
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
		// s.SetParent(d)
	}
	return d, err
}

// ArtLayer represents an Art Layer in a photoshop document.
type ArtLayer struct {
	name string
	// TextItem  string
	bounds  [2][2]int
	parent  Group
	visible bool
}

type ArtLayerJSON struct {
	Name    string
	Bounds  [2][2]int
	Parent  Group
	Visible bool
}

func (a *ArtLayer) UnmarshalJSON(b []byte) error {
	tmp := &ArtLayerJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	a.name = tmp.Name
	a.bounds = tmp.Bounds
	a.parent = tmp.Parent
	a.visible = tmp.Visible
	return nil
}

func (a *ArtLayer) Name() string {
	return a.name
}

func (a *ArtLayer) Bounds() [2][2]int {
	return a.bounds
}

// X1 returns the layer's leftmost x value.
func (a *ArtLayer) X1() int {
	return a.bounds[0][0]
}

// X2 returns the layer's rightmost x value.
func (a *ArtLayer) X2() int {
	return a.bounds[1][0]
}

// Y1 returns the layer's topmost Y value.
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
	js := fmt.Sprintf("%s.visible=%v;",
		strings.TrimRight(JSLayer(a.Path()), ";"), b)
	log.Printf("Setting %s.Visible to %v\n", a.name, b)
	DoJs("compilejs.jsx", js)
}

// Visible returns whether or not the layer is currently hidden.
func (a *ArtLayer) Visible() bool {
	return a.visible
}

// SetPos snaps the given layer boundry to the given point.
// Valid options for bound are: TL, TR, BL, BR
// TODO: Improve
func (a *ArtLayer) SetPos(x, y int, bound string) {
	var lyrX, lyrY int
	lyrX = a.X1()
	if bound != "TL" {
		lyrY = a.Y2()
	} else { // "BL"
		lyrY = a.Y1()
	}
	byt, err := DoJs("moveLayer.jsx", JSLayer(a.Path()), fmt.Sprint(x-lyrX), fmt.Sprint(y-lyrY))
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byt, &a)
}

type LayerSet struct {
	name      string
	parent    Group
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

type LayerSetJSON struct {
	Name      string
	Parent    Group
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (l *LayerSet) UnmarshalJSON(b []byte) error {
	tmp := &LayerSetJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	l.name = tmp.Name
	l.parent = tmp.Parent
	l.artLayers = tmp.ArtLayers
	l.layerSets = tmp.LayerSets
	return nil
}

func (l *LayerSet) Name() string {
	return l.name
}

func (l *LayerSet) ArtLayers() []*ArtLayer {
	return l.artLayers
}

// ArtLayer returns the first top level ArtLayer matching
// the given name.
func (l *LayerSet) ArtLayer(name string) *ArtLayer {
	for _, lyr := range l.artLayers {
		if lyr.name == name {
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
	byt, err := DoJs("getLayerSet.jsx", JSLayer(path))
	var out *LayerSet
	err = json.Unmarshal(byt, &out)
	if flag.Lookup("test.v") != nil {
		// log.Println(string(byt))
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
		// log.Println("\t", set.name)
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
