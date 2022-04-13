import {Function, Result, Script, Variable} from "../pb/fengine_pb";
import * as library from "../sdk/db";
import {VM} from "vm2";
import _ from "lodash";

type MsgType = void | number | string | boolean | Uint8Array;
type Func = (input: any) => MsgType;
type ThingReference = {
  [key: string]: MsgType | Func
}

export class Executor {
  private script: Script;

  exec(script: Script): Result {
    try {
      this.script = script;
      const fn = script.getFunction();
      const {sandbox, code: sandboxCode, attributes} = this.buildSandbox();
      const {args, params} = Executor.parseInput(fn.getInputList());
      const code = `((${params})=>{try{${fn.getCode()}}catch(_e_){return _e_}})(${args.join()})`;
      console.dir(sandbox);
      console.log(`>---\n${sandboxCode}\n${code}\n---<`);

      const vm = new VM({sandbox});
      console.time("exec");
      let output = Executor.wrap(vm.run(sandboxCode + code));
      console.timeEnd("exec");
      console.log(`done! out:`);
      this.compareAttributes(sandbox.me, attributes);

      return new Result().setOutput(output);
    } catch (e) {
      return new Result().setOutput(new Variable().setString(e.message));
    }
  }

  buildSandbox() {
    let code = "";
    const me: ThingReference = {};
    const attributes: ThingReference = {};

    this.script.getAttributesList().forEach(attr => {
      let name = attr.getName();
      let value: MsgType = Executor.readVarValue(attr);
      me[name] = value;
      attributes[name] = value;
    });

    this.script.getRefereeMap().forEach((fn: Function, name: string) => {
      let _code = fn.getCode();
      if (_code) {
        const {params} = Executor.parseInput(fn.getInputList(), false);
        code += `me['${name}'] = (${params}) => {${_code}};\n`;
      } else {
        me[name] = (input: any) => {
          console.log(input);
        };
      }
    });

    return {sandbox: {me, ...library}, code, attributes};
  }

  private static parseInput(input: Variable[], hasArgs: boolean = true) {
    const args: MsgType[] = [];
    const params: string[] = [];
    input.forEach(inp => {
      params.push(inp.getName());
      if (hasArgs) {
        args.push(Executor.readVarValue(inp, true));
      } else {
        console.log(`pi: ${inp.getName()}`);
      }
    });
    if (!hasArgs) console.log(`> ${JSON.stringify(args)}`);
    return {args, params};
  }

  private static readVarValue(input: Variable, isParam: boolean = false): MsgType {
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
        return isParam ? `'${input.getString()}'` : input.getString();
      case input.hasBinary():
        return input.getBinary_asU8();
    }
  }

  private static wrap(output: any) {
    const variable = new Variable();
    switch (typeof output) {
      case "object":
        if (output instanceof Error) {
          console.log(`error: ${output.message}`);
          variable.setJson(JSON.stringify({error: output.message}));
        } else {
          variable.setJson(JSON.stringify(output));
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

  compareAttributes(me: ThingReference, attributes: ThingReference) {
    for (let i in attributes) {
      if (!_.isEqual(attributes[i], me[i])) {
        console.log(`>>> ${i}: ${attributes[i]} -> ${me[i]}`);
      }
    }
  }
}
