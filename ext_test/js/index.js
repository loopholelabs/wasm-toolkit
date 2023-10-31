
// --muximport 
// env/ext_mux
//    0:env/ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_New
//    1:env/ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_HttpConnector_Fetch
//
// --muxexport
// ext_resize
//    0:ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_Resize


let WRITE_BUFFER = new Uint8Array().buffer;

let EXT_PARAM_BUFFER = new Uint8Array().buffer;

function asHex(buff) {
  let a = new Uint8Array(buff);
  let s = "";
  for(let i=0;i<a.length;i++) {
    n = a[i];
    s = s + n.toString(16).padStart(2, "0")
  }
  return s;
}

function fromHex(data) {
  let a = new Uint8Array(data.length / 2);
  for(let i=0;i<(data.length/2);i++) {
    val = parseInt(data.substring(i*2, (i*2)+2), 16);
    a[i] = val;
  }
  return a.buffer;
}

function run() {
	console.log("#JS Run");

  console.log("#JS INPUT " + asHex(WRITE_BUFFER))

  // Call extension

  EXT_PARAM_BUFFER = new Uint8Array(fromHex("0cf601")).buffer;
  let ev = scale_ext_mux([BigInt(0), 0, scale_address_of(EXT_PARAM_BUFFER), EXT_PARAM_BUFFER.byteLength]);
  console.log("#JS NEW RETURN " + ev);

  // Now call fetch...

  EXT_PARAM_BUFFER = new Uint8Array(fromHex("050a1368747470733a2f2f6966636f6e6669672e6d65")).buffer;
  let fv = scale_ext_mux([BigInt(1), ev, scale_address_of(EXT_PARAM_BUFFER), EXT_PARAM_BUFFER.byteLength])
  console.log("#JS FETCH RETURN " + fv);

  console.log("#JS FETCHED " + asHex(EXT_WRITE_BUFFER));

  // Use that in the reply maybe...

  WRITE_BUFFER = new Uint8Array(fromHex("050a3d48656c6c6f20776f726c642066726f6d207761736d20326130303a323363363a363630643a633530313a346238613a663031383a623332343a37346332")).buffer;

  const ptr = scale_address_of(WRITE_BUFFER);
  len = WRITE_BUFFER.byteLength;

  console.log("#JS ptr = " + ptr + " len = " + len);

  return (BigInt(ptr) << BigInt(32)) | BigInt(len);
}

function main() {
	console.log("#JS Main");
}

function initialize() {
	console.log("#JS Initialize");
  return BigInt(0);  
}

function resize(len) {
	console.log("#JS Resize (" + len + ")");

  WRITE_BUFFER = new Uint8Array(len).buffer;

  const ptr = scale_address_of(WRITE_BUFFER);

  console.log("#JS Resize -> " + ptr);

  return ptr;
}

//
EXT_WRITE_BUFFER = new Uint8Array().buffer;

function ext_resize(id, len) {
	console.log("#JS Ext_resize " + id + " (" + len + ")");
  EXT_WRITE_BUFFER = new Uint8Array(len).buffer;
  const ptr = scale_address_of(EXT_WRITE_BUFFER);
  return ptr;
}

exports = {
	run: run,
	main: main,
  resize: resize,
	initialize: initialize,
  ext_resize: ext_resize
}
