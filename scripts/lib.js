// Opens and returns a file, overwriting new data.
function newFile(path) {
	var f = File(path)
	if(f.exists)
	    f.remove()
	f.encoding = "UTF8"
	f.open("e", "TEXT", "????")
	return f
}

// Moves a layer
function positionLayer(lyr, x, y, alignment){
     if(lyr.iisBackgroundLayer||lyr.positionLocked) return  
     var layerBounds = lyr.bounds;  
     var layerX = layerBounds[0].value;  
     if (alignment == 'top' || alignment == null)
         var layerY = layerBounds[1].value;  
     else if (alignment == 'bottom')
        var layerY = layerBounds[3].value;
     var deltaX = x-layerX;  
     var deltaY = y-layerY;  
     lyr.translate(deltaX, deltaY);  
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