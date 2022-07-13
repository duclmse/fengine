import {credentials, ServerCredentials} from "@grpc/grpc-js";
import {getServer} from "./executor/server";
import {initClient} from "./executor/grpc_client";

require("dotenv").config();

if (require.main === module) {
  const {ADDRESS, RESOLVER_ADDRESS} = process.env;
  const server = getServer();
  server.bindAsync(ADDRESS || "localhost:1234", ServerCredentials.createInsecure(), () => {
    console.log(`Server started at ${ADDRESS}`);
    server.start();
  });
  initClient(RESOLVER_ADDRESS || "localhost:1235", credentials.createInsecure());
}
