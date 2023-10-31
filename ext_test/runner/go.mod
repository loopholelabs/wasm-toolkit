module github.com/loopholelabs/scale-ext-example

go 1.20

replace github.com/loopholelabs/scale => /home/jimmy/code/scale/scale

require signature v0.1.0

require (
	HttpFetch v0.0.0-00010101000000-000000000000
	github.com/loopholelabs/scale v0.3.19
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/hashicorp/hcl/v2 v2.18.1 // indirect
	github.com/loopholelabs/polyglot v1.1.3 // indirect
	github.com/loopholelabs/scale-extension-interfaces v0.1.0 // indirect
	github.com/loopholelabs/scale-signature-interfaces v0.1.7 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/tetratelabs/wazero v1.5.0 // indirect
	github.com/zclconf/go-cty v1.14.1 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace signature v0.1.0 => /home/jimmy/.config/scale/signatures/local_testsig_latest_e6ddebc792ee929e2654b4281baca1376e05bf5a96d4bdf63a05a2aab5f9e749_signature/golang/host

replace HttpFetch => /home/jimmy/.config/scale/extensions/local_testext_latest_5c7d22390f9101d459292d76c11b5e9f66c327b1766aae34b9cc75f9f40e8206_extension/golang/host
