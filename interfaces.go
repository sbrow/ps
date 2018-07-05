package ps

// Group represents a Document or LayerSet.
type Group interface {
	Name() string
	Parent() Group
	SetParent(Group)
	Path() string
	ArtLayer(name string) *ArtLayer
	LayerSet(name string) *LayerSet
	ArtLayers() []*ArtLayer
	LayerSets() []*LayerSet
	MustExist(name string) Layer
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

// Layer represents an ArtLayer or LayerSet Object.
type Layer interface {
	Bounds() [2][2]int
	MarshalJSON() ([]byte, error)
	Name() string
	Parent() Group
	Path() string
	Refresh() error
	SetParent(g Group)
	SetPos(x, y int, bound string)
	SetVisible(b bool) error
	UnmarshalJSON(b []byte) error
	Visible() bool
}
