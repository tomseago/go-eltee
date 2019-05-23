import proto from "../api/api_pb";

export function apiCallStart(call, req) {
    return {
        type: apiCallStart.type,
        call,
        req,
    };
}
apiCallStart.type = "API_CALL_START";


export function apiCallStarted(call, req) {
    return {
        type: apiCallStarted.type,
        call,
        req,
    };
}
apiCallStarted.type = "API_CALL_STARTED";


export function apiCallFailed(call, req, error) {
    return {
        type: apiCallFailed.type,
        call,
        req,
        error,
    };
}
apiCallFailed.type = "API_CALL_FAILED";

export function apiCallData(call, req, data) {
    return {
        type: apiCallData.type,
        call,
        data,
    };
}
apiCallData.type = "API_CALL_DATA";


export function callStateNames() {
    // return apiCallStart("stateNames", new proto.Void({}));
    const p = new proto.Void({});
    const x = apiCallStart("stateNames", p);
    return x;
}
callStateNames.type = "CALL_STATE_NAMES";
