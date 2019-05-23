import React, { Fragment } from "react";

import Button from "@material-ui/core/Button";

import Log from "../../lib/logger";
import Loading from "../../common/loading";
import ErrorComp, { ErrorBoundary } from "../../common/error";
import { StateNames, ApiCall } from "../../api";
import proto from "../../api/api_pb";

export default function StatesPage(props) {

    const req = new proto.Void();

    return (
        <ErrorBoundary>            
            <ApiCall name="stateNames" req={req}>
                { (loading, data, error, refresh) => {
                    // TODO: There is opportunity to add another
                    // composition component here to do the loading and
                    // error conditions that are pretty standard
                    if (loading) {
                        return <Loading />;
                    }

                    if (error) {
                        return <ErrorComp>{error}</ErrorComp>;
                    }

                    Log.info(data);
                    if (!data) {
                        return <ErrorComp>No data</ErrorComp>;
                    }

                    const items = data.listList.map(name => (<li key={name}>{name}</li>));

                    return (
                        <Fragment>
                            <ul>{items}</ul>
                            <Button variant="contained" onClick={() => refresh()}>Refresh</Button>
                        </Fragment>
                    );
                }}
            </ApiCall>
        </ErrorBoundary>
    );
}


