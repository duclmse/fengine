const grpc = require("@grpc/grpc-js");
const services = require("../pb/fengine_grpc_pb");
const executor = require("../executor/engine");

function getServer() {
  const server = new grpc.Server();
  server.addService(services.FEngineExecutorService, {execute});
  return server;
}

function execute(call, callback) {
  callback(null, executor.exec(call.request));
}

module.exports = {
  getServer,
};
