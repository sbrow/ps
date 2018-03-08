package ps

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPkgPath(t *testing.T) {
	out := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "sbrow", "ps")
	if filepath.Join(pkgpath) != out {
		t.Fatal(filepath.Join(pkgpath), out)
	}
}

func TestStart(t *testing.T) {
	err := Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestOpen\"")
	}
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestClose\"")
	}
	err := Close(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestQuit\"")
	}
	err := Quit(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoJs(t *testing.T) {
	out := []byte("F:\\TEMP\\js_out.txt\r\narg\r\nargs\r\n")
	script := "test.jsx"
	ret, err := DoJs(script, "arg", "args")
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != string(out) {
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", ret, out)
		t.Fatal(fail)
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

func TestWait(t *testing.T) {
	Wait("Waiting...")
}

func TestDoAction_Crop(t *testing.T) {
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
	err = DoAction("DK", "Crop")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoAction_Undo(t *testing.T) {
	err := DoAction("DK", "Undo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveAs(t *testing.T) {
	err := SaveAs("F:\\TEMP\\test.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLayers(t *testing.T) {
	byt, err := Layers("Areas/TitleBackground")
	// _, err := Layers("Areas/TitleBackground")
	if err != nil {
		t.Fatal(err)
	}
	for _, lyr := range byt {
		fmt.Println(lyr.Name)
		fmt.Println(lyr.Bounds)
	}

}

func TestLayer(t *testing.T) {
	lyr, err := Layer("Areas/TitleBackground")
	// _, err := Layer("Areas/TitleBackground")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(lyr.Name)
	fmt.Println(lyr.Bounds)

}

func TestApplyDataset(t *testing.T) {
	out := []byte("done!\r\n")
	ret, err := ApplyDataset("Anger")
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != string(out) {
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", ret, out)
		t.Fatal(fail)
	}
	err = Quit(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoJs_HideLayer(t *testing.T) {
	_, err := DoJs("hideLayers.jsx", "Areas/TitleBackground")
	if err != nil {
		t.Fatal(err)
	}
}
