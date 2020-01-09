import proto from "../api/api_pb";
import Log from "../lib/logger";

// Records a start for a specific call
export function apiCallStarted(call) {
    return {
        type: apiCallStarted.type,
        call,
        at: Date.now(),
    };
}
apiCallStarted.type = "API_CALL_STARTED";


export function apiCallFailed(call, req, error) {
    return {
        type: apiCallFailed.type,
        call,
        req,
        error,
        at: Date.now(),
    };
}
apiCallFailed.type = "API_CALL_FAILED";

export function apiCallData(call, req, data) {
    return {
        type: apiCallData.type,
        call,
        req,
        data,
        at: Date.now(),
    };
}
apiCallData.type = "API_CALL_DATA";


/**
 * This action is for making API calls that should only have one instance
 * outstanding at a time and which have a global handler.
 *
 * @param call the name of the api call to make
 * @param req the api request message
 */
export function maybeDoCall(call, req) {
    return {
        type: maybeDoCall.type,
        call,
        req,
    };
}
maybeDoCall.type = "MAYBE_DO_CALL";

export function maybeCallStateNames() {
    return maybeDoCall("stateNames", new proto.Void({}));
}

export function maybeCallControlPoints(wsName) {
    const req = new proto.StringMsg();
    req.setVal(wsName);
    return maybeDoCall(`controlPoints:${wsName}`, req);
}

/**
 * For api calls that don't need to check for existing in flight calls
 * use this action. Two handlers are provided for success and failure.
 * The success handler is given an already copied state object so all it
 * has to do is modify the parts of it that it wants to. It doesn't need
 * to be a pure function.
 * @param call the name of the rpc call
 * @param req the message to send as the request
 * @param success (stateCopy, resp, req, call) => {} success handler
 * @param failure (stateCopy, error, req, call) => {} failure handler
 * @returns {{call: *, success: *, failure: *, type: string, req: *}}
 */
export function doCall(call, req, success, failure) {
    return {
        type: doCall.type,
        call,
        req,
        success,
        failure,
    };
}
doCall.type = "DO_CALL";


export function handleCallResult(call, req, result, handler) {
    return {
        type: handleCallResult.type,
        call,
        req,
        result,
        handler,
    };
}
handleCallResult.type = "HANDLE_CALL_RESULT";
