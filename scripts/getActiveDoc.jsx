#include lib.js
var stdout = newFile(arguments[0]);
var doc = app.activeDocument;
stdout.writeln(('{"Name": "'+doc.name+'", "Height":'+doc.height+
    ', "Width":'+doc.width+', "ArtLayers": [').replace(/ px/g, ""));
function layers(lyrs) {
	if (typeof lyrs === 'undefined')
		return;
	for (var i = 0; i < lyrs.length; i++) {
		var lyr = lyrs[i];
		stdout.write(('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
	                     lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
	                     lyr.bounds[3] + ']], "Visible": ' + lyr.visible+', "Text":').replace(/ px/g, ""));
	if (lyr.kind == LayerKind.TEXT)
		stdout.write('"'+lyr.textItem.contents+'"');
	else
		stdout.write("null");
	stdout.write("}")
	if (i+1 != lyrs.length)
		stdout.write(',');
	stdout.writeln();
	}
}
layers(doc.artLayers)
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