package ps

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type TextItem struct {
	contents string
	size     float64
	// color    Color
	font   string
	parent *ArtLayer
}

type TextItemJSON struct {
	Contents string
	Size     float64
	// Color    [3]int
	Font string
}

func (t *TextItem) Contents() string {
	return t.contents
}

func (t *TextItem) Size() float64 {
	return t.size
}

// MarshalJSON implements the json.Marshaler interface, allowing the TextItem to be
// saved to disk in JSON format.
func (t *TextItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&TextItemJSON{
		Contents: t.contents,
		Size:     t.size,
		// Color:    t.color.RGB(),
		Font: t.font,
	})
}

func (t *TextItem) UnmarshalJSON(b []byte) error {
	tmp := &TextItemJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	t.contents = tmp.Contents
	t.size = tmp.Size
	// t.color = RGB{tmp.Color[0], tmp.Color[1], tmp.Color[2]}
	t.font = tmp.Font
	return nil
}

func (t *TextItem) SetText(txt string) {
	if txt == t.contents {
		return
	}
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	bndtext := "[[' + lyr.bounds[0] + ',' + lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + lyr.bounds[3] + ']]"
	js := fmt.Sprintf(`%s.textItem.contents='%s';var lyr = %[1]s;stdout.writeln(('%[3]s').replace(/ px/g, ''));`,
		lyr, txt, bndtext)
	byt, err := DoJs("compilejs.jsx", js)
	var bnds *[2][2]int
	json.Unmarshal(byt, &bnds)
	if err != nil || bnds == nil {
		log.Println("text:", txt)
		log.Println("js:", js)
		fmt.Printf("byt: '%s'\n", string(byt))
		log.Panic(err)
	}
	t.contents = txt
	t.parent.bounds = *bnds
}

func (t *TextItem) SetSize(s float64) {
	if t.size == s {
		return
	}
	lyr := strings.TrimRight(JSLayer(t.parent.Path()), ";")
	js := fmt.Sprintf("%s.textItem.size=%f;", lyr, s)
	_, err := DoJs("compilejs.jsx", js)
	if err != nil {
		t.size = s
	}
}

// TODO: Documentation for Format(), make to textItem
func (t *TextItem) Fmt(start, end int, font, style string) {
	var err error
	if !t.parent.Visible() {
		return
	}
	_, err = DoJs("fmtText.jsx", fmt.Sprint(start), fmt.Sprint(end),
		font, style)
	if err != nil {
		log.Panic(err)
	}
}
