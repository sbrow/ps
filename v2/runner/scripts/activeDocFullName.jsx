#include lib.js
var stdout = newFile(arguments[0]);stdout.writeln(app.activeDocument.fullName);stdout.close();