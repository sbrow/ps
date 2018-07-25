package ps

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// LayerSet holds a group of Layer objects and a group of LayerSet objects.
type LayerSet struct {
	name      string
	bounds    [2][2]int
	parent    Group
	current   bool // Whether we've checked this layer since we loaded from disk.
	visible   bool
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

// LayerSetJSON is an exported version of LayerSet, that allows LayerSets to be
// saved to and loaded from JSON.
type LayerSetJSON struct {
	Name      string
	Bounds    [2][2]int
	Visible   bool
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

// MarshalJSON returns the LayerSet in JSON form.
func (l *LayerSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&LayerSetJSON{
		Name:      l.name,
		Bounds:    l.bounds,
		Visible:   l.visible,
		ArtLayers: l.artLayers,
		LayerSets: l.layerSets,
	})
}

// UnmarshalJSON loads the json data into this LayerSet.
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

// Name returns the name of the LayerSet
func (l LayerSet) Name() string {
	return l.name
}

// ArtLayers returns the LayerSet's ArtLayers.
func (l *LayerSet) ArtLayers() []*ArtLayer {
	if Mode != Fast {
		for _, lyr := range l.artLayers {
			if !lyr.current {
				if err := lyr.Refresh(); err != nil {
					log.Println(err)
				}
				lyr.current = true
			}
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
				err := lyr.Refresh()
				if err != nil {
					if err = l.Refresh(); err != nil {
						log.Panic(err)
					}
					if err := lyr.Refresh(); err != nil {
						log.Panic(err)
					}
				}
			}
			return lyr
		}
	}
	return nil
}

// LayerSets returns the LayerSets contained within
// this set.
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

// MustExist returns a Layer from the set with the given name, and
// panics if it doesn't exist.
//
// If there is a LayerSet and an  ArtLayer with the same name,
// it will return the LayerSet.
func (l *LayerSet) MustExist(name string) Layer {
	set := l.LayerSet(name)
	if set == nil {
		lyr := l.ArtLayer(name)
		if lyr == nil {
			log.Panicf("no Layer found at \"%s%s\"", l.Path(), name)
		}
		return lyr
	}
	return set
}

// Bounds returns the furthest corners of the LayerSet.
func (l LayerSet) Bounds() [2][2]int {
	return l.bounds
}

// SetParent puts this LayerSet into the given group.
func (l *LayerSet) SetParent(g Group) {
	l.parent = g
}

// Parent returns this layerSet's parent.
func (l *LayerSet) Parent() Group {
	return l.parent
}

// Path returns the layer path to this Set.
func (l *LayerSet) Path() string {
	if l.parent == nil {
		return l.name
	}
	return fmt.Sprintf("%s%s/", l.parent.Path(), l.name)
}

// NewLayerSet grabs the LayerSet with the given path and returns it.
func NewLayerSet(path string, g Group) (*LayerSet, error) {
	path = strings.Replace(path, "//", "/", -1)
	byt, err := DoJS("getLayerSet.jsx", JSLayer(path), JSLayerMerge(path))
	if err != nil {
		return nil, err
	}
	var out *LayerSet
	err = json.Unmarshal(byt, &out)
	if err != nil {
		log.Println(JSLayer(path))
		log.Println(string(byt))
		return nil, err
	}
	out.SetParent(g)
	log.Printf("Loading ActiveDocument/%s\n", out.Path())
	for _, lyr := range out.artLayers {
		lyr.SetParent(out)
	}
	for i, set := range out.layerSets {
		var s *LayerSet
		s, err = NewLayerSet(fmt.Sprintf("%s%s/", path, set.Name()), out)
		if err != nil {
			log.Fatal(err)
		}
		out.layerSets[i] = s
		s.SetParent(out)
	}
	out.current = true
	return out, err
}

// Visible returns whether or not the LayerSet is currently visible.
func (l LayerSet) Visible() bool {
	return l.visible
}

// SetVisible makes the LayerSet visible.
func (l *LayerSet) SetVisible(b bool) error {
	if l.visible == b {
		return nil
	}
	js := fmt.Sprintf("%s.visible=%v;", JSLayer(l.Path()), b)
	if _, err := DoJS("compilejs.jsx", js); err != nil {
		return err
	}
	l.visible = b
	return nil
}

// SetPos snaps the given layerset boundary to the given point.
// Valid options for bound are: "TL", "TR", "BL", "BR"
func (l *LayerSet) SetPos(x, y int, bound string) {
	if !l.visible || (x == 0 && y == 0) {
		return
	}
	path := JSLayer(l.Path())
	mrgPath := JSLayerMerge(l.Path())
	byt, err := DoJS("LayerSetBounds.jsx", path, mrgPath)
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
	byt, err = DoJS("moveLayer.jsx", JSLayer(l.Path()), fmt.Sprint(x-lyrX),
		fmt.Sprint(y-lyrY), JSLayerMerge(l.Path()))
	if err != nil {
		fmt.Println("byte:", string(byt))
		panic(err)
	}
	var lyr LayerSet
	if err = json.Unmarshal(byt, &lyr); err != nil {
		fmt.Println("byte:", string(byt))
		log.Panic(err)
	}
	l.bounds = lyr.bounds
}

// Refresh syncs the LayerSet with Photoshop.
func (l *LayerSet) Refresh() error {
	var tmp *LayerSet
	byt, err := DoJS("getLayerSet.jsx", JSLayer(l.Path()))
	if err != nil {
		if err = DoAction("Undo", "DK"); err != nil {
			return err
		}
		if byt, err = DoJS("getLayerSet.jsx", JSLayer(l.Path())); err != nil {
			return err
		}
	}
	err = json.Unmarshal(byt, &tmp)
	if err != nil {
		log.Println("Error in LayerSet.Refresh() \"", string(byt), "\"", "for", l.Path())
		return err
	}
	tmp.SetParent(l.Parent())
	for _, lyr := range l.artLayers {
		err = lyr.Refresh()
		if err != nil {
			l.artLayers = tmp.artLayers
			break
		}
	}
	for _, set := range l.layerSets {
		if err = set.Refresh(); err != nil {
			return err
		}
	}
	l.name = tmp.name
	l.bounds = tmp.bounds
	l.visible = tmp.visible
	l.current = true
	return nil
}
