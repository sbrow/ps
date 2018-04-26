#include lib.js
var stdout = newFile(arguments[0]);
var doc = app.activeDocument;
stdout.writeln(('{"Name": "'+doc.name+'", "Height":'+doc.height+
    ', "Width":'+doc.width+', "ArtLayers": [').replace(/ px/g, ""));

stdout.writeln(layers(doc.artLayers))
stdout.writeln('], "LayerSets": [');
function lyrSets(sets, nm) {
	if (typeof sets === 'undefined')
		return;
	for (var i = 0; i < sets.length; i++) {
		var set = sets[i];
		var name = nm + set.name + "/";
		stdout.write('{"Name": "' + set.name + '", "Visible":'+ set.visible +'}');
		if (i+1 != sets.length)
			stdout.write(',');
	}
}
lyrSets(doc.layerSets)
stdout.write(']}');
stdout.close();