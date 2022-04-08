const grpc = require("../pb/fengine_grpc_pb");

class TableImpl {
  constructor(name) {
    this.name = name;
  }

  Select(filter) {
    console.log(`calling Select ${filter} into ${this.name}`);
  }

  Insert(row) {
    console.log(`calling Insert ${row} into ${this.name}`);
  }

  Update(info) {
    const {set, where} = info;
    console.log(`calling Update ${set} where ${where} in ${this.name}`);
  }

  Delete(filter) {
    console.log(`calling Delete ${filter} from ${this.name}`);
  }
}

function Table(name) {
  return new TableImpl(name);
}

if (require.main === module) {
  let table = Table("asdsds");
  console.log(table.name);
  return table;
}

module.exports = {
  Table: Table
};
