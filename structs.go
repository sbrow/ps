package ps

import (
	"encoding/json"
	"fmt"
)

type Group interface {
	Parent() Group
	GetArtLayers() []*ArtLayer
	GetLayerSets() []*LayerSet
}

type Document struct {
	Name      string
	Height    int
	Width     int
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
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

// ActiveDocument returns a Document object from Photoshop's active document.
func ActiveDocument() (Document, error) {
	byt, err := DoJs("getActiveDocument.jsx")
	var d Document
	err = json.Unmarshal(byt, &d)
	fmt.Println(string(byt))
	return d, err
}

func GetDocument() (*Document, error) {
	byt, err := DoJs("getActiveDoc.jsx")
	var d *Document
	err = json.Unmarshal(byt, &d)
	for _, lyr := range d.ArtLayers {
		lyr.Parent = d
	}
	for _, set := range d.LayerSets {
		set.parent = d
		for _, lyr := range set.ArtLayers {
			lyr.Parent = set
		}
	}
	fmt.Println(string(byt))
	return d, err
}

type ArtLayer struct {
	Name string
	// TextItem  string
	Bounds    [2][2]int
	Parent    Group
	Path      string
	Visiblity bool
}

// func (a *ArtLayer) setParent(c Group) {
// 	ArtLayer.Parent = c
// }

// Layer returns an ArtLayer from the active document given a specified
// path string.
func Layer(path string) (ArtLayer, error) {
	byt, err := DoJs("getLayer.jsx", JSLayer(path))
	var out ArtLayer
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return ArtLayer{}, err
	}
	out.Path = path
	return out, err
}

// SetVisible makes the layer visible.
func (a *ArtLayer) SetVisible() {
	js := JSLayer(a.Path) + fmt.Sprintf(".visible=%s;", true)
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
	byt, err := DoJs("moveLayer.jsx", JSLayer(a.Path), fmt.Sprint(x-lyrX), fmt.Sprint(y-lyrY))
	if err != nil {
		panic(err)
	}
	fmt.Println("byte", string(byt))
	fmt.Println("bounds", a.Bounds)
	json.Unmarshal(byt, &a)
	fmt.Println("after", a.Bounds)
}

type LayerSet struct {
	Name      string
	parent    Group
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

func (l *LayerSet) GetArtLayers() []*ArtLayer {
	return l.ArtLayers
}

func (l *LayerSet) GetLayerSets() []*LayerSet {
	return l.LayerSets
}

func (l *LayerSet) Parent() Group {
	return l.parent
}

func GetLayerSet(path string) (LayerSet, error) {
	byt, err := DoJs("getLayerSet.jsx", JSLayer(path))
	var out LayerSet
	fmt.Println(string(byt)) // TODO: Debug
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return LayerSet{}, err
	}
	for _, lyr := range out.ArtLayers {
		lyr.Parent = &out
	}
	return out, err
}
