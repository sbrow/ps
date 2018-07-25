#include lib.js
var stdout = newFile(arguments[0]);
var set = eval(arguments[1]);
stdout.writeln('{"Name": "'+set.name+'", "Visible": '+ set.visible +', "ArtLayers":[');
stdout.flush();
var str = layers(set.artLayers);
str = str.replace(/\r/g, "\\r");
stdout.writeln(str);
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
stdout.write(', "Bounds": [[],[]]');
stdout.write("}");
stdout.close();