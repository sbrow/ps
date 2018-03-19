#include lib.js

var stdout = newFile(arguments[0]);
var doc = app.activeDocument;

stdout.write(('{"Name": "' + doc.name +'", "Height":' +doc.height +
               ', "Width":' + doc.width + ", ").replace(/ px/g, ""));

function layersNsets(obj) {
	stdout.write('"ArtLayers": [');
	lyrss(obj.artLayers, "")
	stdout.write('], "LayerSets": [');
	lyrSets(obj.layerSets, "");
	// stdout.write('], "LayerSets": [');
}

function lyrss(lyrs, set) {
	if (typeof lyrs === 'undefined')
		return;
	for (var i = 0; i < lyrs.length; i++) {
		var lyr = lyrs[i];
		stdout.write(('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
	                     lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
	                     lyr.bounds[3] + ']], "Path": "' + set + 
	                     '", "Visible": ' + lyr.visible + '}').replace(/ px/g, ""));
		if (i+1 != lyrs.length)
			stdout.write(',');
	}
}
function lyrSets(sets, nm) {
	if (typeof sets === 'undefined')
		return;
	for (var i = 0; i < sets.length; i++) {
		var set = sets[i];
		var name = nm + set.name + "/";
		stdout.write('{"Name": "' + set.name + '", "LayerSets": [');
		lyrSets(set.layerSets, name);
		stdout.write('], "Layers": [');
		lyrss(set.artLayers, name);
		stdout.write(']}');
		if (i+1 != sets.length)
			stdout.write(',');
	}
}

layersNsets(doc)
stdout.writeln(']}');

alert(doc.layerSets.getByName("Group 2").layerSets.getByName("Group 1").layers.getByName("Layer 1").name)
stdout.close();