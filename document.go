package ps

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Document represents a Photoshop document (PSD file).
type Document struct {
	name      string
	height    int
	width     int
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

// DocumentJSON is an exported version of Document that
// allows Documents to be saved to and loaded from JSON.
type DocumentJSON struct {
	Name      string
	Height    int
	Width     int
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

// MarshalJSON returns the Document in JSON format.
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DocumentJSON{Name: d.name, Height: d.height,
		Width: d.width, ArtLayers: d.artLayers, LayerSets: d.layerSets})
}

// UnmarshalJSON loads JSON data into this Document.
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

// Parent returns the Group that contains d.
func (d *Document) Parent() Group {
	return nil
}

// Height returns the height of the document, in pixels.
func (d *Document) Height() int {
	return d.height
}

// ArtLayers returns this document's ArtLayers, if any.
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
				if err := set.Refresh(); err != nil {
					log.Panic(err)
				}
			}
			return set
		}
	}
	return nil
}

// ActiveDocument returns document currently focused in Photoshop.
//
// TODO: Reduce cylcomatic complexity
func ActiveDocument() (*Document, error) {
	log.Println("Loading ActiveDoucment")
	d := &Document{}

	byt, err := DoJS("activeDocName.jsx")
	if err != nil {
		return nil, err
	}
	d.name = strings.TrimRight(string(byt), "\r\n")
	if Mode != Safe {
		err = d.Restore()
		return d, err
	}
	log.Println("Loading manually (This could take awhile)")
	byt, err = DoJS("getActiveDoc.jsx")
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(byt, &d)
	if err != nil {
		d.Dump()
		log.Panic(err)
	}
	for _, lyr := range d.artLayers {
		lyr.SetParent(d)
	}
	for i, set := range d.layerSets {
		var s *LayerSet
		s, err = NewLayerSet(set.Path()+"/", d)
		if err != nil {
			log.Fatal(err)
		}
		d.layerSets[i] = s
		s.SetParent(d)
	}
	d.Dump()
	return d, err
}

// Restore loads document data from a JSON file.
func (d *Document) Restore() error {
	byt, err := ioutil.ReadFile(d.Filename())
	if err == nil {
		log.Println("Previous version found, loading")
		err = json.Unmarshal(byt, &d)
	}
	return err
}

// SetParent does nothing, as the document is a top-level object
// and therefore can't have a parent group.
// The function is needed to implement the group interface.
func (d *Document) SetParent(g Group) {}

// Path returns the root path ("") for all the layers.
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
		strings.TrimRight(d.name, "\r\n")+".txt")
}

// Dump saves the document to disk in JSON format.
func (d *Document) Dump() {
	log.Println("Dumping to disk")
	log.Println(d.Filename())
	f, err := os.Create(d.Filename())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()
	byt, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.Write(byt); err != nil {
		log.Println(err)
	}
}
