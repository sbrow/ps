#include lib.js

// var saveFile = File(arguments[0]);
var arg = 'app.activeDocument.layerSets.getByName("Text").artLayers.getByName("short");';
var set = eval(arg);
// set.textItem.size=10;
alert(set.textItem.font)
// var doc=app.activeDocument
// doc.layerSets.getByName("ResolveGem").merge();
// alert(doc.artLayers.getByName("ResolveGem").bounds); 
// doc.activeHistoryState=doc.historyStates[doc.historyStates.length-2]
// setFormatting(0,6, "Arial", "Bold");
// saveFile.close();