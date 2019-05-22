import React, { Component, Fragment } from "react";
import Log from "../../lib/logger";

class DebugPage extends Component {
    state = {}

    render() {
        Log.info("DebugPage render");
        
        return (
            <Fragment>
                <h1>I am the debug page</h1>
            </Fragment>
        );
    }
}


export default DebugPage;
