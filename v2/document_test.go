package ps

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// wd is the working directory
var wd string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("runtime.Caller(0) returned !ok")
	}
	wd = filepath.Dir(file)

	f, err := os.Create("test.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
}

func TestDocument_Dump(t *testing.T) {
	// Must be true for test to be valid.
	Mode = Normal

	tests := []struct {
		name string
	}{
		{"Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Get rid of old Dump.
			os.Remove(filepath.Join(wd, tt.name+".json"))

			// Generate a fresh Doc (loaded slowly).
			want, err := Open(filepath.Join(wd, tt.name+".psd"))
			if err != nil {
				t.Fatal(err)
			}

			// Dump the contents.
			want.Dump()
			// Grab a new version of the doc (loaded from json).
			got, err := Open(filepath.Join(wd, tt.name+".psd"))
			if err != nil {
				t.Fatal(err)
			}
			got.layerSets[0].current = true
			if !reflect.DeepEqual(got, want) {
				t.Errorf("wanted: %+v\ngot: %+v", want, got)
			}
		})
	}
}

func TestDocument_Save(t *testing.T) {
	file := filepath.Join(wd, "Test.psd")
	d, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}

	layerName := "Group 1"
	lyr := d.LayerSet(layerName)
	if lyr == nil {
		t.Fatalf("LayerSet '%s' was not found", layerName)
	}

	// Change a layer name.
	_, err = DoJS(filepath.Join("SetName"), JSLayer(lyr.Path()), "Group 2")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err = DoAction("DK", "Undo"); err != nil {
			t.Error(err)
		}
		err = d.Save()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = d.Save()
	if err != nil {
		t.Fatal(err)
	}

	d2, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(d, d2) {
		t.Errorf("wanted: %+v\ngot: %+v", d, d2)
	}
}
