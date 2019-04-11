States = new (function() {

    this.reqCurrentState = function reqCurrentState() {
        var msg = {
            Msg: "reqCurrentState"
        }

        Connection.sendJSON(msg);
    }

    this.setCP = function setCP(cpName, value) {
        var msg = {
            Msg: "setCP",
            Body: {
                Name: cpName,
                Value: value,
            }
        }

        Connection.sendJSON(msg);        
    }

    this.registerHandlers = function registerHandlers() {
        Connection.registerHandler("currentState", this.handleCurrentState.bind(this))
    }

    ////
    this.handleCurrentState = function handleCurrentState(type, body) {
        Log.info("Current State ---------- start");

        Log.info(body);
        // var el = $("#fixtureList");
        // el.empty();

        // for(key in body) {
        //     Log.info(key, "->", body[key]);
        //     el.append("<div class='fixture'>"+key+" <a href='//dmx_override?fixture="+key+"'>DMX</a></div>");
        // }

        Log.info("----------------------- end");
    }

    ////

    this.start = function start() {
        this.registerHandlers();
    }

    /////
})();

class ControlPoint {
    constructor(name) {
        this.name = name;
    }

    toValue() {        
    }
}

class ColorPoint extends ControlPoint {
    constructor(name) {
        super(name);
        this.components = {};
    }

    toValue() {
        return this.components;
    }

    rgb(r,g,b) {
        this.components.red = parseFloat(r) / 255.0;
        this.components.green = parseFloat(g) / 255.0;
        this.components.blue = parseFloat(b) / 255.0;
    }
}

