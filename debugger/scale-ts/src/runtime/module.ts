/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
// import { argv, env } from 'node:process';
import { Host } from './host';
import { Context } from "./context"

import { WASI } from 'wasi';

import { openSync } from 'node:fs';

let stackPrefix = "";
let functionArgsIndex = 0;
let functionHeader = "";
let functionArgDetails: string[] = [];

let showmem = false;

let memoryData: Uint8Array[] = [];

let functionNames: string[] = [];
 functionNames[0] = "$__wasm_call_ctors";
 functionNames[1] = "$_*github.com/loopholelabs/polyglot-go.Buffer_.Bytes";
 functionNames[2] = "$_*github.com/loopholelabs/polyglot-go.Buffer_.Len";
 functionNames[3] = "$_*github.com/loopholelabs/polyglot-go.Buffer_.Reset";
 functionNames[4] = "$_*github.com/loopholelabs/polyglot-go.Buffer_.Write";
 functionNames[5] = "$runtime.sliceAppend";
 functionNames[6] = "$runtime.slicePanic";
 functionNames[7] = "$github.com/loopholelabs/polyglot-go.decodeUint32";
 functionNames[8] = "$runtime.lookupPanic";
 functionNames[9] = "$github.com/loopholelabs/polyglot-go.decodeString";
 functionNames[10] = "$runtime.alloc";
 functionNames[11] = "$_*github.com/loopholelabs/polyglot-go.Decoder_.Bytes";
 functionNames[12] = "$runtime.nilPanic";
 functionNames[13] = "$_*github.com/loopholelabs/polyglot-go.Decoder_.Error";
 functionNames[14] = "$_*github.com/loopholelabs/polyglot-go.Decoder_.Map";
 functionNames[15] = "$_*github.com/loopholelabs/polyglot-go.Decoder_.Nil";
 functionNames[16] = "$_*github.com/loopholelabs/polyglot-go.Decoder_.String";
 functionNames[17] = "$github.com/loopholelabs/polyglot-go.encodeUint32";
 functionNames[18] = "$github.com/loopholelabs/polyglot-go.encodeString";
 functionNames[19] = "$_*github.com/loopholelabs/polyglot-go.encoder_.Bytes";
 functionNames[20] = "$_*github.com/loopholelabs/polyglot-go.encoder_.Map";
 functionNames[21] = "$_*github.com/loopholelabs/polyglot-go.encoder_.Nil";
 functionNames[22] = "$github.com/loopholelabs/scale-signature-http.NewHttpRequest";
 functionNames[23] = "$github.com/loopholelabs/scale-signature-http.NewHttpResponse";
 functionNames[24] = "$github.com/loopholelabs/scale-signature-http.NewHttpStringList";
 functionNames[25] = "$_*github.com/loopholelabs/scale-signature-http.HttpStringList_.decode";
 functionNames[26] = "$_*github.com/loopholelabs/scale-signature-http.HttpStringList_.internalEncode";
 functionNames[27] = "$_*github.com/loopholelabs/scale-signature-http.Context_.GuestContext";
 functionNames[28] = "$_*github.com/loopholelabs/scale-signature-http.Context_.Response";
 functionNames[29] = "$runtime.memequal";
 functionNames[30] = "$runtime.hash32";
 functionNames[31] = "$runtime.runtimePanic";
 functionNames[32] = "$runtime.printstring";
 functionNames[33] = "$runtime.printnl";
 functionNames[34] = "$runtime.putchar";
 functionNames[35] = "$runtime.markRoots";
 functionNames[36] = "$_runtime.gcBlock_.state";
 functionNames[37] = "$_runtime.gcBlock_.markFree";
 functionNames[38] = "$runtime.growHeap";
 functionNames[39] = "$runtime.startMark";
 functionNames[40] = "$_runtime.gcBlock_.setState";
 functionNames[41] = "$runtime.calculateHeapAddresses";
 functionNames[42] = "$runtime.looksLikePointer";
 functionNames[43] = "$_runtime.gcBlock_.findHead";
 functionNames[44] = "$runtime.hashmapMake";
 functionNames[45] = "$runtime.fastrand";
 functionNames[46] = "$runtime.hashmapStringPtrHash";
 functionNames[47] = "$runtime.hashmapStringEqual";
 functionNames[48] = "$malloc";
 functionNames[49] = "$runtime.hashmapBinarySet";
 functionNames[50] = "$runtime.hashmapSet";
 functionNames[51] = "$runtime.nilMapPanic";
 functionNames[52] = "$runtime.hashmapNext";
 functionNames[53] = "$runtime.hashmapGet";
 functionNames[54] = "$free";
 functionNames[55] = "$runtime.hashmapBinaryGet";
 functionNames[56] = "$runtime.hashmapBinaryDelete";
 functionNames[57] = "$runtime._panic";
 functionNames[58] = "$runtime.printitf";
 functionNames[59] = "$calloc";
 functionNames[60] = "$realloc";
 functionNames[61] = "$runtime.hashmapStringSet";
 functionNames[62] = "$_start";
 functionNames[63] = "$runtime.stringConcat";
 functionNames[64] = "$run";
 functionNames[65] = "$interface:_ErrorWriteBuffer:func:_named:error__basic:uint32_basic:uint32__FromReadBuffer:func:___named:error__ToWriteBuffer:func:___basic:uint32_basic:uint32__.ErrorWriteBuffer$invoke";
 functionNames[66] = "$resize";


export class Module {
    private _code: Buffer;
    private _wasmMod: WebAssembly.Module;
    private _next: Module | null;

    private time_start: number;
    private time_run: number;

    constructor(code: Buffer, next: Module | null) {
      this._code = code;
      this._next = next;
      this._wasmMod = new WebAssembly.Module(this._code);        
      this.time_start = 0;
      this.time_run = 0;
    }

    ShowStats() {
      console.log("Start:  " + this.time_start);
      console.log("Resize: " + Context.time_resize);
      console.log("Run:    " + this.time_run);
    }

    // Run this module, with an optional next module
    run(context: Context): Context {

        // For now, send it in via stdin.
        // TODO: Fix this...
        const fd = openSync('example.in', 'r');

        const wasi = new WASI({
            stdin: fd,
            args: [],
            env: {},
        });

        let wasmModule: WebAssembly.Instance;

        let nextModule = this._next;
        const WASI_ESUCCESS = 0;
        const WASI_EBADF = 8;
        const WASI_EINVAL = 28;
        const WASI_ENOSYS = 52;
    
        let showMemDiff = (fname: string) => {
          let prevMem = memoryData.pop();
          if (prevMem==undefined) {
            console.log("Error!")
          } else {

            // Now do a diff of the memory with current memory...
            const mem = wasmModule.exports.memory as WebAssembly.Memory;
            const memData = new Uint8Array(mem.buffer);

            for(let i=0;i<memData.length; i+=16) {
              let changed = false;
              for (let j=0;j<16;j++) {
                if (memData[i+j]!=prevMem[i+j]) {
                  changed = true;
                }
              }

              if (changed) {
                // Show some memory data...
                let prev = getAsciiHex(prevMem.slice(i, i+16));
                let now = getAsciiHex(memData.slice(i, i+16));
                console.log("Memory Diff " + fname + " " + i.toString(16).padStart(10, '0') + " " + prev.hex + " | " + prev.ascii + " -> " + now.hex + " | " + now.ascii);
              }
            }
          }
        }

        let getAsciiHex = (data: Uint8Array): {hex:string, ascii:string} => {
	        const hex = [...data].map(x => x.toString(16).padStart(2, '0')).join('');

          let ascii:string = ""

          const allowed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-=+!$%^&*()\"'@#~[]{},.<>?|";

          for(let i=0;i<data.length;i++) {
            if (allowed.indexOf(String.fromCharCode(data[i]))!=-1) {
              ascii = ascii + String.fromCharCode(data[i]);
              } else {
              ascii = ascii + " ";
            }
          }

          return {hex: hex, ascii: ascii};
        }

        let getMemoryData = (p: number): {hex:string, ascii:string} => {
          const decoder = new TextDecoder();
          const mem = wasmModule.exports.memory as WebAssembly.Memory;
          const memData = new Uint8Array(mem.buffer);
          const data = memData.slice(p, p + 16);    // For now, just do 16 bytes

          return getAsciiHex(data);
        }

        const importObject = {
          wasi_snapshot_preview1: /*
          {
            proc_exit: (rval: number) => {
              console.log("WASI proc_exit");          
              return WASI_ENOSYS;
            },
            fd_read: () => {
              return WASI_EBADF;
            },
            fd_seek: (fd: number, offset: number, whence: number, newOffsetPtr: number) => {
              console.log("WASI fd_seek");
            },
            fd_fdstat_get: (fd: number, bufPtr: number) => {
              console.log("WASI fd_fdstat_get");
              return WASI_ESUCCESS;
            },
            fd_close: (fd: number) => {
              console.log("WASI fd_close");
              return WASI_ENOSYS;
            },
            clock_time_get: () => {
              console.log("WASI clock_time_get");
              return WASI_ESUCCESS;
            },
            environ_sizes_get: (environCount: number, environBufSize: number) => {
              console.log("WASI environ_sizes_get");
              return WASI_ESUCCESS;
            },
            environ_get:(environ: number, environBuf: number) => {
              console.log("WASI environ_get");
              return WASI_ESUCCESS;
            },
            fd_write: (fd: number, iovs: number, iovsLen: number, nwritten: number) => {
              // Look at the memory, and update the number written...
              const mem = wasmModule.exports.memory as WebAssembly.Memory;
              const memData = new Uint8Array(mem.buffer);

              // Get the io vectors
              let bytesWritten = 0;
              let iovs_ptr = iovs;

              for (let vec = 0;vec < iovsLen; vec ++) {
                const v = memData.slice(iovs_ptr, iovs_ptr + 8);
                const dv = new DataView(v.buffer);
                const ptr = dv.getUint32(0, true);
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
            }

            },
            */  
            wasi.wasiImport,
//            wasi_snapshot_preview1: wasi.exports,

            env: {
                debug_sf: function(id: number) {
                  functionHeader = functionNames[id] + "(";
                  functionArgsIndex = 0;
                  functionArgDetails = [];
                },
                debug_sfp_i32: function(id: number, v: number) {
                  if (functionArgsIndex>0) functionHeader = functionHeader + ", ";
                  functionHeader = functionHeader + "i32:" + v;
                  if (showmem) {
                    let memdata = getMemoryData(v);
                    functionArgDetails.push("[" + functionArgsIndex + "] i32:" + v.toString(16).padStart(10) + " -> " + memdata.hex + " | " + memdata.ascii);
                  }
                  functionArgsIndex++;
                },
                debug_sfp_i64: function(id: number, v: number) {
                  if (functionArgsIndex>0) functionHeader = functionHeader + ", ";
                  functionHeader = functionHeader + "i64:" + v;
                  functionArgsIndex++;
                },
                debug_sfp_f32: function(id: number, v: number) {
                  if (functionArgsIndex>0) functionHeader = functionHeader + ", ";
                  functionHeader = functionHeader + "f32:" + v;
                  functionArgsIndex++;
                },
                debug_sfp_f64: function(id: number, v: number) {
                  if (functionArgsIndex>0) functionHeader = functionHeader + ", ";
                  functionHeader = functionHeader + "f64:" + v;
                  functionArgsIndex++;
                },
                debug_sfp: function(id: number) {
                  console.log("--DEBUG-- " + stackPrefix + "->" + functionHeader + ")");
                  for (let ll of functionArgDetails) {
                    console.log("          " + stackPrefix + "  " + ll);
                  }
                  // TODO: Now show some of the memory details?
                  stackPrefix = stackPrefix + " ";

                  // Push the memory onto our stack
                  const mem = wasmModule.exports.memory as WebAssembly.Memory;
                  const memBufferSnap = new ArrayBuffer(mem.buffer.byteLength);
                  new Uint8Array(memBufferSnap).set(new Uint8Array(mem.buffer));
                  let memSnap = new Uint8Array(memBufferSnap);
                  memoryData.push(memSnap);

                },
                debug_ef: function(id: number) {
                  stackPrefix = stackPrefix.substring(0, stackPrefix.length - 1);
                  console.log("--DEBUG-- " + stackPrefix + "<-" + functionNames[id]);
                  if (showmem) showMemDiff(functionNames[id]);
                },
                debug_ef_i32: function(v: number, id: number): number {
                  stackPrefix = stackPrefix.substring(0, stackPrefix.length - 1);
                  console.log("--DEBUG-- " + stackPrefix + "<-" + functionNames[id] + " i32:" + v);
                  if (showmem) {
                    let memdata = getMemoryData(v);
                    console.log("          " + stackPrefix + "i32:" + v.toString(16).padStart(10) + " -> " + memdata.hex + " | " + memdata.ascii);
                    showMemDiff(functionNames[id]);
                  }
                  return v;
                },
                debug_ef_i64: function(v: number, id: number): number {
                  stackPrefix = stackPrefix.substring(0, stackPrefix.length - 1);
                  console.log("--DEBUG-- " + stackPrefix + "<-" + functionNames[id] + " i64:" + v);
                  if (showmem) showMemDiff(functionNames[id]);
                  return v;
                },
                debug_ef_f32: function(v: number, id: number): number {
                  stackPrefix = stackPrefix.substring(0, stackPrefix.length - 1);
                  console.log("--DEBUG-- " + stackPrefix + "<-" + functionNames[id] + " f32:" + v);
                  if (showmem) showMemDiff(functionNames[id]);
                  return v;
                },
                debug_ef_f64: function(v: number, id: number): number {
                  stackPrefix = stackPrefix.substring(0, stackPrefix.length - 1);
                  console.log("--DEBUG-- " + stackPrefix + "<-" + functionNames[id] + " f64:" + v);
                  if (showmem) showMemDiff(functionNames[id]);
                  return v;
                },
                next: function(ptr: number, len: number): BigInt {
                    const mem = wasmModule.exports.memory as WebAssembly.Memory;
//                    console.log("Next called with " + ptr + " / " + len);
                    
                    if (nextModule != null) {
                        let c = Context.readFrom(mem, ptr, len);
                        let rc = nextModule.run(c);
                        let v = rc.writeTo(mem, wasmModule.exports.malloc as Function);
                        return Host.packMemoryRef(v.ptr, v.len);
                    } else {
                      let packed = Host.packMemoryRef(ptr, len);
//			                console.log("No next module, just return input");
//                      console.log("Next got " + ptr + " " + len + " => " + packed);
                      return packed;
                    }
                }
            }
        };

        wasmModule = new WebAssembly.Instance(this._wasmMod, importObject);

//        console.log("EXPORTS", wasmModule.exports);

        let ctime_start = (new Date()).getTime();

        console.log("Calling wasi.start");
        
        wasi.start(wasmModule);

        let etime_start = (new Date()).getTime();
        this.time_start += (etime_start - ctime_start);

        console.log("wasi.start took " + (etime_start - ctime_start).toFixed(2) + "\n");

        let retContext: Context = context;

        const ITERATIONS = 1;

        for(let ii=0;ii<ITERATIONS;ii++) {

//          console.log("Using resize to allocate and write to memory...");
          const mem = wasmModule.exports.memory as WebAssembly.Memory;
          
          let v = context.writeTo(mem, wasmModule.exports.resize as Function);
//          console.log("Called resize and written to memory " + v.ptr + ", " + v.len);

//          console.log("Calling run function with input " + v.ptr + " " + v.len);
          let ctime_run = (new Date()).getTime();
          const runfn = wasmModule.exports.run as Function;
          let packed = runfn(v.ptr, v.len);
          let etime_run = (new Date()).getTime();
          this.time_run += (etime_run - ctime_run);

  //        console.log("Returned from run " + v.ptr + ", " + v.len + " => " + packed + " took " + (etime_run - ctime_run).toFixed(2));
          let [outContextPtr, outContextLen] = Host.unpackMemoryRef(packed);

          
          retContext = Context.readFrom(mem, outContextPtr, outContextLen);
        }
//        console.log("Return ptr=" + outContextPtr + " len=" + outContextLen);

        return retContext;
    }
}
