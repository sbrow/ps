// Opens and returns a file, overwriting new data.
function newFile(path) {
	var f = File(path)
	f.encoding = "UTF8"
	f.open("w")
	return f
}

File.prototype.flush = function() {
    this.close()
    this.open("a")
};
function flush(file) {
    file.close()
    file.open("a")
}


// Prints an error message.
function err(e) {
    return 'ERROR: ' + e.message + ' at ' + e.fileName + ':' + e.line;
}

function bounds(lyr) {
    return ('"Bounds": [[' + lyr.bounds[0] + ',' +
    lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
    lyr.bounds[3] + ']]').replace(/ px/g, "");
}

function Undo() {
    var desc = new ActionDescriptor();
    var ref = new ActionReference();
    ref.putEnumerated( charIDToTypeID( "HstS" ), charIDToTypeID( "Ordn" ), charIDToTypeID( "Prvs" ));
    desc.putReference(charIDToTypeID( "null" ), ref);
    executeAction( charIDToTypeID( "slct" ), desc, DialogModes.NO );
}

/**
* The setFormatting function sets the font, font style, point size, and RGB color of specified
* characters in a Photoshop text layer.
*
* @param start (int) the index of the insertion point *before* the character you want.,
* @param end (int) the index of the insertion point following the character.
* @param fontName is a string for the font name.
* @param fontStyle is a string for the font style.
* @param fontSize (Number) the point size of the text.
* @param colorArray (Array) is the RGB color to be applied to the text.
*/      
function setFormatting(start, end, fontName, fontStyle, fontSize, colorArray) {
    if(app.activeDocument.activeLayer.kind == LayerKind.TEXT){
        var activeLayer = app.activeDocument.activeLayer;
        fontSize = activeLayer.textItem.size;
        colorArray = [0, 0, 0];
        if(activeLayer.kind == LayerKind.TEXT){
            if((activeLayer.textItem.contents != "")&&(start >= 0)&&(end <= activeLayer.textItem.contents.length)){
                var idsetd = app.charIDToTypeID( "setd" );
                var action = new ActionDescriptor();
                var idnull = app.charIDToTypeID( "null" );
                var reference = new ActionReference();
                var idTxLr = app.charIDToTypeID( "TxLr" );
                var idOrdn = app.charIDToTypeID( "Ordn" );
                var idTrgt = app.charIDToTypeID( "Trgt" );
                reference.putEnumerated( idTxLr, idOrdn, idTrgt );
                action.putReference( idnull, reference );
                var idT = app.charIDToTypeID( "T   " );
                var textAction = new ActionDescriptor();
                var idTxtt = app.charIDToTypeID( "Txtt" );
                var actionList = new ActionList();
                var textRange = new ActionDescriptor();
                var idFrom = app.charIDToTypeID( "From" );
                textRange.putInteger( idFrom, start );
                textRange.putInteger( idT, end );
                var idTxtS = app.charIDToTypeID( "TxtS" );
                var formatting = new ActionDescriptor();
                var idFntN = app.charIDToTypeID( "FntN" );
                formatting.putString( idFntN, fontName );
                var idFntS = app.charIDToTypeID( "FntS" );
                formatting.putString( idFntS, fontStyle );
                var idSz = app.charIDToTypeID( "Sz  " );
                var idPnt = app.charIDToTypeID( "#Pnt" );
                formatting.putUnitDouble( idSz, idPnt, fontSize );
                var idClr = app.charIDToTypeID( "Clr " );
                var colorAction = new ActionDescriptor();
                var idRd = app.charIDToTypeID( "Rd  " );
                colorAction.putDouble( idRd, colorArray[0] );
                var idGrn = app.charIDToTypeID( "Grn " );
                colorAction.putDouble( idGrn, colorArray[1]);
                var idBl = app.charIDToTypeID( "Bl  " );
                colorAction.putDouble( idBl, colorArray[2] );
                var idRGBC = app.charIDToTypeID( "RGBC" );
                formatting.putObject( idClr, idRGBC, colorAction );
                textRange.putObject( idTxtS, idTxtS, formatting );
                actionList.putObject( idTxtt, textRange );
                textAction.putList( idTxtt, actionList );
                action.putObject( idT, idTxLr, textAction );
                app.executeAction( idsetd, action, DialogModes.NO );
            }
        }
    }
}

function layers(lyrs) {
    if (typeof lyrs === 'undefined')
        return;
    var str = ""; 
    for (var i = 0; i < lyrs.length; i++) {
        var lyr = lyrs[i];
        str += ('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
                         lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
                         lyr.bounds[3] + ']], "Visible": ' + lyr.visible+', "TextItem": ').replace(/ px/g, "");
        if (lyr.kind == LayerKind.TEXT) {
            str += ('{"Contents": "'+lyr.textItem.contents+'",').replace(/\r/g, '\\r');
            var size = Number(lyr.textItem.size.replace(/ p[tx]/g, ''));
            if (lyr.textItem.size.includes("px")) {
                size = (size / app.activeDocument.resolution) * 72;
            }
            str += ' "Size": '+size+',';
            str += ' "Font": "'+lyr.textItem.font+'"}\n'
        } else
            str += "null";
    str += "}";
    if (i+1 != lyrs.length)
        str += ',';
    }
    return str
}