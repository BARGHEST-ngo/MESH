# `mkctr`: cross platform container builder for go

`mkctr` is a small go binary which uses `GOOS= GOARCH= go build` directly to compile go binaries and then uses [go-containerregistry](https://github.com/google/go-containerregistry) to create and publish the new containers based on the desired platforms.

This is inspired by [ko](https://github.com/google/ko) which is awesome but doesn't support multiple binaries in a single container.

## Usage

```bash
mkctr \
  --base="alpine:latest" \
  --gopaths="\
    tailscale.com/cmd/tailscale:/usr/local/bin/tailscale, \
    tailscale.com/cmd/tailscaled:/usr/local/bin/tailscaled" \
  --tags="latest" \
  --repos="tailscale/tailscale" \
  [--files=foo.txt:/var/lib/foo.txt,bar.txt:/var/lib/bar.txt] \
  [--target=<target>] \ # e.g. flyio, local
  [--push] \
  [--] [<cmd>...]
```

`mkctr` auto discovers `GOOS`/`GOARCH` from the specified base image. If the base image supports multiple platforms, binaries are compiled for each platform as long as it's one of `linux/amd64`, `linux/386`, `linux/arm`, `linux/arm64`. Multi-arch base image must be either an [OCI image index](https://github.com/opencontainers/image-spec/blob/main/image-index.md) or [Docker manifest list](https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-2.md#manifest-list).
`mkctr` produces image of the same media type as the base image and uses the media type of the base image, or of the individual image references in case of a multi-arch image, to determine the media type of the layer it builds.


## Maturity
This is under active development. While Tailscale uses it, backwards compatability is not guaranteed, and some functionality is missing.
