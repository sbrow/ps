var saveFile = File(arguments[0]);
if(saveFile.exists)
    saveFile.remove();

saveFile.encoding = "UTF8";
saveFile.open("e", "TEXT", "????");
for (var i = 0; i < arguments.length; i++) {
	saveFile.writeln(arguments[i])
}
saveFile.close();