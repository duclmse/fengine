import {Type, Value, Variable} from "../pb/fengine_pb";
import ValueCase = Value.ValueCase;

export function wrap(value: any, name: string | null = null): Variable {
  const variable = new Variable();
  if (name) {
    variable.setName(name);
  }
  switch (typeof value) {
    case "object":
      return variable.setType(Type.JSON).setJson(
        JSON.stringify(value instanceof Error ? {error: value.message} : value));
    case "boolean":
      return variable.setType(Type.BOOL).setBool(value);
    case "number":
      if (Number.isInteger(value)) {
        if (value < 4294967296) {
          return variable.setType(Type.I32).setI32(value);
        } else {
          return variable.setType(Type.I64).setI64(value);
        }
      }
      return variable.setType(Type.F64).setF64(value);
    case "string":
      return variable.setType(Type.STRING).setString(value);
    default:
      return variable;
  }
}

export function unwrap(variable: Variable) {
  switch (variable.getType()) {
    case Type.I32:
      return variable.getI32();
    case Type.I64:
      return variable.getI64();
    case Type.F32:
      return variable.getF32();
    case Type.F64:
      return variable.getF64();
    case Type.BINARY:
      return variable.getBinary();
    case Type.BOOL:
      return variable.getBool();
    case Type.JSON:
      return JSON.parse(variable.getJson());
    case Type.STRING:
      return variable.getString();
    default:
      return null;
  }
}

export function unwrapValue(value: Value) {
  switch (value.getValueCase()) {
    case ValueCase.I32:
      return value.getI32();
    case ValueCase.I64:
      return value.getI64();
    case ValueCase.F32:
      return value.getF32();
    case ValueCase.F64:
      return value.getF64();
    case ValueCase.BINARY:
      return value.getBinary();
    case ValueCase.BOOL:
      return value.getBool();
    case ValueCase.JSON:
      return JSON.parse(value.getJson());
    case ValueCase.STRING:
      return value.getString();
    default:
      return null;
  }
}
