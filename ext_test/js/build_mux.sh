#!/bin/bash

cat index.js | ~/code/scale/scale/compile/typescript/builder/target/release/js_builder

# Now add the other bits we need...

../../cmd/cmd customs -i scale.wasm -o module_out.wasm \
--muximport "env/ext_mux,0:env/ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_New,1:env/ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_HttpConnector_Fetch" \
--muxexport "ext_resize,0:ext_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_Resize"
