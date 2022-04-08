const path = require("path");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const {Server: Server1} = require("@grpc/grpc-js");

const executor = require("../executor/engine1")

function execute(call, callback) {
  callback(null, executor.exec(call.request));
}

/**
 * Get a new server with the handler functions in this file bound to the methods it serves.
 * @return {Server1} The new server object
 */
function getServer(protoPath) {
  const packageDefinition = protoLoader.loadSync(path.resolve(protoPath), {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  });
  const fengine = grpc.loadPackageDefinition(packageDefinition).viot;
  const server = new grpc.Server();
  server.addService(fengine.FEngineExecutor.service, {
    execute: execute
  });
  return server;
}

module.exports = {
  getServer
};
