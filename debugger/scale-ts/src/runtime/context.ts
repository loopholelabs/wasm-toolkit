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

import { Kind, decodeError } from "@loopholelabs/polyglot-ts";
import { HttpContext as pgContext, HttpRequest as pgRequest, HttpResponse as pgResponse, HttpStringList as pgStringList } from "@loopholelabs/scale-signature-http";

// Wrapper for context
export class Context {
    private _context: pgContext;

    public static time_resize: number = 0;

    constructor(ctx: pgContext) {
        // TODO: Abstract out...
        this._context = ctx;
    }

    context(): pgContext {
        return this._context;
    }

    // Write a context into WebAssembly memory and return a ptr/length
    writeTo(mem: WebAssembly.Memory, mallocfn: Function): {ptr:number, len:number} {
        let inContextBuff = new Uint8Array();
        let encoded = this._context.encode(inContextBuff);

        let ctime_resize = (new Date()).getTime();
        let encPtr = mallocfn(encoded.length);
        let etime_resize = (new Date()).getTime();
        Context.time_resize += (etime_resize - ctime_resize);


//	      const hex = [...encoded].map(x => x.toString(16).padStart(2, '0')).join('');
//	      console.log("writeTo: HEX data = " + hex);

//        console.log("Writing to memory at " + encPtr);
        const memData = new Uint8Array(mem.buffer);
        memData.set(encoded, encPtr);  // Writes the context into memory
        return {ptr: encPtr, len: encoded.length};
    }

    // Read a context from WebAssembly memory
    public static readFrom(mem : WebAssembly.Memory, ptr: number, len: number): Context {
        const memData = new Uint8Array(mem.buffer);

        let inContextBuff = memData.slice(ptr, ptr + len);

//	      const hex = [...inContextBuff].map(x => x.toString(16).padStart(2, '0')).join('');
//	      console.log("readFrom: HEX data = " + hex);


        // Is it an error?
        if (inContextBuff.length > 0 && inContextBuff[0] === Kind.Error) {
          const e = decodeError(inContextBuff).value;
          throw (e);
        }

        let c = pgContext.decode(inContextBuff);
        return new Context(c.value);
    }
}
