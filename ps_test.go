package ps

import (
	"fmt"
	"os"
	"path/filepath"
	_ "strings"
	"testing"
)

func TestPkgPath(t *testing.T) {
	out := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "sbrow", "ps")
	if filepath.Join(pkgpath) != out {
		t.Fatal(filepath.Join(pkgpath), out)
	}
}

func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestOpen\"")
	}
	_, err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	out := []byte("hello,\r\nworld!\r\n")
	msg, err := run("test", "hello,", "world!")
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != string(out) {
		fail := fmt.Sprintf("TestRun faild.\ngot:\n\"%s\"\nwant:\n\"%s\"\n", msg, out)
		t.Fatal(fail)
	}
}

func TestQuit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestQuit\"")
	}
	Quit(2)
}

func TestWait(t *testing.T) {
	Wait("Waiting...")
	fmt.Println()
}

func TestJS(t *testing.T) {
	out := []byte("F:\\TEMP\\js_out.txt\r\narg\r\nargs\r\n")
	script := "test.jsx"
	ret, err := Js(script, "arg", "args")
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != string(out) {
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", ret, out)
		t.Fatal(fail)
	}
}
