import { combineReducers } from "redux";

import Log from "../lib/logger";

import {apiCallStarted, apiCallFailed, apiCallData, handleCallResult} from "./actions";

// Our application state
export const initialState = {
    // A list of names of states
    stateNames: [],

    // Control points indexed by the state name
    cpsByState: {},

    // API calls that could be in flight, indexed by call name
    apiCalls: {},
};

function updateStateCps(existing = {}, action) {
    Log.info("Applying cpl updates from ", action);

    const wsName = action.req.getVal();
    Log.debug("wsName = ", wsName);

    existing[wsName] = action.data;
    Log.debug(existing);

    return existing;
}

// Doesn't need to be pure. Shit's already copied
function handleApiData(existing = {}, action) {
    const out = existing;

    let { call } = action;
    const ix = call.indexOf(":");
    if (ix !== -1) {
        call = call.substring(0, ix);
    }

    switch (call) {
    case "stateNames":
        out.stateNames = action.data.getListList();
        break;

    case "controlPoints":
        out.cpsByState = updateStateCps(out.cpsByState, action);
        break;

    default:
    }
}


export default function reducers(existing = initialState, action) {
    let out = existing;

    switch (action.type) {
    case apiCallStarted.type:
        out = {
            ...existing,
            apiCalls: {
                ...existing.apiCalls,
                [action.call]: {
                    call: action.call,
                    isLoading: true,
                    startedAt: action.at,
                },
            },
        };
        break;

    case apiCallData.type:
        {
            Log.info("existing.apiOps ", existing.apiCalls);
            const { startedAt } = existing.apiCalls[action.call];

            out = {
                ...existing,
                apiCalls: {
                    ...existing.apiCalls,
                    [action.call]: {
                        call: action.call,
                        isLoading: false,
                        elapsedTime: action.at - startedAt,
                        updatedAt: action.at,
                    },
                },
            };
            handleApiData(out, action);
        }
        break;

    case apiCallFailed.type:
        out = {
            ...existing,
            apiCalls: {
                ...existing.apiCalls,
                [action.call]: {
                    call: action.call,
                    isLoading: false,
                    lastError: action.error,
                },
            },
        };
        break;

    case handleCallResult.type:
        out = {
            ...existing,
        };
        if (action.handler) {
            action.handler(out, action.result, action.req, action.call);
        }
        break;

    default:
        // Keep out as it is
        break;
    }

    return out;
}

// ////
