#include lib.js
var stdout = newFile(arguments[0])
var lyrs = getLayers(arguments[1])
var vis = arguments[2] == "true"
try {
	for (var i = 0; i < lyrs.layers.length; i++)
	    lyrs.layers[i].visible = vis
} catch (e) {
	stdout.writeln(err(e))
}
stdout.close()