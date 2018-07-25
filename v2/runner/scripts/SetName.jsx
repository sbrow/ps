#include lib.js
var stdout=newFile(arguments[0]);
var obj=eval(arguments[1]);
obj.name=arguments[2];
stdout.close();