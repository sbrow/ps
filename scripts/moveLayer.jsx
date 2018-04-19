#include lib.js
var stdout = newFile(arguments[0]);
var lyr = eval(arguments[1]);
lyr.translate((Number)(arguments[2]), (Number)(arguments[3]));
if (lyr.typename == 'LayerSet') {
	alert(lyr.name + "\n" + lyr.typename)
	alert(lyr)
	// lyr.merge()
	// lyr=eval(arguments[4])
	// Undo();
}
stdout.writeln('{' + bounds(lyr) + '}')
stdout.close();