// Opens and returns a file, overwriting new data.
function newFile(path) {
	var f = File(path)
	if(f.exists)
	    f.remove()
	f.encoding = "UTF8"
	f.open("e", "TEXT", "????")
	return f
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