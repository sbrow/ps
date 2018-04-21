#include lib.js

// var saveFile = File(arguments[0]);
var arg = 'app.activeDocument.layerSets.getByName("ResolveGem");';
var set = eval(arg);
set.visible=false;
alert(set.visible)
// var doc=app.activeDocument
// doc.layerSets.getByName("ResolveGem").merge();
// alert(doc.artLayers.getByName("ResolveGem").bounds); 
// doc.activeHistoryState=doc.historyStates[doc.historyStates.length-2]
// setFormatting(0,6, "Arial", "Bold");
// saveFile.close();