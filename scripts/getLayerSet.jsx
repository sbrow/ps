#include lib.js
var stdout = newFile(arguments[0]);
var set = eval(arguments[1]);
stdout.writeln('{"Name": "'+set.name+'", "Visible": '+ set.visible +', "ArtLayers":[');
stdout.flush();
for (var i = 0; i < set.artLayers.length; i++) {
	var lyr = set.artLayers[i];
	stdout.write(('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
		lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
	    lyr.bounds[3] + ']], "Visible": ' + lyr.visible + ',"Text":').replace(/ px/g, ""));
	if (lyr.kind == LayerKind.TEXT)
		stdout.write('"'+lyr.textItem.contents.replace(/\r/g, "\\r")+'"');
	else
		stdout.write("null");
	stdout.write("}")
	if (i != set.artLayers.length - 1)
		stdout.writeln(",");
	stdout.flush();
}
stdout.writeln("]");
stdout.write(', "LayerSets": [')
for (var i = 0; i < set.layerSets.length; i++) {
	var s = set.layerSets[i];
	stdout.write('{"Name":"' + s.name + '", "Visible": ' + s.visible  + '}');
	if (i < set.layerSets.length - 1)
		stdout.writeln(",");
	stdout.flush()
}
stdout.writeln(']')
// app.activeDocument.activeLayer=set;
// set.merge();
// set=eval(arguments[2]);
stdout.write(', "Bounds": [[],[]]');
// stdout.write((', "Bounds": [[' + set.bounds[0] + ',' +
		// set.bounds[1] + '],[' + set.bounds[2] + ',' + 
	    // set.bounds[3] + ']]').replace(/ px/g, ""));
stdout.write("}");
// Undo();
stdout.close();