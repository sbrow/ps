package ps

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// TextItem holds the text element of a TextLayer.
type TextItem struct {
	contents string
	size     float64
	font     string
	parent   *ArtLayer
}

// TextItemJSON is the exported version of TextItem
// that allows it to be marshaled and unmarshaled
// into JSON.
type TextItemJSON struct {
	Contents string
	Size     float64
	Font     string
}

// Contents returns the raw text of the TextItem.
func (t TextItem) Contents() string {
	return t.contents
}

// Size returns the font size of the TextItem
func (t TextItem) Size() float64 {
	return t.size
}

// MarshalJSON implements the json.Marshaler interface, allowing the TextItem to be
// saved to disk in JSON format.
func (t *TextItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&TextItemJSON{
		Contents: t.contents,
		Size:     t.size,
		Font:     t.font,
	})
}

// UnmarshalJSON loads the JSON data into the TextItem
func (t *TextItem) UnmarshalJSON(data []byte) error {
	tmp := &TextItemJSON{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	t.contents = tmp.Contents
	t.size = tmp.Size
	t.font = tmp.Font
	return nil
}

// SetText sets the text to the given string.
func (t *TextItem) SetText(txt string) {
	if txt == t.contents {
		return
	}
	var err error
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	bndtext := "[[' + lyr.bounds[0] + ',' + lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + lyr.bounds[3] + ']]"
	js := fmt.Sprintf(`%s.textItem.contents='%s';var lyr = %[1]s;stdout.writeln(('%[3]s').replace(/ px/g, ''));`,
		lyr, txt, bndtext)
	var byt []byte
	if byt, err = DoJS("compilejs.jsx", js); err != nil {
		log.Panic(err)
	}
	var bnds *[2][2]int
	err = json.Unmarshal(byt, &bnds)
	if err != nil || bnds == nil {
		log.Println("text:", txt)
		log.Println("js:", js)
		fmt.Printf("byt: '%s'\n", string(byt))
		log.Panic(err)
	}
	t.contents = txt
	t.parent.bounds = *bnds
}

// SetSize sets the size of the TextItem's font.
func (t *TextItem) SetSize(s float64) {
	if t.size == s {
		return
	}
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	js := fmt.Sprintf("%s.textItem.size=%f;", lyr, s)
	_, err := DoJS("compilejs.jsx", js)
	if err != nil {
		t.size = s
	}
}

// Fmt applies the given font and style to all characters
// in the range [start, end].
func (t *TextItem) Fmt(start, end int, font, style string) {
	if !t.parent.Visible() {
		return
	}
	_, err := DoJS("fmtText.jsx", fmt.Sprint(start), fmt.Sprint(end), font, style)
	if err != nil {
		log.Panic(err)
	}
}
