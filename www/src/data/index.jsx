import { createStore, applyMiddleware } from "redux";

import Log from "../lib/logger";
import reducers, { initialState } from "./reducers";
import { apiCallStart, apiCallStarted, apiCallFailed, apiCallData } from "./actions";

import { client } from "../api";

const apiAdapter = ({ getState, dispatch }) => next => (action) => {
    Log.info("apiAdapter action=", action);
    switch (action.type) {
    case apiCallStart.type:
        // This is a point at which we could check for an existing
        // in flight call, put we're not doing that right at this moment
        Log.info("apiCallStart action=", action);        
        if (!client[action.call]) {
            Log.error("Could not find an api method named ", action.call);
            return null;
        }

        // This counts as a start 
        dispatch(apiCallStarted(action.call, action.req));
        client[action.call](action.req, null)
            .then((resp) => {
                Log.info("Got a resp ", action.call, resp);
                dispatch(apiCallData(action.call, action.req, resp.toObject()));
            })
            .catch((err) => {
                Log.error("Error from api call: ", action.call, err);
                dispatch(apiCallFailed(action.call, action.req, err));
            });
        break;

    default:
        break;
    }

    const out = next(action);

    Log.info("apiAdapter action=", action, "  out=", out);
    return out;
};

const store = createStore(reducers, initialState, applyMiddleware(apiAdapter));
export default store;

