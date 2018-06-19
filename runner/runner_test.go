package runner

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var scripts string

func init() {
	var ok bool
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("package path not found")
	}
	pkgpath = filepath.Dir(file)
	scripts = filepath.Join(pkgpath, "scripts")
}
func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []byte
		wantErr bool
	}{
		{"dojs", []string{filepath.Join(scripts, "test.jsx"), filepath.Join(scripts, "test.txt"), "arg1", "arg2"},
			[]byte(filepath.Join(scripts, "test.txt") + "\r\narg1\r\narg2\r\n"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Run(tt.name, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			f, err := os.Open(tt.args[1])
			if err != nil {
				t.Error(err)
				return
			}
			got, err := ioutil.ReadAll(f)
			if err != nil {
				t.Error(err)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("Run() = \n%v, want \n%v", string(got), string(tt.want))
			}
		})
	}
}
