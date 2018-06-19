#include lib.js
var f = newFile(arguments[0]);
for (var i = 0; i < arguments.length; i++) {
	f.writeln(arguments[i]);
}
f.close();