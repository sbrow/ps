#include lib.js
var stdout = newFile(arguments[0])
eval(arguments[1]);
stdout.close()