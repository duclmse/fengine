const pb = require("../pb/fengine_pb");
const {library} = require("../sdk");
const {VM} = require("vm2");

function exec(script) {
  try {
    const fn = script.getFunction();
    let {sandbox, code: sandboxCode} = buildSandbox(script);
    const vm = new VM({require: {external: true}, sandbox});
    const {args, params} = parseInput(fn.getInputList());
    const code = `((${params})=>{try{${fn.getCode()}}catch(_e_){return _e_}})(${args.join()})`;
    console.log(JSON.stringify(sandbox));
    console.log(`>---\n${sandboxCode}${code}\n---<`);

    let res = wrap(new pb.Result(), vm.run(sandboxCode + code));
    return res
  } catch (e) {
    console.log(e);
    return JSON.stringify({name: "error", type: pb.Type.String, str: e.message});
  }
}

function buildSandbox(script) {
  const me = {};
  let code = "";
  script.getAttributesList().forEach(attr => {
    me[attr.getName()] = readVarValue(attr);
  });

  script.getRefereeMap().forEach((fn, name) => {
    console.log(`referee ${name}: ${fn.getInputList()}`);
    let _code = fn.getCode();
    if (_code) {
      const {args} = parseInput(fn.getInputList(), false);
      code += `me['${name}'] = (${args}) => {${_code}};\n`;
    } else {
      me[name] = input => {
        console.log(input);
      };
    }
  });

  return {sandbox: {me, ...library}, code};
}

function parseInput(input, hasArgs = true) {
  const args = [];
  const params = [];
  input.forEach(inp => {
    params.push(inp.getName());
    if (hasArgs) args.push(readVarValue(inp, true));
    else console.log(inp.getName());
  });
  if (!hasArgs) console.log(JSON.stringify(args));
  return {args, params};
}

function readVarValue(input, isParam) {
  switch (true) {
    case input.hasI32():
      return input.getI32();
    case input.hasI64():
      return input.getI64();
    case input.hasF32():
      return input.getF32();
    case input.hasF64():
      return input.getF64();
    case input.hasBol():
      return input.getBol();
    case input.hasStr():
      return isParam ? `'${input.getStr()}'` : input.getStr();
    case input.hasBin():
      return input.getBin_asU8();
  }
}

function wrap(result) {
  const variable = new pb.Variable();
  switch (typeof result) {
    case "object":
      variable.setType(pb.Type.OBJECT);
      if (result instanceof Error) {
        console.log(`error: ${result.message}`);
        variable.setStr(JSON.stringify({error: result.message}));
      } else {
        variable.setStr(JSON.stringify(result));
      }
      break;
    case "boolean":
      variable.setType(pb.Type.BOOL);
      variable.setBol(result);
      break;
    case "number":
      if (Number.isInteger(result)) {
        if (result < 4294967296) {
          variable.setType(pb.Type.INT32);
          variable.setI32(result);

        } else {
          variable.setType(pb.Type.INT64);
          variable.setI64(result);
        }
        break;
      }
      variable.setType(pb.Type.DOUBLE);
      variable.setF64(result);
      break;
    case "string":
      variable.setType(pb.Type.STRING);
      variable.setStr(result);
      break;
  }
  return variable;
}

module.exports = {exec};
