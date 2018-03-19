#include lib.js

var stdout = newFile(arguments[0]);
var set = eval(arguments[1]);
stdout.writeln('{"Name": "'+set.name+'", "ArtLayers":[');
for (var i = 0; i < set.artLayers.length; i++) {
	var lyr = set.artLayers[i];
	stdout.write(('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
		lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
	    lyr.bounds[3] + ']], "Visible": ' + lyr.visible + '}').replace(/ px/g, ""));
	if (i != set.artLayers.length - 1)
		stdout.writeln(",");
}
stdout.write("]}")
stdout.close();