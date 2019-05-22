/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = require('./api_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ElTeeClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ElTeePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.StringMsg>}
 */
const methodInfo_ElTee_Ping = new grpc.web.AbstractClientBase.MethodInfo(
  proto.StringMsg,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.StringMsg.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.StringMsg)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.StringMsg>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.ping =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/Ping',
      request,
      metadata || {},
      methodInfo_ElTee_Ping,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.StringMsg>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.ping =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/Ping',
      request,
      metadata || {},
      methodInfo_ElTee_Ping);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.Void,
 *   !proto.ProfilesResponse>}
 */
const methodInfo_ElTee_ProfileLibrary = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ProfilesResponse,
  /** @param {!proto.Void} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ProfilesResponse.deserializeBinary
);


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ProfilesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ProfilesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.profileLibrary =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/ProfileLibrary',
      request,
      metadata || {},
      methodInfo_ElTee_ProfileLibrary,
      callback);
};


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ProfilesResponse>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.profileLibrary =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/ProfileLibrary',
      request,
      metadata || {},
      methodInfo_ElTee_ProfileLibrary);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.Void,
 *   !proto.StringMsg>}
 */
const methodInfo_ElTee_StateNames = new grpc.web.AbstractClientBase.MethodInfo(
  proto.StringMsg,
  /** @param {!proto.Void} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.StringMsg.deserializeBinary
);


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.StringMsg)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.StringMsg>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.stateNames =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/StateNames',
      request,
      metadata || {},
      methodInfo_ElTee_StateNames,
      callback);
};


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.StringMsg>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.stateNames =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/StateNames',
      request,
      metadata || {},
      methodInfo_ElTee_StateNames);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.ControlPointList>}
 */
const methodInfo_ElTee_ControlPoints = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ControlPointList,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.ControlPointList.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ControlPointList)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ControlPointList>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.controlPoints =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/ControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_ControlPoints,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ControlPointList>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.controlPoints =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/ControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_ControlPoints);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.FixturePatchMap>}
 */
const methodInfo_ElTee_FixturePatches = new grpc.web.AbstractClientBase.MethodInfo(
  proto.FixturePatchMap,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.FixturePatchMap.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.FixturePatchMap)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.FixturePatchMap>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.fixturePatches =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/FixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_FixturePatches,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.FixturePatchMap>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.fixturePatches =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/FixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_FixturePatches);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ControlPointList,
 *   !proto.Void>}
 */
const methodInfo_ElTee_SetControlPoints = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.ControlPointList} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.ControlPointList} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.setControlPoints =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/SetControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_SetControlPoints,
      callback);
};


/**
 * @param {!proto.ControlPointList} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.setControlPoints =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/SetControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_SetControlPoints);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ControlPointList,
 *   !proto.Void>}
 */
const methodInfo_ElTee_RemoveControlPoints = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.ControlPointList} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.ControlPointList} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.removeControlPoints =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/RemoveControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveControlPoints,
      callback);
};


/**
 * @param {!proto.ControlPointList} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.removeControlPoints =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/RemoveControlPoints',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveControlPoints);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.FixturePatchMap,
 *   !proto.Void>}
 */
const methodInfo_ElTee_SetFixturePatches = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.FixturePatchMap} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.FixturePatchMap} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.setFixturePatches =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/SetFixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_SetFixturePatches,
      callback);
};


/**
 * @param {!proto.FixturePatchMap} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.setFixturePatches =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/SetFixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_SetFixturePatches);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.FixturePatchMap,
 *   !proto.Void>}
 */
const methodInfo_ElTee_RemoveFixturePatches = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.FixturePatchMap} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.FixturePatchMap} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.removeFixturePatches =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/RemoveFixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveFixturePatches,
      callback);
};


/**
 * @param {!proto.FixturePatchMap} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.removeFixturePatches =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/RemoveFixturePatches',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveFixturePatches);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.Void>}
 */
const methodInfo_ElTee_ApplyState = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.applyState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/ApplyState',
      request,
      metadata || {},
      methodInfo_ElTee_ApplyState,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.applyState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/ApplyState',
      request,
      metadata || {},
      methodInfo_ElTee_ApplyState);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.Void,
 *   !proto.StringMsg>}
 */
const methodInfo_ElTee_LoadableStateNames = new grpc.web.AbstractClientBase.MethodInfo(
  proto.StringMsg,
  /** @param {!proto.Void} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.StringMsg.deserializeBinary
);


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.StringMsg)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.StringMsg>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.loadableStateNames =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/LoadableStateNames',
      request,
      metadata || {},
      methodInfo_ElTee_LoadableStateNames,
      callback);
};


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.StringMsg>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.loadableStateNames =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/LoadableStateNames',
      request,
      metadata || {},
      methodInfo_ElTee_LoadableStateNames);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.Void>}
 */
const methodInfo_ElTee_LoadLoadableState = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.loadLoadableState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/LoadLoadableState',
      request,
      metadata || {},
      methodInfo_ElTee_LoadLoadableState,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.loadLoadableState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/LoadLoadableState',
      request,
      metadata || {},
      methodInfo_ElTee_LoadLoadableState);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.Void>}
 */
const methodInfo_ElTee_SaveLoadableState = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.saveLoadableState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/SaveLoadableState',
      request,
      metadata || {},
      methodInfo_ElTee_SaveLoadableState,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.saveLoadableState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/SaveLoadableState',
      request,
      metadata || {},
      methodInfo_ElTee_SaveLoadableState);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.Void,
 *   !proto.Void>}
 */
const methodInfo_ElTee_SaveAllStates = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.Void} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.saveAllStates =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/SaveAllStates',
      request,
      metadata || {},
      methodInfo_ElTee_SaveAllStates,
      callback);
};


/**
 * @param {!proto.Void} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.saveAllStates =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/SaveAllStates',
      request,
      metadata || {},
      methodInfo_ElTee_SaveAllStates);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.Void>}
 */
const methodInfo_ElTee_AddState = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.addState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/AddState',
      request,
      metadata || {},
      methodInfo_ElTee_AddState,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.addState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/AddState',
      request,
      metadata || {},
      methodInfo_ElTee_AddState);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.StringMsg,
 *   !proto.Void>}
 */
const methodInfo_ElTee_RemoveState = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.StringMsg} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.removeState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/RemoveState',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveState,
      callback);
};


/**
 * @param {!proto.StringMsg} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.removeState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/RemoveState',
      request,
      metadata || {},
      methodInfo_ElTee_RemoveState);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.SrcDest,
 *   !proto.Void>}
 */
const methodInfo_ElTee_CopyStateTo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.SrcDest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.copyStateTo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/CopyStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_CopyStateTo,
      callback);
};


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.copyStateTo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/CopyStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_CopyStateTo);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.SrcDest,
 *   !proto.Void>}
 */
const methodInfo_ElTee_MoveStateTo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.SrcDest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.moveStateTo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/MoveStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_MoveStateTo,
      callback);
};


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.moveStateTo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/MoveStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_MoveStateTo);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.SrcDest,
 *   !proto.Void>}
 */
const methodInfo_ElTee_ApplyStateTo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.Void,
  /** @param {!proto.SrcDest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.Void.deserializeBinary
);


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.Void)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.Void>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ElTeeClient.prototype.applyStateTo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ElTee/ApplyStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_ApplyStateTo,
      callback);
};


/**
 * @param {!proto.SrcDest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.Void>}
 *     A native promise that resolves to the response
 */
proto.ElTeePromiseClient.prototype.applyStateTo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ElTee/ApplyStateTo',
      request,
      metadata || {},
      methodInfo_ElTee_ApplyStateTo);
};


module.exports = proto;

