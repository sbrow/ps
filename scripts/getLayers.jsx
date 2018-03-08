//arguments = [ 'F:\\TEMP\\js_out.txt\r\narg\r\nargs\r\n', 'Areas/TitleBackground']
var saveFile = File(arguments[0])
if(saveFile.exists)
    saveFile.remove()
saveFile.encoding = "UTF8"
saveFile.open("e", "TEXT", "????")
try {
    var doc = app.activeDocument
    var splitPath = arguments[1].split('/')
    var bottomLayerSet = doc.layerSets.getByName(splitPath[0])
    for (var i = 1; i < splitPath.length; i++) {
        try {bottomLayerSet = bottomLayerSet.layerSets.getByName(splitPath[i])}
        catch (e) {bottomLayerSet = { layers: [bottomLayerSet.layers.getByName(splitPath[i])] }}
    }
    saveFile.writeln('['); 
    for (var l = 0; l < bottomLayerSet.layers.length; l++) {
        var lyr = bottomLayerSet.layers[l]
        saveFile.write('{"Name":"' + lyr.name + '", "Bounds": [["' + lyr.bounds[0] + '","' +
                         lyr.bounds[1] + '","' + lyr.bounds[2] + '"],["' + lyr.bounds[3] + '"]]}'); 
        if (l != bottomLayerSet.layers.length - 1)
            saveFile.write(',');
        saveFile.writeln();

    }
    saveFile.writeln(']'); 
} catch (e) {
    if (e.message.indexOf('User') == -1)
        alert('ERROR: ' + e.message + ' at ' +  e.fileName + ':' + e.line);
    else
        throw new Exception('User cancelled the operation');
}
saveFile.close();