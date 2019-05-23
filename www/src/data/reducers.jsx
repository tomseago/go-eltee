import { combineReducers } from "redux";

import { apiCallStarted, apiCallFailed, apiCallData } from "./actions";

const initialWorldStates = {
    names: [],
    updatedAt: 0,
    isLoading: false,
};

function worldStates(existing = initialWorldStates, action) {
    let out = existing;

    if (action.call === "stateNames") {
        switch (action.type) {
        case apiCallStarted.type:
            out = {
                ...existing,
                isLoading: true,                
            };
            break;

        case apiCallData.type:
            out = {
                ...existing,
                isLoading: false,
                names: action.data.listList,
                lastError: false,
            };
            // out = Object.assign({}, ...existing, {
            //     names: action.data.getListList(),
            //     lastError: false,
            // });
            break;

        case apiCallFailed.type:
            out = {
                ...existing,
                isLoading: false,
                lastError: action.error,
            };
            // out = Object.assign({}, ...existing, {
            //     isLoading: false,
            //     lastError: action.error,
            // });
            break;

        default:
            // Keep out as it is
            break;
        }
    }

    return out;
}

const reducers = combineReducers({
    worldStates,
});

export default reducers;

export const initialState = {
    worldStates: initialWorldStates,
};
