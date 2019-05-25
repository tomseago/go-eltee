import Log from "../lib/logger";
import client from "./client";
import proto from "./api_pb";

import { apiCallStarted, apiCallData, apiCallFailed } from "../data/actions";


export function findApiOp(state, name) {
    if (!state.apiOps) {
        return { call: name };
    }

    return state.apiOps[name] || { call: name };
}

function isUpdateNeeded(state, call) {
    const opState = state.apiCalls[call];

    if (!opState) {
        return true;
    }

    if (opState.isLoading) {
        return false;
    }

    // As long as it's not loading, then sure???
    return true;
}

function makeCall(dispatch, call, req, handler) {
    if (!client[call]) {
        throw new Error(`Could not find an apiCall named ${call}`);
    }

    dispatch(apiCallStarted(call, req));    
    client[call](req, null)
        .then((resp) => {
            Log.info("Got a resp ", call, resp);

            handler(call, req, resp);
            dispatch(apiCallData(call, req, resp));
        })
        .catch((err) => {
            Log.error("Error from api call: ", call, err);
            dispatch(apiCallFailed(call, req, err));
        });
}

export function ensureStateNames(dispatch, state) {
    const call = "stateNames";

    if (!isUpdateNeeded(state, call)) {
        Log.info("Don't need to do ", call);
        return;
    }

    const req = new proto.Void();

    // We need to do it
    makeCall(dispatch, call, req, (call, req, resp) => {

    });
}

export function ensureControlPoints(dispatch, state, wsName) {

}
