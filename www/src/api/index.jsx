import React from "react";

import { ElTeeClient } from "./api_grpc_web_pb";
import pb from "./api_pb";



const client = new ElTeeClient("http://localhost:9090", null, null);

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
        const req = new pb.Void();
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
