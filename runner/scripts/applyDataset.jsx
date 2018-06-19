var saveFile = File(arguments[0]);
if(saveFile.exists)
    saveFile.remove();
var idAply = charIDToTypeID("Aply");
var desc1 = new ActionDescriptor();
var idnull = charIDToTypeID("null");
var ref1 = new ActionReference();
var iddataSetClass = stringIDToTypeID("dataSetClass");
ref1.putName(iddataSetClass, arguments[1]);
desc1.putReference(idnull, ref1);
executeAction(idAply, desc1, DialogModes.NO);
saveFile.encoding = "UTF8";
saveFile.open("e", "TEXT", "????");
saveFile.writeln("done!");
saveFile.close();