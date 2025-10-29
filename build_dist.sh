#!/usr/bin/env sh

set -eu

go="go"
if [ -n "${TS_USE_TOOLCHAIN:-}" ]; then
	go="./tool/go"
fi

eval `CGO_ENABLED=0 GOOS=$($go env GOHOSTOS) GOARCH=$($go env GOHOSTARCH) $go run ./cmd/mkversion`

if [ "$#" -ge 1 ] && [ "$1" = "shellvars" ]; then
	cat <<EOF
VERSION_MINOR="$VERSION_MINOR"
VERSION_SHORT="$VERSION_SHORT"
VERSION_LONG="$VERSION_LONG"
VERSION_GIT_HASH="$VERSION_GIT_HASH"
EOF
	exit 0
fi

tags="${TAGS:-}"
ldflags="-X tailscale.com/version.longStamp=${VERSION_LONG} -X tailscale.com/version.shortStamp=${VERSION_SHORT}"

# build_dist.sh arguments must precede go build arguments.
while [ "$#" -gt 1 ]; do
	case "$1" in
	--extra-small)
		if [ ! -z "${TAGS:-}" ]; then
			echo "set either --extra-small or \$TAGS, but not both"
			exit 1
		fi
		shift
		ldflags="$ldflags -w -s"
		tags="${tags:+$tags,},$(GOOS= GOARCH= $go run ./cmd/featuretags --min --add=osrouter)"
		;;
	--min)
	    # --min is like --extra-small but even smaller, removing all features,
		# even if it results in a useless binary (e.g. removing both netstack +
		# osrouter). It exists for benchmarking purposes only.
		shift
		ldflags="$ldflags -w -s"
		tags="${tags:+$tags,},$(GOOS= GOARCH= $go run ./cmd/featuretags --min)"
		;;
	--box)
		if [ ! -z "${TAGS:-}" ]; then
			echo "set either --box or \$TAGS, but not both"
			exit 1
		fi
		shift
		tags="${tags:+$tags,}ts_include_cli"
		;;
	--custom-tailscaled)
		if [ ! -z "${TAGS:-}" ]; then
			echo "set either --custom-tailscaled or \$TAGS, but not both"
			exit 1
		fi
		shift
		# Custom build with most features but omitting cloud/platform integrations, CLI, and specialized features
		tags="${tags:+$tags,}ts_omit_aws,ts_omit_cloud,ts_omit_kube,ts_omit_synology,ts_omit_appconnectors,ts_omit_cli,ts_omit_completion,ts_omit_cliconndiag,ts_omit_clientupdate,ts_omit_c2n,ts_omit_oauthkey,ts_omit_outboundproxy,ts_omit_peerapiclient,ts_omit_peerapiserver,ts_omit_portlist,ts_omit_relayserver,ts_omit_wakeonlan,ts_omit_tap,ts_omit_bird"
		;;
	*)
		break
		;;
	esac
done

exec $go build ${tags:+-tags=$tags} -trimpath -ldflags "$ldflags" "$@"
