let bytes = new Uint8Array([
  0x11, 0x63, 0x00, 0x23, 0xCE, 0x11, 0x1D, 0xC1, 0x1A, 0xC1,
  0x0C, 0x10, 0x1C, 0x51, 0xCE, 0x11, 0x1D, 0xC2, 0x1A, 0xC2,
  0x0C, 0x20, 0x2C, 0x52, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15,
  0x0F, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x0F, 0xAA, 0xBB,
  0xCC, 0xDD, 0xEE, 0xFF, 0x0F
]);
let pointer = 0;
let connectionType = bytes[pointer++];

let metadata = {
  connection: (connectionType >> 4) & 0x0F,
  type: connectionType & 0x0F,
  battery: bytes[pointer++] & 0xFF,
  reserved: bytes[pointer++] & 0xFF
};
let hex = a => a.toString(16).toUpperCase();
let {cells, wifis, ptr} = readCellWifi(bytes, pointer);
console.log("----------------------------------------------");
console.log(metadata);
console.log(cells);
console.log(wifis);

console.log(ptr);

function readCellWifi(bytes /*: Uint8Array*/, pointer/*: number*/) {
  let qty = bytes[pointer++];
  let cellQty = (qty >> 4) & 0x0F;
  let wifiQty = qty & 0x0F;
  console.log(`cell=${cellQty}; wifi=${wifiQty}`);
  let {cells, ptr: ptr1} = readCells(bytes, cellQty, pointer);
  let {wifis, ptr: ptr2} = readWifis(bytes, wifiQty, ptr1);
  return {cells, wifis, ptr: ptr2};
}

function readCells(bytes /*: Uint8Array*/, qty/*: number*/, ptr/*: number*/) {
  let cells = [];
  for (let i = 0; i < qty; i++) {
    console.log(ptr);
    let c1 = bytes[ptr++] & 0xFF;
    let c2 = bytes[ptr++] & 0xFF;
    let c3 = bytes[ptr++] & 0xFF;
    let c4 = bytes[ptr++] & 0xFF;

    let cell = {
      cellId: ((c1 << 24) | (c2 << 16) | (c3 << 8) | (c4)) >>> 0, // 32bits
      lac: ((bytes[ptr++] << 8) | bytes[ptr++]) & 0xFFFF, // 16bits

      mcc: ((bytes[ptr++] << 4) | (bytes[ptr] >> 4)) & 0xFFFF, // 12bits
      mnc: ((bytes[ptr++] << 8) | bytes[ptr++]) & 0xFFFF,

      sig: (bytes[ptr++]) & 0xFFFF
    };
    console.log(`${hex(c1)} ${hex(c2)} ${hex(c3)} ${hex(c4)}`);
    console.log(hex(cell.cellId));
    cells.push(cell);
  }
  return {cells, ptr};
}

function readWifis(bytes /*: Uint8Array*/, qty/*: number*/, ptr/*: number*/) {
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

function MAC() {

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
