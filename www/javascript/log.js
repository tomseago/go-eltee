Log = new (function() { 
    this.buffered = []

    this.internal = function(level, args) {

        out = [level];
        for(var i=0; i<args.length; i++) {
            out.push(args[i]);
        }

        if (this.buffered.length > 10) {
            this.buffered.unshift();
        }
        this.buffered.push(out.join(" "));
    }

    this.error = function() {
        console.error.apply(console, arguments);
        this.internal("ERROR", arguments);
    }

    this.err = function() {
        console.error.apply(console, arguments)
        this.internal("ERROR", arguments)
    }

    this.warning = function() {
        console.warn.apply(console, arguments)
        this.internal(" WARN", arguments)
    }

    this.warn = function() {
        console.warn.apply(console, arguments)
        this.internal(" WARN", arguments)
    }

    this.info = function() {
        console.log.apply(console, arguments)
        this.internal(" INFO", arguments)
    }

    this.log = function() {
        console.log.apply(console, arguments)
        this.internal(" INFO", arguments)
    }

    this.notice = function() {
        console.log.apply(console, arguments)
        this.internal("NOTIC", arguments)
    }

    this.debug = function() {
        console.debug.apply(console, arguments)        
        this.internal("DEBUG", arguments)
    }

})()
