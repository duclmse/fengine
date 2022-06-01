let Binary = require("./binary");

let schemaJson = `{
  "name": "payload",
  "type": "object",
  "schemas": [
    {"name": "connection", "type": "byte", "le": true},
    {"name": "type", "type": "byte"},
    {"name": "battery", "type": "byte"},
    {
      "name": "cell",
      "type": "[]object",
      "schemas": [
        {"name": "cellId", "type": "[4]byte"},
        {"name": "lac", "type": "int16"},
        {"name": "mcc", "type": "int16"},
        {"name": "mnc", "type": "int16"},
        {"name": "sig", "type": "byte"}
      ]
    },
    {
      "name": "wifi",
      "type": "[]object",
      "schemas": [
        {"name": "mac", "type": "[6]byte"},
        {"name": "sig", "type": "byte"}
      ]
    },
    {"name": "test", "type": "[]int16"}
  ]
}`;

let bytes = new Uint8Array([
  0x01, 0x01, 0x63, 0x02, 0xce, 0x11, 0x1d, 0xc1, 0x1a, 0x00, 0xc1, 0x0c, 0x10, 0x1c, 0x51, 0xce, 0x11, 0x1d, 0xc2,
  0x1a, 0x00, 0xc2, 0x0c, 0x20, 0x2c, 0x52, 0x03, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x0f, 0xa1, 0xb2, 0xc3, 0xd4,
  0xe5, 0xf6, 0x0f, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x0f, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
]);
console.log(`length = ${bytes.length}`);
console.time("all");
let payload = new Binary(schemaJson, bytes);
let object = payload.toObject();
console.timeLog("all");

console.time("log");
console.log(`----------------------`);
console.log(JSON.stringify(object));
console.log(`----------------------`);
console.log(payload.isDone());
console.timeEnd("log");
