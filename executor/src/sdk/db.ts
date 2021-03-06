import {
  DataRow,
  DeleteRequest,
  DeleteResult,
  InsertRequest,
  InsertResult,
  ResultRow,
  SelectRequest,
  SelectResult,
  UpdateRequest,
  UpdateResult,
  Value,
  Variable
} from "../pb/fengine_pb";
import {getClient} from "../executor/grpc_client";
import {unwrapValue, wrap} from "./utils";

export function Table(name: string): TableImpl {
  return new TableImpl(name);
}

export class TableImpl {
  private readonly _name: string;

  constructor(name: string) {
    this._name = name;
  }

  get name() {
    return this._name;
  }

  static toVars(value: RowInfo) {
    let vars: Variable[] = [];
    for (let k in value) {
      vars.push(wrap(value[k], k));
    }
    return vars;
  }

  Select(info: SelectInfo): Promise<ResultSet> {
    let {fieldNames, filter, limit, offset, orderBy, groupBy} = info;
    let req = new SelectRequest()
      .setTable(this._name)
      .setFieldList(fieldNames || ["*"])
      .setFilter(JSON.stringify(filter))
      .setLimit(limit || 1000)
      .setOffset(offset || 0)
      .setOrderByList(orderBy || [])
      .setGroupByList(groupBy || []);
    return new Promise<ResultSet>((resolve, reject) => {
      getClient().select(req, (err, res) => err == null ? resolve(new ResultSet(res)) : reject(err));
    });
  }

  Insert(info: InsertInfo): Promise<InsertResult> {
    let rows = info.rows.map(value => new DataRow().setValuesList(TableImpl.toVars(value)));
    let req = new InsertRequest()
      .setTable(this._name)
      .setRowList(rows);
    return new Promise<InsertResult>((resolve, reject) => {
      getClient().insert(req, (err, res) => err == null ? resolve(res) : reject(err));
    });
  }

  Update(info: UpdateInfo): Promise<UpdateResult> {
    const {row, filter} = info;
    let req = new UpdateRequest()
      .setTable(this._name)
      .setFieldList(TableImpl.toVars(row))
      .setFilter(JSON.stringify(filter));
    return new Promise<UpdateResult>((resolve, reject) => {
      getClient().update(req, (err, res) => err == null ? resolve(res) : reject(err));
    });
  }

  Delete(filter: object): Promise<DeleteResult> {
    let req = new DeleteRequest()
      .setTable(this._name)
      .setFilter(JSON.stringify(filter));
    return new Promise<DeleteResult>((resolve, reject) => {
      getClient().delete(req, (err, res) => err == null ? resolve(res) : reject(err));
    });
  }
}

interface Index {
  [key: string]: number;
}

export class ResultSet {
  private readonly length: number;
  private readonly index: Index = {};
  private readonly data: ResultRow[];
  private readonly cols: string[];
  private row: Value[] | undefined;
  private rowIndex: number;

  constructor(result: SelectResult) {
    let cols = result.getColumnList();
    this.cols = cols;
    this.index = ResultSet.getIndex(cols);
    let data = result.getRowList();
    this.data = data;
    this.rowIndex = -1;
    this.length = data.length;
  }

  private static getIndex(cols: string[]) {
    let index: Index = {};
    for (let i = 0, l = cols.length; i < l; i++) {
      index[cols[i]] = i;
    }
    return index;
  }

  GetColumns = () => this.cols;

  Next(): boolean {
    if (this.rowIndex + 1 < this.length) {
      this.rowIndex++;
      this.row = this.data[this.rowIndex].getValueList();
      return true;
    }
    return false;
  }

  Get(field: string): Value | undefined {
    return this.row ? this.row[this.index[field]] : undefined;
  }

  * [Symbol.iterator]() {
    for (let i of this.data) {
      yield new Row(i.getValueList(), this.index, this.cols);
    }
  }

  filter(callback: (value: Row, index: number) => boolean): Row[] {
    let arr: Row[] = [];
    let i = 0;
    for (let row of this) {
      if (callback(row, i++)) {
        arr.push(row);
      }
    }
    return arr;
  }

  map<U>(callback: (value: Row, index: number) => U): U[] {
    let arr: U[] = [];
    let i = 0;
    for (let row of this) {
      arr.push(callback(row, i++));
    }
    return arr;
  }

  reduce<U>(callback: (previousValue: U, currentValue: Row, index: number) => U, initial: U): U {
    let result = initial;
    let i = 0;
    for (let row of this) {
      result = callback(result, row, i++);
    }
    return result;
  }
}

export class Row {
  private readonly values: Value[];
  private readonly index: Index;
  private readonly cols: string[];

  constructor(values: Value[], index: Index, cols: string[]) {
    this.values = values;
    this.index = index;
    this.cols = cols;
  }

  Get(field: string) {
    return this.values[this.index[field]];
  }

  ToObject() {
    let obj: RowInfo = {};
    for (let i = 0, l = this.cols.length; i < l; i++) {
      obj[this.cols[i]] = this.values[i];
    }
    return obj;
  }

  * [Symbol.iterator]() {
    for (let i of this.values) {
      yield unwrapValue(i);
    }
  }
}

export interface SelectInfo {
  fieldNames: string[] | undefined;
  filter: object | undefined;
  limit: number | undefined;
  offset: number | undefined;
  groupBy: string[] | undefined;
  orderBy: string[] | undefined;
}

export interface InsertInfo {
  rows: RowInfo[];
}

export interface RowInfo {
  [key: string]: any;
}

interface UpdateInfo {
  row: RowInfo;
  filter: object;
}
