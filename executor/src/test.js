const typedArray = new Uint8Array([1, 2, 3, 4]);


const proto = require("protobufjs");

const any = new Any();
const binarySerialized = [];
any.pack(binarySerialized, "foo.Bar");
console.log(any.getTypeName());
