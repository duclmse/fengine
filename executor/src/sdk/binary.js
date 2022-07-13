export class Binary {
  private readonly schema: Schema;
  private readonly bytes: Uint8Array;
  private pointer: number = 0;

  constructor(schema: string | object, bytes: Iterable<number>) {
    this.schema = Binary.parseSchema(schema);
    Binary.validateSchema(this.schema);

    if (bytes instanceof Uint8Array) {
      this.bytes = bytes;
    } else if (Array.isArray(bytes)) {
      this.bytes = new Uint8Array(bytes);
    } else {
      throw new Error("bytes must be an int array or an instance of Uint8Array");
    }
  }

  static validateSchema(schema: Schema) {
    const {name, type, schemas} = schema;
    if (!name) throw new Error(`name is invalid or undefined`);
    if (!type) throw new Error(`type of ${name} is undefined`);
    switch (type) {
      case "[]object":
        schema.isArray = true;
        schema.subtype = "object";
      case "object":
        if (schemas == undefined) return false;
        for (let s of schemas) {
          if (!Binary.validateSchema(s)) return false;
        }
        return true;
      default:
        if (Binary.isValidType(type)) return true;
        return Binary.checkDefinedLengthArray(schema, type);
    }
  }

  static parseSchema(schema: string | object) {
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

  static checkDefinedLengthArray(schema: Schema, type: string) {
    let matched = /\[(\d*)](\w+)/.exec(type);
    if (!matched) return false;
    let [_, length, subtype] = matched;
    let b = Binary.isValidType(subtype);
    schema.length = length ? +length : 0;
    schema.isArray = true;
    schema.subtype = subtype;
    return b;
  }

  static isValidType(type: string) {
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

  static hex(a: number) {
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

  readObjectArray(schema: Schema, length: number) {
    let array = [];
    let {name, schemas} = schema;
    for (let i = 0; i < length; i++) {
      array.push(this.readObject({name, type: "object", schemas}));
    }
    return array;
  }

  readObject(schema: Schema) {
    const {type, isArray, length, subtype, schemas} = schema;
    if (isArray) {
      return this.readArray(schema, subtype, length);
    }
    if (type !== "object") {
      return this.readField(type, isArray, length, subtype);
    }
    if (schemas == undefined) {
      throw new Error("");
    }
    let result = {};
    for (let s of schemas) {
      let {name} = s;
      result[name] = this.readObject(s);
    }
    return result;
  }

  readArray<T>(schema: Schema, subtype: string | undefined, length: number): T[] {
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

  readGeneralArray<T>(length: number, reader: () => T): T[] {
    let array: T[] = [];
    for (let i = 0; i < length; i++) {
      array.push(reader());
    }
    return array;
  }

  readField(type: string | undefined, isArray: boolean, length: number, subtype: string | undefined) {
    switch (type) {
      case "byte":
      case "int8":
        return this.readByte(le);
      case "int16":
        return this.readInt16(le);
      case "int32":
        return this.readInt32(le);
      case "int64":
        return this.readInt64(le);
      case "float32":
        return this.readFloat32(le);
      case "float64":
        return this.readFloat64(le);
      case "string":
        return this.readString();
    }
  }

  readByte(le?: boolean): number {
    return this.bytes[this.pointer++] & 0xFF;
  }

  readInt16(le: boolean = false) {
    let b1 = this.bytes[this.pointer++] & 0xFF;
    let b2 = this.bytes[this.pointer++] & 0xFF;
    return (b1 << 8) | b2;
  }

  readInt32(le: boolean = false) {
    let i1 = this.readInt16(le);
    let i2 = this.readInt16(le);

    return (i1 << 16) | i2;
  }

  readInt64(le: boolean = false) {
    let i1 = this.readInt32(le);
    let i2 = this.readInt32(le);

    return (i1 * 4294967296) + i2;
  }

  readFloat32(le: boolean = false) {
    let array = this.readByteArray(4);
    return Buffer.from(array).readFloatBE(0);
  }

  readFloat64(le: boolean = false) {
    let array = this.readByteArray(8);
    return le ? Buffer.from(array).readDoubleLE(0) : Buffer.from(array).readDoubleBE(0);
  }

  readString(length: number) {
    let l = length || this.readByte();
    let array = this.readByteArray(l);
    return Buffer.from(array).toString("utf-8");
  }

  isDone() {
    console.log(`pointer - length = ${this.pointer - this.bytes.length}`);
    return this.pointer - this.bytes.length;
  }
}

export class Schema {
  private _name: string = "";
  private _type: string | undefined;
  private _subtype: string | undefined;
  private _schemas: Schema[] | undefined;
  private _isArray: boolean = false;
  private _length: number = 0;

  get name(): string {
    return this._name;
  }

  set name(value: string) {
    this._name = value;
  }

  get type() {
    return this._type;
  }

  set type(value) {
    this._type = value;
  }

  get subtype() {
    return this._subtype;
  }

  set subtype(value) {
    this._subtype = value;
  }

  get schemas(): Schema[] | undefined {
    return this._schemas;
  }

  set schemas(value: Schema[] | undefined) {
    this._schemas = value;
  }

  get isArray(): boolean {
    return this._isArray;
  }

  set isArray(value: boolean) {
    this._isArray = value;
  }

  get length(): number {
    return this._length;
  }

  set length(value: number) {
    this._length = value;
  }
}
