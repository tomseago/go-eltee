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

export function callControlPoints(wsName) {
    return apiCallStart("controlPoints", new proto.StringMsg({ val: wsName }));
}
callControlPoints.type = "CALL_CONTROL_POINTS";

// /////////

// export function ensureRecentStateNames() {
//     return {
//         type: ensureRecentStateNames.type,
//     };
// }
// ensureRecentStateNames.type = "ENSURE_STATE_NAMES";


export function maybeDoOp(opVal, req) {
    return {
        type: maybeDoOp.type,
        op: opVal,
        req,
    };
}
maybeDoOp.type = "MAYBE_DO_OP";

export function maybeCallStateNames(opVal) {
    return maybeDoOp(opVal, new proto.Void({}));
}

export function maybeCallControlPoints(opVal, wsName) {
    return maybeDoOp(opVal, new proto.StringMsg({ val: wsName }));
}
