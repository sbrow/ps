#target Photshop

try {
    /**
     * Make window
     * @todo read data from the template PSD's data sets.
     * @type {Window}
     */
    var win = new Window('dialog','Sync Skirmish Cards');  
    var panel = win.add('tabbedpanel');

    //Make main tab tab that updates folders at a time.
    var main = panel.add('tab', undefined, 'Main');

    //Get our list of deck names.
    var deckNames = (JSON_FOLDER.toString().replace(/\/[^,]*\//g, '').replace(/.json/g, '')).split(',');
    var decks = [];
    var deckJSONs = {};
    var index = 0;
    for (var i in deckNames) {
        if (deckNames[i] != 'Formatting' && deckNames[i] != 'desktop.ini' && deckNames[i] != 'old') {
            var filePath = '{0}/{1}.json'.format(JSON_PATH, deckNames[i]);
            var fileObj = loadJSON(filePath, true);
            deckJSONs[deckNames[i]] = fileObj;
            decks[index++] = deckNames[i];
            var deckTab = panel.add('tab', undefined, deckNames[i]);
            for (var card in deckJSONs[deckNames[i]]) {                
                deckTab.add('checkbox', undefined, deckJSONs[deckNames[i]][card].Name);
            }
        }
    }
/*
    //Button to select everything.
    main.everything = main.add('checkbox', undefined, 'Everything');
    main.everything.onClick = function() {
        main.leaders.value = main.everything.value;
        main.decks.allDecks.value = main.everything.value;
        main.decks.allDecks.onClick();
    };*/

    //Button to select all leaders.
    //main.leaders = main.add('checkbox', undefined, 'Leaders')

    //Section listing all the available decks
    main.decks = main.add('panel', undefined, "Decks"); 

    //Button that selects all Decks.
    main.decks.allDecks = main.decks.add('checkbox', undefined, 'All');

    //Generates checkbox for each deck.
    for (var deck in decks)
        var deckBox = main.decks.add('checkbox', undefined, decks[deck]);

    //Generates a checkbox for each card.
    for (var i=1; i<decks.length+1; i++) {
        main.decks.children[i].onClick = function() {
            for (var k=1; k<main.decks.children.length; k++)
                if (main.decks.children[k].text == panel.children[k].text)
                    for (var j=0; j<panel.children[k].children.length; j++)
                        panel.children[k].children[j].value = main.decks.children[k].value;
        };
    }

    main.decks.allDecks.onClick = function() {
        for (var box=1; box < decks.length+1; box++) { 
            main.decks.children[box].value = main.decks.allDecks.value;
            main.decks.children[box].onClick();
        }
     }

    //Button to cancel the operation.
    main.cancel = main.add('Button', undefined, 'Cancel');

    //Button to compile the chosen cards.
    main.compile = main.add ('Button',undefined, 'Compile');
    main.compile.onClick = function() {
        LOG.debug('', $.fileName);

        var ret = {};
        var leaders = HEROES_JSON;
        for (var i=1; i<panel.children.length; i++) {
            if (i != 2) {
            ret[decks[i-1]] = deckJSONs[decks[i-1]];
            for (var j=panel.children[i].children.length-1; j>=0; j--)
                if (panel.children[i].children[j].value == false)
                    ret[decks[i-1]].splice(j, 1);
            }
        }
        for (var i=panel.children[2].children.length-1; i>=0; i--) {
            if (panel.children[2].children[i].value == false)
                leaders.splice(i, 1);
        }
        if (isEmpty(ret) &&  isEmpty(leaders))
            alert('Must select something');
        else {
            win.close();
            compileDecks(ret, leaders);
        }
    }  

    win.center();  

    function compileDecks(decks, leaders) {
        LOG.log('SYNCING SKIRMISH FILES...', '-'); 
        if (decks)
            try {
                var t0 = new Date().getTime()
                sync(decks, true, LOG);
                var t1 = new Date().getTime()
                LOG.log("TIME: " + (t1 - t0))
            } catch (e) {
                LOG.err(e);
            }
        if(leaders) {
            syncLeaders(leaders, LOG);
        }
    }
    win.show();
} catch (e) {
    alert(e.LineNumber + " " + e.message);
    LOG.err(e);
    }
