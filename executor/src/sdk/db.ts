export class TableImpl {
  private readonly name: string;

  constructor(name: string) {
    this.name = name;
  }

  get getName() {
    return this.name;
  }

  Select(filter: Filter) {
    console.log(`calling Select ${JSON.stringify(filter)} from ${this.name}`);
  }

  Insert(row: any) {
    console.log(`calling Insert ${JSON.stringify(row)} into ${this.name}`);
  }

  Update(info: { set: any, where: any }) {
    const {set, where} = info;
    console.log(`calling Update ${set} where ${where} in ${this.name}`);
  }

  Delete(filter: Filter) {
    console.log(`calling Delete ${filter} from ${this.name}`);
  }
}

export class Filter {

}

export function Table(name: string): TableImpl {
  return new TableImpl(name);
}

