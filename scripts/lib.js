// Opens and returns a file, overwriting new data.
function newFile(path) {
	var f = File(path)
	if(f.exists)
	    f.remove()
	f.encoding = "UTF8"
	f.open("e", "TEXT", "????")
	return f
}

// Returns an array of ArtLayers from a layerSet or an ArtLayer.
function getLayers(path) {
    try {
        var doc = app.activeDocument
        var path = path.split('/')
        var lyrs = doc.layerSets.getByName(path[0])
        for (var i = 1; i < path.length; i++) {
            try {
                lyrs = lyrs.layerSets.getByName(path[i])
            } catch (e) {
                lyrs = { layers: [lyrs.layers.getByName(path[i])] } }
        }
        return lyrs
    } catch (e) {
        if (e.message.indexOf('User') == -1)
            alert(err(e));
        else
            throw new Exception('User cancelled the operation');
    }
}

function err(e) {
    return 'ERROR: ' + e.message + ' at ' + e.fileName + ':' + e.line;
}