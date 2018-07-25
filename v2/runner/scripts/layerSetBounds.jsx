#include lib.js
var stdout = newFile(arguments[0]);
var set = eval(arguments[1]);
app.activeDocument.activeLayer=set;
set.merge();
set=eval(arguments[2]);
stdout.write(('[[' + set.bounds[0] + ',' +
		set.bounds[1] + '],[' + set.bounds[2] + ',' + 
	    set.bounds[3] + ']]').replace(/ px/g, ""));
Undo();
stdout.close();