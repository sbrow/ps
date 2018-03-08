package ps

// type layer interface {
// 	Name() string
// 	TextItem() []string
// }

type ArtLayer struct {
	Name     string
	TextItem string
	Bounds   [2][2]int
	LayerSet string
}

func (a *ArtLayer) SetVisible() {
	DoJs("setLayerVisibility.jsx", a.LayerSet+"/"+a.Name, "true")
}

// func (a *ArtLayer) Name() string {
// 	return a.name
// }

// func (a *ArtLayer) TextItem() string {
// 	return a.textItem
// }
