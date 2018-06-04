if(app.activeDocument.activeLayer.kind == LayerKind.TEXT){
    var activeLayer = app.activeDocument.activeLayer;
    if(activeLayer.kind == LayerKind.TEXT){
        var start = parseInt(arguments[1]);
        var end = parseInt(arguments[2]);
        var fontName = arguments[3];
        var fontStyle = arguments[4];
        var fontSize = activeLayer.textItem.size;
        var colorArray = [0, 0, 0];
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
