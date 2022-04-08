// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var pb_fengine_pb = require('../pb/fengine_pb.js');
var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');

function serialize_viot_DeleteRequest(arg) {
  if (!(arg instanceof pb_fengine_pb.DeleteRequest)) {
    throw new Error('Expected argument of type viot.DeleteRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_DeleteRequest(buffer_arg) {
  return pb_fengine_pb.DeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_InsertRequest(arg) {
  if (!(arg instanceof pb_fengine_pb.InsertRequest)) {
    throw new Error('Expected argument of type viot.InsertRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_InsertRequest(buffer_arg) {
  return pb_fengine_pb.InsertRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_Result(arg) {
  if (!(arg instanceof pb_fengine_pb.Result)) {
    throw new Error('Expected argument of type viot.Result');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_Result(buffer_arg) {
  return pb_fengine_pb.Result.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_Script(arg) {
  if (!(arg instanceof pb_fengine_pb.Script)) {
    throw new Error('Expected argument of type viot.Script');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_Script(buffer_arg) {
  return pb_fengine_pb.Script.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_ScriptRequest(arg) {
  if (!(arg instanceof pb_fengine_pb.ScriptRequest)) {
    throw new Error('Expected argument of type viot.ScriptRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_ScriptRequest(buffer_arg) {
  return pb_fengine_pb.ScriptRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_SelectRequest(arg) {
  if (!(arg instanceof pb_fengine_pb.SelectRequest)) {
    throw new Error('Expected argument of type viot.SelectRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_SelectRequest(buffer_arg) {
  return pb_fengine_pb.SelectRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_viot_UpdateRequest(arg) {
  if (!(arg instanceof pb_fengine_pb.UpdateRequest)) {
    throw new Error('Expected argument of type viot.UpdateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_viot_UpdateRequest(buffer_arg) {
  return pb_fengine_pb.UpdateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var FEngineExecutorService = exports.FEngineExecutorService = {
  execute: {
    path: '/viot.FEngineExecutor/Execute',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.Script,
    responseType: pb_fengine_pb.Result,
    requestSerialize: serialize_viot_Script,
    requestDeserialize: deserialize_viot_Script,
    responseSerialize: serialize_viot_Result,
    responseDeserialize: deserialize_viot_Result,
  },
};

exports.FEngineExecutorClient = grpc.makeGenericClientConstructor(FEngineExecutorService);
var FEngineThingService = exports.FEngineThingService = {
  resolveService: {
    path: '/viot.FEngineThing/ResolveService',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.ScriptRequest,
    responseType: pb_fengine_pb.Script,
    requestSerialize: serialize_viot_ScriptRequest,
    requestDeserialize: deserialize_viot_ScriptRequest,
    responseSerialize: serialize_viot_Script,
    responseDeserialize: deserialize_viot_Script,
  },
};

exports.FEngineThingClient = grpc.makeGenericClientConstructor(FEngineThingService);
var FEngineDataService = exports.FEngineDataService = {
  select: {
    path: '/viot.FEngineData/Select',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.SelectRequest,
    responseType: pb_fengine_pb.Script,
    requestSerialize: serialize_viot_SelectRequest,
    requestDeserialize: deserialize_viot_SelectRequest,
    responseSerialize: serialize_viot_Script,
    responseDeserialize: deserialize_viot_Script,
  },
  insert: {
    path: '/viot.FEngineData/Insert',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.InsertRequest,
    responseType: pb_fengine_pb.Script,
    requestSerialize: serialize_viot_InsertRequest,
    requestDeserialize: deserialize_viot_InsertRequest,
    responseSerialize: serialize_viot_Script,
    responseDeserialize: deserialize_viot_Script,
  },
  update: {
    path: '/viot.FEngineData/Update',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.UpdateRequest,
    responseType: pb_fengine_pb.Script,
    requestSerialize: serialize_viot_UpdateRequest,
    requestDeserialize: deserialize_viot_UpdateRequest,
    responseSerialize: serialize_viot_Script,
    responseDeserialize: deserialize_viot_Script,
  },
  delete: {
    path: '/viot.FEngineData/Delete',
    requestStream: false,
    responseStream: false,
    requestType: pb_fengine_pb.DeleteRequest,
    responseType: pb_fengine_pb.Script,
    requestSerialize: serialize_viot_DeleteRequest,
    requestDeserialize: deserialize_viot_DeleteRequest,
    responseSerialize: serialize_viot_Script,
    responseDeserialize: deserialize_viot_Script,
  },
};

exports.FEngineDataClient = grpc.makeGenericClientConstructor(FEngineDataService);
