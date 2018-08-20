#include lib.js
alert(app.activeDocument.path)
var f = newFile(arguments[0]);
for (var i = 0; i < arguments.length; i++) {
	f.writeln(arguments[i]);
}
f.close();