import {sendUnaryData, Server, ServerUnaryCall} from "@grpc/grpc-js";
import {FEngineExecutorService} from "../pb/fengine_grpc_pb";
import {Result, Script} from "../pb/fengine_pb";
import {Executor} from "./engine";
import {Cache} from "./cache";

const cache = new Cache();
const executor = new Executor(cache);

const server = new Server();
server.addService(FEngineExecutorService, {
  execute,
  upsertService,
  deleteService
});

export function getServer() {
  return server;
}


function execute(call: ServerUnaryCall<Script, Result>, callback: sendUnaryData<Result>) {
  callback(null, executor.exec(call.request));
}

function upsertService(call: ServerUnaryCall<Script, Result>, callback: sendUnaryData<Result>) {
  callback(null, executor.upsertService(call.request));
}

function deleteService(call: ServerUnaryCall<Script, Result>, callback: sendUnaryData<Result>) {
  callback(null, executor.deleteService(call.request));
}
