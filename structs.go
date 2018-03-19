package ps

import (
	"encoding/json"
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
	GetArtLayers() []*ArtLayer
	GetLayerSets() []*LayerSet
}

// Document represents a Photoshop document (PSD file).
type Document struct {
	name      string
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
	d.Height = tmp.Height
	d.Width = tmp.Width
	d.ArtLayers = tmp.ArtLayers
	d.LayerSets = tmp.LayerSets
	return nil
}

type DocumentJSON struct {
	Name      string
	Height    int
	Width     int
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (d *Document) Name() string {
	return d.name
}
func (d *Document) GetArtLayers() []*ArtLayer {
	return d.ArtLayers
}

func (d *Document) GetLayerSets() []*LayerSet {
	return d.LayerSets
}

func (d *Document) Parent() Group {
	return nil
}
func (d *Document) SetParent(g Group) {}

func (d *Document) Path() string {
	return ""
}

func ActiveDocument() (*Document, error) {
	byt, err := DoJs("getActiveDoc.jsx")
	var d *Document
	err = json.Unmarshal(byt, &d)
	for _, lyr := range d.ArtLayers {
		lyr.SetParent(d)
	}
	for i, set := range d.LayerSets {
		s, err := NewLayerSet(set.Path() + "/")
		if err != nil {
			log.Fatal(err)
		}
		d.LayerSets[i] = s
		s.SetParent(d)
	}
	return d, err
}

type ArtLayer struct {
	name string
	// TextItem  string
	Bounds    [2][2]int
	parent    Group
	Visiblity bool
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
	a.Bounds = tmp.Bounds
	a.parent = tmp.Parent
	a.Visiblity = tmp.Visible
	return nil
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
func (a *ArtLayer) SetVisible() {
	js := JSLayer(a.Path()) + fmt.Sprintf(".visible=%s;", true)
	DoJs("compilejs.jsx", js)
}

// Position moves the layer to pos(x, y), measuring from
// the top or bottom left-hand corner.
// TODO: Improve
func (a *ArtLayer) Position(x, y int, align string) {
	var lyrX, lyrY int
	lyrX = a.Bounds[0][0]
	if align != "bottom" {
		lyrY = a.Bounds[0][1]
	} else {
		lyrY = a.Bounds[1][1]
	}
	byt, err := DoJs("moveLayer.jsx", JSLayer(a.Path()), fmt.Sprint(x-lyrX), fmt.Sprint(y-lyrY))
	if err != nil {
		panic(err)
	}
	fmt.Println("byte", string(byt))
	fmt.Println("bounds", a.Bounds)
	json.Unmarshal(byt, &a)
	fmt.Println("after", a.Bounds)
}

type LayerSet struct {
	name      string
	parent    Group
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
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
	l.ArtLayers = tmp.ArtLayers
	l.LayerSets = tmp.LayerSets
	return nil
}

func (l *LayerSet) Name() string {
	return l.name
}

func (l *LayerSet) GetArtLayers() []*ArtLayer {
	return l.ArtLayers
}

func (l *LayerSet) GetLayerSets() []*LayerSet {
	return l.LayerSets
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

func NewLayerSet(path string) (*LayerSet, error) {
	byt, err := DoJs("getLayerSet.jsx", JSLayer(path))
	// fmt.Println(string(byt))
	var out *LayerSet
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return &LayerSet{}, err
	}
	for _, lyr := range out.ArtLayers {
		lyr.SetParent(out)
	}
	for i, set := range out.LayerSets {
		s, err := NewLayerSet(fmt.Sprintf("%s%s/", path, set.Name()))
		if err != nil {
			log.Fatal(err)
		}
		out.LayerSets[i] = s
		s.SetParent(out)
	}
	return out, err
}

// SetVisible makes the LayerSet visible.
func (s *LayerSet) SetVisible(b bool) {
	fmt.Println(s.Path())
	fmt.Println(JSLayer(s.Path()))
	js := JSLayer(strings.TrimRight(s.Path(), ";") + fmt.Sprintf(".visible=%v;", b))
	DoJs("compilejs.jsx", js)
}
