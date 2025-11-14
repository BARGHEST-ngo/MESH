## Building on arm64 devices

Since we call the LocalAPI differently to the current Tailscale implimentation, CGO-generated arg structs are not properly aligned on arm64 which is the majority of physical devices. 

This causes a 'bulkBarrierPreWrite' crash when the garbage collector runs during the LocalAPI calls.

This fix is already in the Go master branch but not yet in any released version or Tailscale's Go Fork (which is required).

Once Tailscale updates their Go fork to include this fix, this manual patch will no longer be needed. 

### Required patch

After the toolchain is downloaded to `~/.cache/tailscale-go`, apply this patch to 
`~/.cache/tailscale-go/src/cmd/cgo/out.go`:

Find line ~1058 (search for `typedef %s %v _cgo_argtype`):
```go
fmt.Fprintf(fgcc, "\ttypedef %s %v _cgo_argtype;\n", ctype.String(), p.packedAttribute())
```

Replace this with:

```go
fmt.Fprintf(fgcc, "\ttypedef %s %v __attribute__((aligned(8))) _cgo_argtype;\n", ctype.String(), p.packedAttribute())
```

then rebuild the Go toolchain:

cd ~/.cache/tailscale-go/src
./make.bash

