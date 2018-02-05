var Path = arguments[0];
alert(Path)
var saveFile = File(Path + "/test.txt");

if(saveFile.exists)
    saveFile.remove();

saveFile.encoding = "UTF8";
saveFile.open("e", "TEXT", "????");
saveFile.writeln("Testing...");
saveFile.close();