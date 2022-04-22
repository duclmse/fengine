import {Function, Parameter, Result, Script, Type, Variable} from "../pb/fengine_pb";
import {Cache} from "./cache";
import * as library from "../sdk/db";
import {VM} from "vm2";
import _ from "lodash";

type MsgType = void | number | string | boolean | Uint8Array;
type Func = (input: any) => MsgType;
type ThingReference = {
  [key: string]: MsgType | Func
}

class E {
  cache: Cache;

  constructor(cache: Cache) {
    this.cache = cache;
  }

  private static buildSandbox(script: Script) {
    let code = "";
    const me: ThingReference = {};
    const attributes: ThingReference = {};

    // script.getAttributesList().forEach(attr => {
    //   let name = attr.getName();
    //   let value: MsgType = E.readVarValue(attr);
    //   me[name] = value;
    //   attributes[name] = value;
    // });

    script.getServicesMap().forEach((fn: Function, name: string) => {
      let _code = fn.getCode();
      if (_code) {
        const params = E.parseParameters(fn.getInputList());
        code += `me['${name}'] = (${params}) => {${_code}};\n`;
      } else {
        me[name] = (input: any) => {
          console.log(input);
        };
      }
    });

    return {sandbox: {me, ...library}, code, attributes};
  }

  private static parseArguments(input: Variable[]) {
    const args: MsgType[] = [];
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
      args.push(E.readVarValue(inp));
    });
    return {args, params};
  }

  private static parseParameters(input: Parameter[]): string[] {
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
    });
    return params;
  }

  private static wrap(output: any, type: Type) {
    const variable = new Variable();
    switch (typeof output) {
      case "object":
        if (type === Type.JSON) {
          variable.setJson(JSON.stringify(output instanceof Error ? {error: output.message} : output));
        } else {
          throw new Error("");
        }
        break;
      case "boolean":
        variable.setBool(output);
        break;
      case "number":
        if (Number.isInteger(output)) {
          if (output < 4294967296) {
            variable.setI32(output);
          } else {
            variable.setI64(output);
          }
          break;
        }
        variable.setF64(output);
        break;
      case "string":
        variable.setString(output);
        break;
    }

    return variable;
  }

  // @ts-ignore
  private static readParams(input: Parameter) {
    switch (input.getType()) {
      case Type.I32:
      case Type.I64:
      case Type.F32:
      case Type.F64:
      case Type.BOOL:
      case Type.JSON:
      case Type.STRING:
      case Type.BINARY:
    }
  }

  private static readVarValue(input: Variable): MsgType {
    // prettier-ignore
    switch (true) {
      case input.hasI32():
        return input.getI32();
      case input.hasI64():
        return input.getI64();
      case input.hasF32():
        return input.getF32();
      case input.hasF64():
        return input.getF64();
      case input.hasBool():
        return input.getBool();
      case input.hasJson():
        return input.getJson();
      case input.hasString():
        return input.getString();
      case input.hasBinary():
        return input.getBinary_asU8();
    }
  }

  private static compareAttributes(me: ThingReference, attributes: ThingReference) {
    for (let i in attributes) {
      if (!_.isEqual(attributes[i], me[i])) {
        // if (attributes[i] !== me[i]) {
        console.log(`>>> ${i}: ${attributes[i]} -> ${me[i]}`);
      }
    }
  }

  exec(script: Script): Result {
    try {
      const fn = script.getFunction()!;
      if (!fn) {
        let json = new Variable().setJson(JSON.stringify({error: "Function is not defined"}));
        return new Result().setOutput(json);
      }

      const {sandbox, code: sandboxCode, attributes} = E.buildSandbox(script);
      const {args, params} = E.parseArguments(fn.getInputList());
      const code = `((${params})=>{try{${fn.getCode()}}catch(_e_){return _e_}})(${args.join()})`;
      console.debug(`${JSON.stringify(sandbox)}>---\n${sandboxCode}\n${code}\n---<`);

      const vm = new VM({sandbox});
      const label = new Date().getTime();
      console.time(`${label}`);
      let output = E.wrap(vm.run(sandboxCode + code), fn.getOutput()!);
      console.timeEnd(`${label}`);
      E.compareAttributes(sandbox.me, attributes);

      return new Result().setOutput(output);
    } catch (e: any) {
      return new Result().setOutput(new Variable().setString(e.message));
    }
  }
}

export {E as Executor};
