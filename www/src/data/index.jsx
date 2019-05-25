import { createStore, applyMiddleware } from "redux";

import Log from "../lib/logger";
import reducers, { initialState } from "./reducers";
import { apiCallStart, apiCallStarted, apiCallFailed, apiCallData, maybeDoOp } from "./actions";

import { client } from "../api";

const CACHE_TIME = 1000;

const apiAdapter = ({ getState, dispatch }) => next => (action) => {
    Log.info("apiAdapter action=", action);
    switch (action.type) {
    case apiCallStart.type:
        {
            let call = action.call;
            const ix = call.indexOf(":");

            if (ix !== -1) {
                call = call.substring(0, ix);
            }
            Log.info("apiCallStart action=", action);        
            if (!client[call]) {
                Log.error("Could not find an api method named ", call);
                return null;
            }

            // This counts as a start 
            dispatch(apiCallStarted(action.call, action.req));
            client[call](action.req, null)
                .then((resp) => {
                    Log.info("Got a resp ", action.call, resp);
                    dispatch(apiCallData(action.call, action.req, resp));
                })
                .catch((err) => {
                    Log.error("Error from api call: ", action.call, err);
                    dispatch(apiCallFailed(action.call, action.req, err));
                });
        }
        break;

    case maybeDoOp.type:
        {
            const { op } = action;
            if (!op) {
                Log.error("maybeDoOp action with no op value");
            } else if (op.isLoading) {
                Log.info(`Op ${op.call} is already loading`);
            } else if (Date.now() - op.updatedAt < CACHE_TIME) {
                Log.info(`Op ${op.call} is to recent: ${Date.now() - op.updatedAt}ms ago`);                
            } else {
                // Seems like we can allow it
                dispatch(apiCallStart(op.call, action.req));
            }
        }
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

