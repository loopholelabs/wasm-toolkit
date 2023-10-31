package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"HttpFetch"
	sig "signature"

	"github.com/loopholelabs/polyglot"
	scale "github.com/loopholelabs/scale"
	scalefunc "github.com/loopholelabs/scale/scalefunc"
)

type FetchExtension struct {
}

// implementor

type HttpConnector struct {
}

func (fe *FetchExtension) New(c *HttpFetch.HttpConfig) (HttpFetch.HttpConnector, error) {
	fmt.Printf("# New(%v)\n", c)
	return &HttpConnector{}, nil
}

func (hc *HttpConnector) Fetch(u *HttpFetch.ConnectionDetails) (HttpFetch.HttpResponse, error) {
	fmt.Printf("# Fetch(%v)\n", u)

	r := HttpFetch.HttpResponse{}
	// Do the actual fetch here...

	resp, err := http.Get(u.Url)
	if err != nil {
		return r, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	return HttpFetch.HttpResponse{
		StatusCode: int32(resp.StatusCode),
		Body:       body,
	}, nil
}

func main() {
	fmt.Printf("Running scale function with ext...\n")

	sgo, err := scalefunc.Read("local-testfngo-latest.scale")
	if err != nil {
		panic(err)
	}

	f := os.Args[1]

	fmt.Printf("Using wasm [%s]\n", f)

	// Replace the function
	wasm, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	sgo.Function = wasm

	testfn(sgo)

}

func testfn(fn *scalefunc.Schema) {
	fmt.Printf("Running scale function with ext... %s\n", fn.Language)

	ext_impl := &FetchExtension{}

	ctx := context.Background()

	// runtime
	config := scale.NewConfig(sig.New).
		WithContext(ctx).
		WithFunctions([]*scalefunc.Schema{fn}).
		WithExtension(HttpFetch.New(ext_impl)).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRawOutput(true)

	r, err := scale.New(config)
	if err != nil {
		panic(err)
	}

	i, err := r.Instance(nil)
	if err != nil {
		panic(err)
	}

	sigctx := sig.New()

	sigctx.Context.MyString = "hello world"

	b := polyglot.NewBuffer()
	sigctx.Context.Encode(b)
	src := b.Bytes()
	h := hex.EncodeToString(src)
	fmt.Printf("INPUT IS %s\n", h)

	err = i.Run(context.Background(), sigctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Data[%s] from scaleFunction: %s\n", fn.Language, sigctx.Context.MyString)

	b = polyglot.NewBuffer()
	sigctx.Context.Encode(b)
	src = b.Bytes()
	h = hex.EncodeToString(src)
	fmt.Printf("OUTPUT IS %s\n", h)

}
