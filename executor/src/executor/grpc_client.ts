import {FEngineDataClient} from "../pb/fengine_grpc_pb";
import {ChannelCredentials} from "@grpc/grpc-js";

let client: FEngineDataClient;

export function initClient(address: string, credential: ChannelCredentials) {
  if (client == null) {
    console.log(`Initialized gRPC client at ${address}`);
    client = new FEngineDataClient(address, credential);
  }
}

export function getClient(): FEngineDataClient {
  return client;
}
