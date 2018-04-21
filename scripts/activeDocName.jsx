#include lib.js
var stdout = newFile(arguments[0]);stdout.writeln(app.activeDocument.name);stdout.close();