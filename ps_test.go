package ps

import (
	"fmt"
	_ "io/ioutil"
	"os"
	"path"
	_ "strings"
	"testing"
)

func TestPkgPath(t *testing.T) {
	fmt.Println(pkgpath)
}

// TODO: Comparison borked
func TestRun(t *testing.T) {
	out := []byte("Testing...\n")
	msg, err := run("test")
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) == string(out) {
		fail := fmt.Sprintf("run(test)\ngot:\t\"%s\"\nwant:\t\"%s\"\n", msg, out)
		t.Fatal(fail)
	} else {
		os.Remove(path.Join(pkgpath, "scripts", "test.txt"))
	}
}

func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestOpen\"")
	}
	Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
}

func TestQuit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestQuit\"")
	}
	Quit(2)
}

func TestWait(t *testing.T) {
	Wait("Waiting...")
}

// TODO: Comparison borked
/*
func TestJS(t *testing.T) {
	out := "Testing...\n"
	_, err := Js(path.Join(Folder, "test.jsx"), Folder)
	if err != nil {
		t.Fatal(err)
	}
	f, err := ioutil.ReadFile(path.Join(Folder, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(string(f), string(out)) != 0 {
		fmt.Println(f)
		fmt.Println([]byte(out))
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", f, out)
		t.Fatal(fail)
	}
}
*/
