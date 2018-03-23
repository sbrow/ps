#include lib.js

var stdout = newFile(arguments[0]);
var lyr = eval(arguments[1]);
stdout.writeln(('{"Name":"' + lyr.name + '","Bounds":[[' + lyr.bounds[0] + ',' +
	lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
    lyr.bounds[3] + ']],"Visible":' + lyr.visible + '}').replace(/ px/g, ""));
stdout.close();