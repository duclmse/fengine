import {Function, MethodId, Parameter, Result, Script, Type, Variable} from "../pb/fengine_pb";
import {Cache} from "./cache";
import * as library from "../sdk/db";
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

  static buildSandbox(script: Script) {
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

  static parseArguments(input: Variable[]) {
    const args: MsgType[] = [];
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
      args.push(E.readVarValue(inp));
    });
    return {args, params};
  }

  static parseParameters(input: Parameter[]): string[] {
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
    });
    return params;
  }

  static wrap(output: any, type: Type) {
    const variable = new Variable();
    switch (typeof output) {
      case "object":
        if (type === Type.JSON) {
          variable.setType(Type.JSON).setJson(
            JSON.stringify(output instanceof Error ? {error: output.message} : output));
        } else {
          throw new Error("");
        }
        break;
      case "boolean":
        variable.setType(Type.BOOL).setBool(output);
        break;
      case "number":
        if (Number.isInteger(output)) {
          if (output < 4294967296) {
            variable.setType(Type.I32).setI32(output);
          } else {
            variable.setType(Type.I64).setI64(output);
          }
          break;
        }
        variable.setType(Type.F64).setF64(output);
        break;
      case "string":
        variable.setType(Type.STRING).setString(output);
        break;
    }

    return variable;
  }

  // @ts-ignore
  static readParams(input: Parameter) {
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

  static readVarValue(input: Variable): MsgType {
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

  static compareAttributes(me: ThingReference, attributes: ThingReference) {
    for (let i in attributes) {
      if (!_.isEqual(attributes[i], me[i])) {
        // if (attributes[i] !== me[i]) {
        console.log(`>>> ${i}: ${attributes[i]} -> ${me[i]}`);
      }
    }
  }

  exec(script: Script): Result {
    return new Result();
  }

  upsertService(request: Script): Result {
    this.cache.set(new MethodId(), new Function());
    const json = JSON.stringify({});
    let variable = new Variable().setType(Type.JSON).setJson(json);
    return new Result().setOutput(variable);
  }

  deleteService(request: Script): Result {
    const json = JSON.stringify({});
    let variable = new Variable().setType(Type.JSON).setJson(json);
    return new Result().setOutput(variable);
  }
}

export {E as Executor};
