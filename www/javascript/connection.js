Connection = new (function(){
    var CLOSED = 0;
    var OPENING = 1;
    var OPEN = 2;

    this.isOpen = CLOSED;
    this.path = "/websocket"

    this._socket = false
    this._sendQ = []

    this.sendText = function sendText(msg) {
        Log.debug("Queing msg", msg)

        this._sendQ.push(msg);
        this.trySend()
    }

    this.sendJSON = function sendJSON(msg) {
        msg = JSON.stringify(msg);
        this.sendText(msg);
    }

    this.trySend = function trySend() {
        Log.info("trySend() isOpen=", this.isOpen);
        if (this.isOpen != OPEN) {
            this.startOpen();
            return;
        }

        // We think it's open
        try {
            while(this._sendQ.length > 0) {
                var msg = this._sendQ.shift();
                Log.debug("Calling send for:", msg)
                this._socket.send(msg);
            }
        } catch (e) {
            Log.error("Failed to send message", e);
            this._sendQ = [];
            this.close();
        }
    }

    this.close = function close() {
        try {
            this._socket.close();
        } catch (e) {
            // ignore issues here            
        }
        this.isOpen = CLOSED;
    }

    this.startOpen = function startOpen() {
        if (this.isOpen != CLOSED) {
            return;
        }

        // Try to open the connection
        this.isOpen = OPENING;

        var u = new URL(document.location)        
        this._socket = new WebSocket("ws://"+u.host+this.path);

        this._socket.addEventListener("open", this.handleOpen.bind(this));
        this._socket.addEventListener("message", this.handleMessage.bind(this));
        this._socket.addEventListener("close", this.handleClose.bind(this));
        this._socket.addEventListener("error", this.handleError.bind(this));
    }

    this.handleOpen = function handleOpen(event) {
        Log.info("handleOpen", event);

        this.isOpen = OPEN;
        setTimeout(this.trySend.bind(this), 1);
    }

    this.handleMessage = function handleMessage(event) {
        Log.info("handleMessage", event);

        // TODO: Distinguish string versus binary messages and deal appropriately
        // For now, we assume everything is a JSON string

        try {
            var decoded = JSON.parse(event.data);
        } catch (e) {
            Log.error("Unable to parse message", e);
            return
        }

        var list = this._handlers[decoded.Msg];
        if (!list) {
            Log.warning("Message has no handlers, ignoring: ", decoded.Msg);
            return
        }

        for (var i=0; i<list.length; i++) {
            try {
                list[i](decoded.Msg, decoded.Body)
            } catch (e) {
                Log.error("Error with handler:", e)
            }
        }
    }

    this.handleClose = function handleClose(event) {
        Log.info("handleClose", event);
        this.isOpen = CLOSED;
        this._socket = false;
    }

    this.handleError = function handleError(event) {
        Log.error("handleError", event);

        this.isOpen = CLOSED;
        this._socket = false;
    }

    this._handlers = {}
    this.registerHandler = function registerHandler(msg, fn) {
        var list = this._handlers[msg];
        if (!list) {
            list = [];
            this._handlers[msg] = list;
        }

        list.push(fn);
    }

})();
