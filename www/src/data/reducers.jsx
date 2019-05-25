import { combineReducers } from "redux";

import Log from "../lib/logger";

import { apiCallStarted, apiCallFailed, apiCallData } from "./actions";

export const initialState = {
    stateNames: [],
    cpsByState: {},
    apiOps: {},
};

function updateStateCps(existing = {}, action) {
    Log.info("Applying cpl updates from ", action);

    const wsName = action.req.getVal();
    

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
        out.stateNames = action.data.listList;
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
            apiOps: {
                ...existing.apiOps,
                [action.call]: {
                    call: action.call,
                    isLoading: true,
                    startedAt: Date.now(),
                    handler: action.handler,
                },
            },
        };
        break;

    case apiCallData.type:
        {
            Log.info("existing.apiOps ", existing.apiOps);
            const { startedAt } = existing.apiOps[action.call];
            const now = Date.now();

            out = {
                ...existing,
                apiOps: {
                    ...existing.apiOps,
                    [action.call]: {
                        call: action.call,
                        isLoading: false,
                        elapsedTime: now - startedAt,
                        updatedAt: now,
                    },
                },
            };
            // Dpesn't have to be pure...
            handleApiData(out, action);
        }
        break;

    case apiCallFailed.type:
        out = {
            ...existing,
            apiOps: {
                ...existing.apiOps,
                [action.call]: {
                    call: action.call,
                    isLoading: false,
                    lastError: action.error,
                },
            },
        };
        break;

    default:
        // Keep out as it is
        break;
    }

    return out;
}

// ////
