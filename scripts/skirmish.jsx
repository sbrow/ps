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

DeckCardPSD.prototype.formatText = function() {
    var speed = this.textLayers.getByName('speed');
    if (speed.visible) {
        this.changeStroke(speed, (speed.textItem.contents == 1) ? [128, 128, 128] : [255, 255, 255],
                          this.colors.banner)
    }
    /**
     * The lowest we allow a text layer to go.
     * @type {int}
     */
    var bottom  = this.doc.height-this.tolerance.flavor_text

    // Get our text layers.
    var short_text = this.setTextLayer('short_text', undefined, null, 'Arial', 'Regular',[this.bold_words, "Bold"]);
    var long_text = this.textLayers.getByName('long_text');
    var flavor_text = this.textLayers.getByName('flavor_text');

    // Position the layers.
    positionLayer(this.short_textBackground, this.short_textBackground.bounds[0], short_text.bounds[3] + this.tolerance.short_text, 'bottom');
    positionLayer(long_text, long_text.bounds[0], this.short_textBackground.bounds[3] + this.tolerance.long_text, 'top');
    positionLayer(flavor_text, flavor_text.bounds[0], bottom, 'bottom');

    /**
     * Make our layers visible
     * @todo  hack, fix.
     */
    short_text.visible = short_text.textItem.contents != "“";
    long_text.visible = long_text.textItem.contents != "“";
    flavor_text.visible = flavor_text.textItem.contents != "“";



    //Hide long_text if too long.
    if (long_text.bounds[3] > this.doc.height - bottom) {
        long_text.visible == false;
    }

    this.log.debug(short_text.bounds)
    this.log.debug(long_text.bounds)
    this.log.debug(flavor_text.bounds)
    //Hide flavor text if too long.
    if ( (long_text.visible && flavor_text.bounds[1] < long_text.bounds[3])
    ||   (short_text.visible && flavor_text.bounds[1] < short_text.bounds[3])) {
        flavor_text.visible = false;
    }
};