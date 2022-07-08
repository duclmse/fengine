import {Error, Function, MethodInfo, Parameter, Result, Script, Type, Variable} from "../pb/fengine_pb";
import * as library from "../sdk/db";
import {wrap} from "../sdk/utils";
import {Cache} from "./cache";
import {VM} from "vm2";
import _ from "lodash";

type MsgType = void | number | string | boolean | Uint8Array;
type Func = (input: any) => MsgType;
type ThingReference = {
  [key: string]: MsgType | Func
}
type Obj = {
  [key: string]: any
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

    script.getAttributesList().forEach(attr => {
      let name = attr.getName();
      let value: MsgType = E.readVarValue(attr);
      me[name] = value;
      attributes[name] = value;
    });

    script.getServicesMap().forEach((fn: Function, name: string) => {
      let _code = fn.getCode();
      if (_code) {
        const params = E.parseParams(fn.getInputList());
        code += `me['${name}'] = async ({${params}}) => {${_code}};\n`;
      } else {
        me[name] = (input: any) => {
          console.log(`${name}(${input})`);
        };
      }
    });

    let p: Obj = {};
    let $input = script.getMethod()?.getInputList()?.reduce((p, e) => {
      p[e.getName()] = this.readVarValue(e);
      return p;
    }, p);

    return {sandbox: {me, ...Object.freeze(library), $input}, code, attributes};
  }

  static parseParams(input: Parameter[] | Variable[]): string[] {
    return input.map(inp => inp.getName());
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

  static compareAttributes(me: ThingReference, attributes: ThingReference) {
    const attrs: Variable[] = [];
    for (let i in attributes) {
      if (!_.isEqual(attributes[i], me[i])) {
        console.log(`>>> ${i}: ${attributes[i]} -> ${me[i]}`);
        attrs.push(wrap(me[i], i));
      }
    }
    return attrs;
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

  async exec(script: Script): Promise<Result> {
    try {
      const fn = script.getMethod()!;
      if (!fn) {
        let json = new Variable().setJson(JSON.stringify({error: "Function is not defined"}));
        return new Result().setOutput(json);
      }

      const {sandbox, code: sandboxCode, attributes} = E.buildSandbox(script);
      const params = E.parseParams(fn.getInputList());
      const code = `(async({${params}})=>{try{return me.${fn.getName()}({${params}})}catch(_e_){return _e_}})($input)`;
      console.debug(`${JSON.stringify(sandbox)}\n<---\n${sandboxCode}${code}\n--->`);

      const vm = new VM({sandbox});
      const label = new Date().getTime();
      console.time(`${label}`);
      let output = await vm.run(sandboxCode + code);
      let wrappedOutput = wrap(output);
      console.timeEnd(`${label}`);
      console.log(`${typeof output} -> ${JSON.stringify(output)}`);
      let attrs = E.compareAttributes(sandbox.me, attributes);

      return new Result().setOutput(wrappedOutput).setAttributesList(attrs);
    } catch (e: any) {
      return new Result().setError(new Error().setCode(1).setMessage(e.message));
    }
  }

  upsertService(request: Script): Result {
    this.cache.set(new MethodInfo(), new Function());
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
