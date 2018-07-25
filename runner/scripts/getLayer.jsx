#include lib.js
app.displayDialogs=DialogModes.NO
var stdout = newFile(arguments[0]);
var lyr = eval(arguments[1]);
var lyrs = [lyr];
stdout.writeln(layers(lyrs))