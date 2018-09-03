var saveFile = File(arguments[0]);
if(saveFile.exists)
    saveFile.remove();
var desc1 = new ActionDescriptor();
var ref1 = new ActionReference();
ref1.putName(stringIDToTypeID("dataSetClass"), arguments[1]);
desc1.putReference(charIDToTypeID("null"), ref1);
desc = executeAction(charIDToTypeID("Aply"), desc1, DialogModes.NO);
saveFile.encoding = "UTF8";
saveFile.open("e", "TEXT", "????");
saveFile.write("args: ");
for (i = 0; i < arguments.length; i++) {
    saveFile.write(arguments[i] + ",");
}
saveFile.writeln();
saveFile.close();