package ps

import (
	"encoding/json"
	"fmt"
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
		strings.TrimRight(d.name, "\r\n")+".txt")
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
