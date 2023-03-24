import * as fs from 'fs';
import { Module } from './runtime/module';
import { HttpContext as Context, HttpRequest as Request, HttpResponse as Response, HttpStringList as StringList } from "@loopholelabs/scale-signature-http";

import { Host } from './runtime/host';
import { Context as ourContext} from './runtime/context';

// Create a context to send in...
var enc = new TextEncoder();
let body = enc.encode("Hello world this is a request body");
let headers = new Map<string, StringList>();
headers.set('content', new StringList(['hello']));
let req1 = new Request('GET', "http://something.com", BigInt(100), 'https', '1.2.3.4', body, headers);
let respBody = enc.encode("Response body");
let respHeaders = new Map<string, StringList>();        
const resp1 = new Response(200, respBody, respHeaders);        
const context = new Context(req1, resp1);

// Now we can use context...

let ff = process.argv[2];

console.log("Loading " + ff);

const modHttpEndpoint = fs.readFileSync(ff);
//const modHttpEndpoint = fs.readFileSync('./example_modules/http-endpoint.wasm');
//const modHttpMiddleware = fs.readFileSync('./example_modules/http-middleware.wasm');
let moduleHttpEndpoint = new Module(modHttpEndpoint, null);
//let moduleHttpMiddleware = new Module(modHttpMiddleware, moduleHttpEndpoint);

// Run the modules...

let ctx = new ourContext(context);

console.log("\nINPUT CONTEXT")
Host.showContext(context);

let ctime = (new Date()).getTime();

console.log("Calling run");
let retContext = moduleHttpEndpoint.run(ctx);
console.log("Run returned");

/*

// Do some benchmarking stuff
for(let i=0;i<1000;i++) {
  let retContext = moduleHttpEndpoint.run(ctx);
}

*/
let etime = (new Date()).getTime();

moduleHttpEndpoint.ShowStats();

console.log("TOTAL " + (etime - ctime).toFixed(2));

console.log("\nOUTPUT CONTEXT");
Host.showContext(retContext.context());
