/* eslint-disable react/no-multi-comp,react/prop-types */
import React, { useState, useEffect } from "react";

import Log from "../lib/logger";
import { ElTeePromiseClient } from "./api_grpc_web_pb";
import proto from "./api_pb";

// const client = new ElTeeClient("http://localhost:9090", null, null);
export const client = new ElTeePromiseClient("http://localhost:9090", null, null);

// Really client could/should be passed in from the outside right???


export const ApiContext = React.createContext(null);

export function ApiProvider(props) {
    const { children } = props;

    return (
        <ApiContext.Provider value={client}>
            {children}
        </ApiContext.Provider>
    );
}

export function withApi(Comp) {
    return props => (
        <ApiContext.Consumer>
            {api => <Comp {...props} api={api} />}
        </ApiContext.Consumer>
    );
}

export class StateNames extends React.Component {
    state = {};

    componentDidMount() {
        this.setState({ loading: true });
        const req = new proto.Void();
        client.stateNames(req, null, (error, sMsg) => {
            this.setState({ loading: false });
            if (error != null) {
                this.setState({ error });
                return;
            }

            this.setState({ data: sMsg.getListList() });
        });
    }

    render() {
        const { loading, data, error } = this.state;
        const { children } = this.props;
        const fn = children;

        if (!fn) {
            return null;
        }

        if (loading) {
            return fn(true);
        }

        return fn(false, data, error);
    }
}

export function ApiCall(props) {
    const { name, req: propReq, children } = props;

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [data, setData] = useState(null);
    const [req, setReq] = useState(propReq);
    const [serial, setSerial] = useState(0);

    useEffect(() => {
        setLoading(true);
        setError(false);

        Log.info("ApiCall name=", name, " req=", req.toObject(), " serial=", serial);
        client[name](req, null)
            .then((resp) => {
                setData(resp.toObject());
                setLoading(false);
            })
            .catch((err) => {
                setError(err);
                setLoading(false);
            });

        // Should there be a cleanup? I don't really think so...
    }, [req, serial]);

    const fn = children;

    // Now let's do our rendering yo!
    if (!fn) {
        return null;
    }

    if (loading) {
        return fn(true);
    }

    return fn(false, data, error, (newReq) => {
        if (newReq) {
            setReq(newReq);
        }

        // Regardless of it being new, we can update it
        setSerial(serial + 1);
    });
}
