#include lib.js
var stdout = newFile(arguments[0])
var set = getLayers(arguments[1])

stdout.writeln('['); 
for (var l = 0; l < set.layers.length; l++) {
    var lyr = set.layers[l]
    var lyrset = arguments[1].replace(lyr.name, "")
    stdout.write(('{"Name":"' + lyr.name + '", "Bounds": [[' + lyr.bounds[0] + ',' +
                     lyr.bounds[1] + '],[' + lyr.bounds[2] + ',' + 
                     lyr.bounds[3] + ']], "LayerSet": "' + lyrset + '"}').replace(/ px/g, "")); 
    if (l != set.layers.length - 1)
        stdout.write(',');
    stdout.writeln();

}
stdout.writeln(']'); 
stdout.close()