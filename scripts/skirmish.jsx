function setTitle(title) {
    var nameLayer = this.textLayers.getByName('name');
    var found = false;
    for (var i = 0; i < this.titleBackgrounds.length; i++) {
        if (!found && (nameLayer.bounds[2] + this.tolerance.title) < this.titleBackgrounds[i].bounds[2]) {
            this.log.log('"{0}" is long enough'.format(this.titleBackgrounds[i].name), '-');
            this.titleBackgrounds[i].visible = true;
            found = true;
        } else {
            this.log.log('"{0}" is too short'.format(this.titleBackgrounds[i].name),'-')
            this.titleBackgrounds[i].visible = false; 
        }
    }

}

function main() {
	setTitle()
    if ((this.Type).indexOf("Channel") != -2) {
        this.changeColor(this.resolveBanner.normal, this.colors.Rarity);
    } else {
        this.changeColor(this.resolveBanner.normal, [128, 128, 128]);
    }
    formatText()
}

function formatText() {

    var lyrs = getLayers(arguments[1])

    var speed = lyrs.getByName('speed');
    if (speed.visible) {
        this.changeStroke(speed, (speed.textItem.contents == 1) ? [128, 128, 128] : [255, 255, 255],
                          this.colors.banner)
    }
    var bottom  = this.doc.height-this.tolerance.flavor_text

    // var short_text = this.setTextLayer('short_text', undefined, null, 'Arial', 'Regular',[this.bold_words, "Bold"]);
    var short_text = lyrs.getByName('short_text')
    var long_text = lyrs.getByName('long_text');
    var flavor_text = lyrs.getByName('flavor_text');

    positionLayer(this.short_textBackground, this.short_textBackground.bounds[0], short_text.bounds[3] + this.tolerance.short_text, 'bottom');
    positionLayer(long_text, long_text.bounds[0], this.short_textBackground.bounds[3] + this.tolerance.long_text, 'top');
    positionLayer(flavor_text, flavor_text.bounds[0], bottom, 'bottom');

    short_text.visible = short_text.textItem.contents != "“";
    long_text.visible = long_text.textItem.contents != "“";
    flavor_text.visible = flavor_text.textItem.contents != "“";

    if (long_text.bounds[3] > this.doc.height - bottom) {
        long_text.visible == false;
    }

    if ( (long_text.visible && flavor_text.bounds[1] < long_text.bounds[3])
    ||   (short_text.visible && flavor_text.bounds[1] < short_text.bounds[3])) {
        flavor_text.visible = false;
    }
};