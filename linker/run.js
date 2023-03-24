var instance;

const WASI_ESUCCESS = 0;
const WASI_EBADF = 8;
const WASI_EINVAL = 28;
const WASI_ENOSYS = 52;

const importObject = {
  wasi_snapshot_preview1: {
    fd_write: (fd, iovs, iovsLen, nwritten) => {

      // Look at the memory, and update the number written...
      const mem = instance.exports.memory;
      const memData = new Uint8Array(mem.buffer);

      // Get the io vectors
      let bytesWritten = 0;
      let iovs_ptr = iovs;

      for (let vec = 0;vec < iovsLen; vec ++) {
        const v = memData.slice(iovs_ptr, iovs_ptr + 8);
        const dv = new DataView(v.buffer);
        let ptr = dv.getUint32(0, true);
        const len = dv.getUint32(4, true);

        iovs_ptr += 8;
        bytesWritten += len;

        const b = memData.slice(ptr, ptr + len);
        const dec = new TextDecoder();
        //const hex = [...b].map(x => x.toString(16).padStart(2, '0')).join('');

        // NOTE: The console.log/error functions add a newline at the end.
        if (fd == 1) {
          console.log(dec.decode(b));
        } else if (fd == 2) {
          console.error(dec.decode(b));
        }
      }

      new DataView(mem.buffer).setUint32(nwritten, bytesWritten, true);
      return WASI_ESUCCESS;
    },

    environ_sizes_get: (environCount, environBufSize)=>{
//      console.log(";; -> environ_sizes_get")
      const mem = instance.exports.memory;
      const dataView = new DataView(mem.buffer);
      dataView.setUint32(environCount, 0, true);
      dataView.setUint32(environBufSize, 0, true);
      return WASI_ESUCCESS
    },

    clock_time_get: ()=>{
//      console.log(";; -> clock_time_get")
      return WASI_ESUCCESS
    },
    args_sizes_get: (argc, argvBufSize)=>{
//      console.log(";; -> args_sizes_get")
      const mem = instance.exports.memory;
      const dataView = new DataView(mem.buffer);
      dataView.setUint32(argc, 0, true);
      dataView.setUint32(argvBufSize, 0, true);
      return WASI_ESUCCESS
    },
    args_get: ()=>{
//      console.log(";; -> args_get")
      return WASI_ESUCCESS
    },
    environ_get: ()=>{
//      console.log(";; -> environ_get")
      return WASI_ESUCCESS
    },
    fd_close: ()=>{
//      console.log(";; -> fd_close")
      return WASI_EBADF
    },
    fd_fdstat_get: ()=>{
//      console.log(";; -> fd_fdstat_get")
      return WASI_EBADF
    },
    fd_prestat_get: ()=>{
//      console.log(";; -> fd_prestat_get")
      return WASI_EBADF
    },
    fd_prestat_dir_name: ()=>{
//      console.log(";; -> fd_prestat_dir_name")
      return WASI_EINVAL
    },
    fd_read: ()=>{
//      console.log(";; -> fd_read")
      return WASI_EBADF
    },
    fd_seek: ()=>{
//      console.log(";; -> fd_seek")
      return WASI_EBADF
    },
    path_open: ()=>{
//      console.log(";; -> path_open")
      return WASI_EBADF
    },
    proc_exit: (rval)=>{
//      console.log(";; -> proc_exit(" + rval + ")")
      return WASI_ENOSYS
    }
  },
  env: {
    h_debug: (n) => {
      console.log("DEBUG: 0x" + n.toString(16) + " " + n);
    },
    mod1_value: () => {
      return 77;
    }
  },
};

const fs = require('fs');

const wasmFile = process.argv[2];

console.log(";; Loading wasm file " + wasmFile);

const wasmBuffer = fs.readFileSync(wasmFile);

WebAssembly.instantiate(wasmBuffer, importObject).then(wasmModule => {
  instance = wasmModule.instance;

  // Exported function live under instance.exports
  const _start1 = wasmModule.instance.exports["mod1._start"];
  const hello1 = wasmModule.instance.exports["mod1.hello"];
  const _start2 = wasmModule.instance.exports["mod2._start"];
  const hello2 = wasmModule.instance.exports["mod2.hello"];
  console.log(";; Running _start1");
  try {
    _start1();
  } catch(e) {
//    console.log("_start1 threw an error", e);
  }
  console.log(";; Running _start2");
  try {
    _start2();
  } catch(e) {
//    console.log("_start2 threw an error", e);
  }

  // Now try running hello()

  console.log(";; Running mod1.hello()");
  const ret1 = hello1();
//  console.log(";; hello function returned ", ret1);

  console.log(";; Running mod2.hello()");
  const ret2 = hello2();
//  console.log(";; hello function returned ", ret2);

}).catch((e) => {
  console.error(e);
});
