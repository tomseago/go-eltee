DmxMonitor = new (function() {

    this.dmxMonStart = function reqFixtures() {
        var msg = {
            Msg: "dmxMonStart"
        }

        Connection.sendJSON(msg);
    }

    this.dmxMonStop = function reqFixtures() {
        var msg = {
            Msg: "dmxMonStop"
        }

        Connection.sendJSON(msg);
    }

    this.registerHandlers = function registerHandlers() {
        Connection.registerHandler("dmxData", this.handleDmxData.bind(this))
    }

    ////
    this.handleDmxData = function handleDmxData(type, body) {
        Log.info("handleDmxData---------- start");

        for (var i=0; i<body.length; i++) {
            var el = $("#dmxVal"+i);
            el.text(body[i]);
        }
    }


    ////

    this.start = function start() {
        this.registerHandlers();

        //this.dmxMonStart();
    }

})();