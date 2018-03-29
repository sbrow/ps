#include lib.js

var stdout = newFile(arguments[0]);
var lyr = eval(arguments[1]);
stdout.write(('{"Name":"' + lyr.name + '","Bounds":[[' + lyr.bounds[0] + ',' +
	lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
    lyr.bounds[3] + ']],"Visible":' + lyr.visible+',"Text":').replace(/ px/g, ""));
if (lyr.kind == LayerKind.TEXT) {
	stdout.write('"'+lyr.textItem.contents.replace(/\r/g, "\\r")+'"');
}
else
	stdout.write(null)
stdout.writeln('}')
stdout.close();