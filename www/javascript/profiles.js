Profiles = new (function() {

    this.reqProfiles = function reqProfiles() {
        var msg = {
            Msg: "reqProfiles"
        }

        Connection.sendJSON(msg);
    }

    this.registerHandlers = function registerHandlers() {
        Connection.registerHandler("profileList", this.handleProfileList.bind(this))
    }

    ////
    this.handleProfileList = function handleProfileList(type, body) {
        Log.info("Profile List---------- start");
        this.profiles = {};

        for(key in body) {
            Log.info(key, "->", body[key]);

            this.profiles[key] = body[key];
        }
        Log.info("----------------------- end");
    }

    ////

    this.start = function start() {
        this.registerHandlers();

        this.reqProfiles();
    }


})();