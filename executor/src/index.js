const {credentials, ServerCredentials} = require("@grpc/grpc-js");
const {getServer} = require("./executor/server");
const {initClient} = require("./executor/grpc_client");

require("dotenv").config();

if (require.main === module) {
  const {ADDRESS, RESOLVER_ADDRESS} = process.env;
  const server = getServer();
  server.bindAsync(ADDRESS, ServerCredentials.createInsecure(), () => {
    console.log(`Server started at ${ADDRESS}`);
    server.start();
  });
  initClient(RESOLVER_ADDRESS, credentials.createInsecure());
}
