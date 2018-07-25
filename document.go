package ps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Document represents a Photoshop document (PSD file).
type Document struct {
	name      string
	fullName  string
	height    int
	width     int
	artLayers []*ArtLayer
	layerSets []*LayerSet
}

// DocumentJSON is an exported version of Document that
// allows Documents to be saved to and loaded from JSON.
type DocumentJSON struct {
	Name      string
	FullName  string
	Height    int
	Width     int
	ArtLayers []*ArtLayer
	LayerSets []*LayerSet
}

// MarshalJSON returns the Document in JSON format.
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DocumentJSON{Name: d.name, FullName: d.fullName, Height: d.height,
		Width: d.width, ArtLayers: d.artLayers, LayerSets: d.layerSets})
}

// UnmarshalJSON loads JSON data into this Document.
func (d *Document) UnmarshalJSON(b []byte) error {
	tmp := &DocumentJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	d.name = tmp.Name
	d.fullName = tmp.FullName
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
// This fulfills the Group interface.
func (d *Document) Name() string {
	return d.name
}

// FullName returns the absolute path to the current document file.
func (d *Document) FullName() string {
	return d.fullName
}

// Parent returns the Group that contains d.
func (d *Document) Parent() Group {
	return nil
}

// Height returns the height of the document, in pixels.
func (d *Document) Height() int {
	return d.height
}

// ArtLayer returns the first top level ArtLayer matching
// the given name.
func (d *Document) ArtLayer(name string) *ArtLayer {
	for _, lyr := range d.artLayers {
		if lyr.name == name {
			if Mode == 0 && !lyr.current {
				err := lyr.Refresh()
				if err != nil {
					log.Panic(err)
				}
			}
			return lyr
		}
	}
	return nil
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
// TODO(sbrow): Reduce cyclomatic complexity of ActiveDocument().
func ActiveDocument() (*Document, error) {
	log.Println("Loading ActiveDoucment")
	d := &Document{}

	byt, err := DoJS("activeDocName.jsx")
	if err != nil {
		return nil, err
	}
	d.name = strings.TrimRight(string(byt), "\r\n")
	if Mode != Safe {
		err = d.Restore(d.DumpFile())
		switch {
		case os.IsNotExist(err):
			log.Println("Previous version not found.")
		case err == nil:
			return d, err
		default:
			return nil, err

		}
	}
	log.Println("Loading manually (This could take awhile)")
	byt, err = DoJS("getActiveDoc.jsx")
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(byt, &d); err != nil {
		d.Dump()
		return nil, err
	}
	for _, lyr := range d.artLayers {
		lyr.SetParent(d)
	}
	for i, set := range d.layerSets {
		var s *LayerSet
		if s, err = NewLayerSet(set.Path()+"/", d); err != nil {
			return nil, err
		}
		d.layerSets[i] = s
		s.SetParent(d)
	}
	d.Dump()
	return d, err
}

// Restore loads document data from a JSON file.
func (d *Document) Restore(path string) error {
	if path == "" {
		path = d.DumpFile()
	}
	byt, err := ioutil.ReadFile(path)
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

// DumpFile returns the path to the json file where
// this document's data gets dumped. See Document.Dump
func (d *Document) DumpFile() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	path := filepath.Join(strings.Replace(d.fullName, "~", usr.HomeDir, 1))
	return strings.Replace(path, ".psd", ".json", 1)
}

// Dump saves the document to disk in JSON format.
func (d *Document) Dump() {
	log.Println("Dumping to disk")
	log.Println(d.DumpFile())
	defer d.Save()
	f, err := os.Create(d.DumpFile())
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

// MustExist returns a Layer from the set with the given name, and
// panics if it doesn't exist.
//
// If there is a LayerSet and an  ArtLayer with the same name,
// it will return the LayerSet.
func (d *Document) MustExist(name string) Layer {
	set := d.LayerSet(name)
	if set == nil {
		lyr := d.ArtLayer(name)
		if lyr == nil {
			log.Panicf("no Layer found at \"%s%s\"", d.Path(), name)
		}
		return lyr
	}
	return set
}

// Save saves the Document in place.
func (d *Document) Save() error {
	js := fmt.Sprintf("var d=app.open(File('%s'));\nd.save();", d.FullName())
	if _, err := DoJS("compilejs", js); err != nil {
		return err
	}
	return nil
}
