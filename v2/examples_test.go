package ps

import "fmt"

func ExampleJSLayer() {
	// The path of a layer inside a top level group.
	path := "Group 1/Layer 1"
	fmt.Println(JSLayer(path))
	// Output: app.activeDocument.layerSets.getByName('Group 1').artLayers.getByName('Layer 1')
}
