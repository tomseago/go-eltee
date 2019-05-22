import React, { Component, Fragment } from "react";
import Log from "../../lib/logger";

class FixturesPage extends Component {
    state = {}

    render() {
        Log.info("FixturesPage render");
        
        return (
            <Fragment>
                <h1>Fixtures are...</h1>
            </Fragment>
        );
    }
}


export default FixturesPage;
