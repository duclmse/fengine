const {VM} = require("vm2");

// const {fengine} = require("../grpc/server")

const Type = {
  Any: 0,
  Int32: 1,
  Int64: 2,
  Float: 3,
  Double: 4,
  Bool: 5,
  String: 6,
  Bytes: 7,
};

const TypeName = {
  0: "Int32",
  1: "Int64",
  2: "Float",
  3: "Double",
  4: "Bool",
  5: "String",
  6: "Bytes",
};

function exec(script) {
  const {"function": func, referee} = script;
  const vm = new VM({
    require: {external: true},
    sandbox: buildSandbox(referee),
  });

  try {
    // console.log(JSON.stringify(script));
    const input = parseInput(func.input);
    const code = `((${input.args}) => {${script.function.code}})(${input.values.join(",")})`;
    console.log(`---\n${code}\n---`);

    return wrap(vm.run(code));
  } catch (e) {
    console.log(e);
    return JSON.stringify({name: "error", type: Type.String, str: e.message});
  }
}

function buildSandbox(referee) {
  const me = {};
  for (let fn in referee) {
    me[fn] = function (input) {
      console.log(input);
    };
  }
  return {me};
}

function parseInput(input) {
  const args = [];
  const values = [];
  for (let inp of input) {
    args.push(inp.name);
    let value = readValue(inp);
    values.push(value);
    console.log(`${inp.name}: ${value}`);
  }
  return {args, values};
}

function readValue(input) {
  switch (input.value) {
    case "i32":
      return input["i32"];
    case "i64":
      return input["i64"];
    case "f32":
      return input["f32"];
    case "f64":
      return input["f64"];
    case "bol":
      return input["bol"];
    case "str":
      return `'${input["str"]}'`;
    case "bin":
      return input["bin"];
  }
}

function wrap(variable) {
  console.log(`-> ${variable}`);
  switch (typeof variable) {
    case "object":
      return {name: "result", type: Type.Any, any: variable};
    case "boolean":
      return {name: "result", type: Type.Bool, bol: variable};
    case "number":
      if (Number.isInteger(variable)) {
        return {name: "result", type: Type.Int64, i64: variable};
      }
      return {name: "result", type: Type.Double, f64: variable};
    case "string":
      return {name: "result", type: Type.String, str: variable};
  }
}

module.exports = {exec};
