require("dotenv").config();

const grpc = require("@grpc/grpc-js");
const {getServer} = require("./grpc/server");

if (require.main === module) {
  const {ADDRESS} = process.env
  const server = getServer();
  server.bindAsync(ADDRESS, grpc.ServerCredentials.createInsecure(), () => {
    console.log(`Server started at: ${ADDRESS}`);
    server.start();
  });
}
