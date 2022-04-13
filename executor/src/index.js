const {ServerCredentials} = require("@grpc/grpc-js");
const {getServer} = require("./grpc/server");

require("dotenv").config();

if (require.main === module) {
  const {ADDRESS} = process.env
  const server = getServer();
  server.bindAsync(ADDRESS, ServerCredentials.createInsecure(), () => {
    console.log(`Server started at: ${ADDRESS}`);
    server.start();
  });
}
