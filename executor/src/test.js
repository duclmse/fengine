const {VM, VMScript} = require('vm2');

const vm = new VM();
const script = new VMScript('Math.rndom()');
script.compile()
console.log("dasds");
console.log(vm.run(script));
console.log(vm.run(script));
