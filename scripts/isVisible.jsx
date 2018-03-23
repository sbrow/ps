#include lib.js
var stdout = newFile(arguments[0]);stdout.writeln(eval(arguments[1]).visible);stdout.close();