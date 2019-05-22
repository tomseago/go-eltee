/* eslint-disable no-empty */

class Logger {
    debug(...args) {
        try {
            console.debug(...args);
        } catch (e) { }
    }

    info(...args) {
        try {
            console.info(...args);
        } catch (e) { }
    }

    warn(...args) {
        try {
            console.warn(...args);
        } catch (e) { }
    }

    error(...args) {
        try {
            console.error(...args);
        } catch (e) { }
    }
}

const log = new Logger();

export default log;
