import Log from "../lib/logger";

import {doCall, maybeDoCall, apiCallStarted, apiCallFailed, apiCallData, handleCallResult} from "./actions";
import { client } from "../api";

// The apiMiddleware allows maybeDoOp actions to be dispatched which
// then asynchronously cause API calls to happen

export const apiMiddleware = store => next => action => {
    if (action.type !== maybeDoCall.type && action.type !== doCall.type) {
        next(action);
        return;
    }

    if (!action.call) {
        Log.error("maybeDoCall action has no call ", action);
        return;
    }

    if (action.type === doCall.type) {
        const rpcName = action.call;

        if (!client[rpcName]) {
            Log.error(`Could not find an apiCall named ${rpcName}`);
            return;
        }

        client[rpcName](action.req, null)
            .then((resp) => {
                next(handleCallResult(action.call, action.req, resp, action.success));
            })
            .catch((err) => {
                Log.error("Error from api call:", action, err);
                next(handleCallResult(action.call, action.req, err, action.failure));
            });

    } else {
        const {apiCalls} = store.getState();

        if (apiCalls[action.call] && apiCalls[action.call].isLoading) {
            // Nothing to do
            Log.info("API call is already loading ", action.op);
            return;
        }

        let rpcName = action.call;
        const ix = rpcName.indexOf(":");
        if (ix != -1) {
            rpcName = rpcName.substring(0, ix);
        }

        if (!client[rpcName]) {
            Log.error(`Could not find an apiCall named ${rpcName}`);
            return;
        }

        // Tell the world we're going to start it
        next(apiCallStarted(action.call));

        client[rpcName](action.req, null)
            .then((resp) => {
                Log.info("Got a resp ", action.call, resp);

                next(apiCallData(action.call, action.req, resp));
            })
            .catch((err) => {
                Log.error("Error from api call: ", action.call, err);
                next(apiCallFailed(action.call, action.req, err));
            });
    }

};

