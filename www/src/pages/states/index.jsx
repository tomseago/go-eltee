import React, { Component, Fragment } from "react";
import Log from "../../lib/logger";

import { StateNames } from "../../api";

class StatesPage extends Component {
    state = {}

    render() {
        Log.info("StatesPage render");
        
        return (
            <Fragment>
                <h1>States are...</h1>
                <StateNames>
                    { (loading, data, error) => {
                        if (loading) {
                            return <h2>Loading...</h2>;
                        }

                        if (error) {
                            return <h2>{error}</h2>;
                        }

                        console.log(data);
                        if (!data) {
                            return <h2>No data</h2>;
                        }

                        const items = data.map(name => (<li key={name}>{name}</li>));
                        return (
                            <ul>{items}</ul>
                        );
                    }}
                </StateNames>
            </Fragment>
        );
    }
}


export default StatesPage;
