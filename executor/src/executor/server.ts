import {sendUnaryData, Server, ServerUnaryCall} from "@grpc/grpc-js";
import {FEngineExecutorService} from "../pb/fengine_grpc_pb";
import {Result, Script} from "../pb/fengine_pb";
import {Executor} from "./engine";

export function getServer() {
  const server = new Server();
  server.addService(FEngineExecutorService, {execute});
  return server;
}

function execute(call: ServerUnaryCall<Script, Result>, callback: sendUnaryData<Result>) {
  callback(null, new Executor().exec(call.request));
}
