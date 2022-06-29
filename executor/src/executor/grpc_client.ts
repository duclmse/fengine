import {FEngineDataClient} from "../pb/fengine_grpc_pb";
import {ChannelCredentials, Client} from "@grpc/grpc-js";

let client: Client;

export function initClient(address: string, credential: ChannelCredentials) {
  if (client == null) {
    client = new FEngineDataClient(address, credential);
  }
}

export function getClient(): Client {
  return client;
}
