module.exports =  class Binary {
  constructor(schema, bytes) {
    this.schema = Binary.parseSchema(schema);
    Binary.validateSchema(this.schema);

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
    const {name, type, schemas} = schema;
    if (!name) throw new Error(`name is invalid or undefined`);
    if (!type) throw new Error(`type of ${name} is undefined`);
    switch (type) {
      case "[]object":
        schema.isArray = true;
        schema.subtype = "object";
      case "object":
        for (let s of schemas) {
          if (!Binary.validateSchema(s)) return false;
        }
        return true;
      default:
        if (Binary.isValidType(type)) return true;
        return Binary.checkDefinedLengthArray(schema, type);
    }
  }

  static parseSchema(schema) {
    switch (typeof schema) {
      case "string":
        return JSON.parse(schema, (k, v) => {
          if (k === "name" && typeof v != "string") throw new Error("name must be a string");
          if (k === "type" && typeof v != "string") throw new Error("type must be a string");
          if (k === "schemas" && !Array.isArray(v)) throw new Error("schemas must be an array");
          return v;
        });
      case "object":
        return schema;
      default:
        throw new Error("schema must be a string or an object");
    }
  }

  static checkDefinedLengthArray(schema, type) {
    let matched = /\[(\d*)](\w+)/.exec(type);
    if (!matched) return false;
    let [_, length, subtype] = matched;
    let b = Binary.isValidType(subtype);
    schema.length = length ? +length : 0;
    schema.isArray = true;
    schema.subtype = subtype;
    return b;
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
      case "string":
      case "object":
        return true;
      default:
        return false;
    }
  }

  static hex(a) {
    let h = a.toString(16).toUpperCase();
    return h.length % 2 === 0 ? h : `0${h}`;
  };

  toObject() {
    return this.readObject(this.schema);
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
      value = ((value * 256) + (this.bytes[this.pointer++] & 0xFF));
    }
    return value;
  }

  readObjectArray(schema, length) {
    let array = [];
    let {name, schemas} = schema;
    for (let i = 0; i < length; i++) {
      array.push(this.readObject({name, type: "object", schemas}));
    }
    return array;
  }

  readObject(schema) {
    const {type, isArray, length, subtype, schemas} = schema;
    if (isArray) {
      return this.readArray(schema, subtype, length);
    }
    if (type === "object") {
      let result = {};
      for (let s of schemas) {
        let {name} = s;
        result[name] = this.readObject(s);
      }
      return result;
    }
    return this.readField(type, isArray, length, subtype);
  }

  readArray(schema, subtype, length) {
    length = length || this.readBytes();
    switch (subtype) {
      case "byte":
      case "int8":
        return this.readGeneralArray(length, () => this.readByte());
      case "int16":
        return this.readGeneralArray(length, () => this.readInt16());
      case "int32":
        return this.readGeneralArray(length, () => this.readInt32());
      case "int64":
        return this.readGeneralArray(length, () => this.readInt64());
      case "float32":
        return this.readGeneralArray(length, () => this.readFloat32());
      case "float64":
        return this.readGeneralArray(length, () => this.readFloat64());
      case "string":
        return this.readGeneralArray(length, () => this.readString());
      case "object":
        return this.readObjectArray(schema, length);
    }
  }

  readGeneralArray(length, reader) {
    let array = [];
    for (let i = 0; i < length; i++) {
      array.push(reader());
    }
    return array;
  }

  readField(type) {
    switch (type) {
      case "byte":
      case "int8":
        return this.readByte();
      case "int16":
        return this.readInt16();
      case "int32":
        return this.readInt32();
      case "int64":
        return this.readInt64();
      case "float32":
        return this.readFloat32();
      case "float64":
        return this.readFloat64();
      case "string":
        return this.readString();
    }
  }

  readByte() {
    return this.bytes[this.pointer++] & 0xFF;
  }

  readInt16() {
    let b1 = this.bytes[this.pointer++] & 0xFF;
    let b2 = this.bytes[this.pointer++] & 0xFF;
    return (b1 << 8) | b2;
  }

  readInt32() {
    let i1 = this.readInt16();
    let i2 = this.readInt16();

    return (i1 << 16) | i2;
  }

  readInt64() {
    let i1 = this.readInt32();
    let i2 = this.readInt32();

    return (i1 * 4294967296) + i2;
  }

  readFloat32() {
    let array = this.readByteArray(4);
    return Buffer.from(array).readFloatBE(0);
  }

  readFloat64() {
    let array = this.readByteArray(8);
    return Buffer.from(array).readDoubleBE(0);
  }

  readString(length) {
    let l = length || this.readByte();
    let array = this.readByteArray(l);
    return Buffer.from(array).toString("utf-8");
  }

  isDone() {
    return this.pointer - this.bytes.length;
  }
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
