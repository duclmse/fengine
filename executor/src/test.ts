import {Table} from "./sdk";
import {initClient} from "./executor/grpc_client";
import {credentials} from "@grpc/grpc-js";

require("dotenv").config();
let {RESOLVER_ADDRESS} = process.env;
initClient(RESOLVER_ADDRESS, credentials.createInsecure());
(async ({s, i}) => {
  let res = await Table("tbl_test").Select({
    filter: {
      $and: [
        {a: {$gt: 10, $lt: 20}}
      ]
    }
  });
  let all = [];
  for (let row of res) {
    let vals = [];
    for (let v of row) {
      vals.push(v);
    }
    all.push(vals);
  }
  console.log(res.GetColumns());
  console.log(all);
})({s: "input", i: 32});
