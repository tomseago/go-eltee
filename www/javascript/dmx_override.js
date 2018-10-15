DmxOverride = new (function() {

    // this.reqFixtures = function reqFixtures() {
    //     var msg = {
    //         Msg: "reqFixtures"
    //     }

    //     Connection.sendJSON(msg);
    // }

    this.registerHandlers = function registerHandlers() {
        // Connection.registerHandler("fixtureList", this.handleFixtureList.bind(this))
    }

    ////
    // this.handleFixtureList = function handleFixtureList(type, body) {
    //     Log.info("Fixture List---------- start");

    //     var el = $("#fixtureList");
    //     el.empty();

    //     for(key in body) {
    //         Log.info(key, "->", body[key]);
    //         el.append("<div class='fixture'>"+key+" <a href='//dmx_override?fixture="+key+"'>DMX</a></div>");
    //     }

    //     Log.info("----------------------- end");
    // }

    ////

    this.start = function start() {
        this.registerHandlers();
    }



})();