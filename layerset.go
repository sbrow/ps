// TODO: Count skipped steps.
package ps

import (
	"encoding/json"
	"errors"
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
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
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
	if err != nil {
		panic(err)
	}
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
