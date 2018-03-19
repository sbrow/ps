#include lib.js
var stdout = newFile(arguments[0]);
var lyr = eval(arguments[1]);
lyr.translate((Number)(arguments[2]), (Number)(arguments[3]));
stdout.writeln('{' + bounds(lyr) + '}')
stdout.close();