let schemaJson = `{
  "name": "payload",
  "type": "object",
  "schemas": [
    {
      "name": "metadata",
      "type": "object",
      "schemas": [
        {"name": "connection", "type": "byte"},
        {"name": "type", "type": "byte"},
        {"name": "battery", "type": "byte"}
      ]
    },
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
      "name": "payload",
      "type": "[]object",
      "schemas": [
        {"name": "mac", "type": "[6]byte"},
        {"name": "sig", "type": "byte"}]
    }
  ]
}`;

console.log(`validateSchema=${validateSchema(schema)}`);

let bytes = new Uint8Array([
  0x01, 0x01, 0x63, 0x00, 0x02,
  0xCE, 0x11, 0x1D, 0xC1, 0x1A, 0x00, 0xC1, 0x0C, 0x10, 0x1C, 0x51,
  0xCE, 0x11, 0x1D, 0xC2, 0x1A, 0x00, 0xC2, 0x0C, 0x20, 0x2C, 0x52,
  0x03, 0x10, 0x11, 0x12, 0x13,
  0x14, 0x15, 0x0F, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x0F,
  0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x0F
]);

let hex = a => {
  let h = a.toString(16).toUpperCase();
  return h.length % 2 === 0 ? h : `0${h}`;
};

class Binary {
  constructor(schema, bytes) {
    this.schema = Binary.validateSchema(schema);
    if (bytes instanceof Uint8Array) {
      this.bytes = bytes;
    } else if (Array.isArray(bytes)) {
      this.bytes = new Uint8Array(bytes);
    } else {
      throw new Error("bytes must be an int array or an instance of Uint8Array");
    }
    this.pointer = 0;
  }

  static validateSchema(schema) {
    switch (typeof schema) {
      case "string":
        JSON.parse(schemaJson, (k, v) => {
          console.log(`${k} ${v}`);
          if (k === "name" && typeof v != "string") throw "name must be a string";
          if (k === "type" && typeof v != "string") throw "type must be a string";
          return v;
        });
        break;
      case "object":
        break;
      default:
        throw new Error("schema must be a string or an object");
    }
    const {name, type, schemas} = schema;

    switch (type) {
      case "object":
      case "[]object":
        if (!Array.isArray(schemas)) return false;
        for (let s of schemas) {
          if (!validateSchema(s)) return false;
        }
        return true;
      default:
        if (Binary.isValidType(type)) return true;
        let arrayDef = Binary.checkDefinedLengthArray(type);
        if (!arrayDef) {
          return false;
        }
        return Binary.isValidType(arrayDef.type);
    }
  }

  static checkDefinedLengthArray(type) {
    let matched = /\[(\d+)](\w+)/.exec(type);
    if (!matched) return null;
    return {length: matched[1], type: matched[2]};
  }

  static isValidType(type) {
    switch (type) {
      case "byte":
      case "int8":
      case "int16":
      case "int32":
      case "int64":
      case "float32":
      case "float64":
        return true;
      default:
        return false;
    }
  }

  toObject() {

  }

  readByteArray(length = 1) {
    let arr = [];
    for (let i = 0; i < length; i++) {
      arr.push(this.bytes[this.pointer++] & 0xFF);
    }
    return new Uint8Array(arr);
  }

  readBytes(length = 1) {
    let value = 0;
    for (let i = 0; i < length; i++) {
      value = (value << 8 | this.bytes[this.pointer++] & 0xFF) >>> 0;
    }
    return value;
  }

  readObjectArray(schema) {

  }

  readObject(schema) {

  }

  readField() {

  }

  readByte() {
    return (this.bytes[this.pointer++] & 0xFF) >>> 0;
  }

  readInt8() {
    return this.bytes[this.pointer++] & 0xFF;
  }

  readInt16() {
    return this.readBytes(2);
  }

  readInt32() {
    return this.readBytes(4);
  }

  readInt64() {
    return this.readBytes(8);
  }

  readFloat32() {
    let array = this.readByteArray(4);
    return Buffer.from(array).readFloatBE(0);
  }

  readFloat64() {
    let array = this.readByteArray(8);
    return Buffer.from(array).readDoubleBE(0);
  }
}

function readMetadata(bytes, pointer) {
  let {value: connection, pointer: p1} = readBytes(bytes, pointer, 1);
  let {value: type, pointer: p2} = readBytes(bytes, p1, 1);
  let {value: battery, pointer: p3} = readBytes(bytes, p2, 1);
  let {value: reserved, pointer: p4} = readBytes(bytes, p3, 1);
  let metadata = {connection, type, battery, reserved};
  return {metadata, pointer: p4};
}

function readCellWifi(bytes /*: UInt8Array*/, pointer/*: number*/) {
  let cellQty = bytes[pointer++] & 0xFF;
  let {cells, ptr: ptr1} = readCells(bytes, cellQty, pointer);
  let wifiQty = bytes[ptr1++] & 0xFF;
  let {wifis, ptr: ptr2} = readWifis(bytes, wifiQty, ptr1);

  return {cells, wifis, ptr: ptr2};
}

function readCells(bytes /*: UInt8Array*/, qty/*: number*/, ptr/*: number*/) {
  let cells = [];
  for (let i = 0; i < qty; i++) {
    let {value: cid, pointer: p1} = readBytes(bytes, ptr, 4);
    let {value: lac, pointer: p2} = readBytes(bytes, p1, 2);
    let {value: mcc, pointer: p3} = readBytes(bytes, p2, 2);
    let {value: mnc, pointer: p4} = readBytes(bytes, p3, 2);
    let {value: sig, pointer: p5} = readBytes(bytes, p4, 1);
    cells.push({cellId: hex(cid), lac, mcc, mnc, sig});
    ptr = p5;
  }
  return {cells, ptr};
}

function readWifis(bytes /*: UInt8Array*/, qty/*: number*/, ptr/*: number*/) {
  let wifis = [];
  for (let i = 0; i < qty; i++) {
    let mac = hex(bytes[ptr++]);
    for (let i = 1; i < 6; i++) {
      mac += `:${hex(bytes[ptr++])}`;
    }
    let sig = bytes[ptr++];
    wifis.push({mac, sig});
  }
  return {wifis, ptr};
}

// function Thing(name) {
//   try {
//     const fn = script.getFunction();
//     if (!fn) {
//       let json = new Variable().setJson(JSON.stringify({error: "Function is not defined"}));
//       return new Result().setOutput(json);
//     }
//
//     const {sandbox, code: sandboxCode, attributes} = E.buildSandbox(script);
//     const {args, params} = E.parseArguments(fn.getInputList());
//     const code = `((${params})=>{try{${fn.getCode()}}catch(_e_){return _e_}})(${args.join()})`;
//     console.debug(`${JSON.stringify(sandbox)}>---\n${sandboxCode}\n${code}\n---<`);
//
//     const vm = new VM({sandbox});
//     const label = new Date().getTime();
//     console.time(`${label}`);
//     let output = E.wrap(vm.run(sandboxCode + code), fn.getOutput());
//     console.timeEnd(`${label}`);
//     E.compareAttributes(sandbox.me, attributes);
//
//     return new Result().setOutput(output);
//   } catch (e) {
//     return new Result().setOutput(new Variable().setString(e.message));
//   }
// }
