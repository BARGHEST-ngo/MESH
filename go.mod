module github.com/BARGHEST-ngo/MESH

go 1.25.6

require (
	github.com/botherder/go-savetime v1.5.0
	github.com/google/uuid v1.6.0
	github.com/i582/cfmt v1.4.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/mvt-project/androidqf_ward v0.0.0-00010101000000-000000000000
	github.com/peterbourgon/ff/v3 v3.4.0
	github.com/toqueteos/webbrowser v1.2.1
	golang.org/x/net v0.49.0
	golang.org/x/sys v0.40.0
	tailscale.com v1.94.1
)

require (
	filippo.io/age v1.2.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	fyne.io/systray v1.11.1-0.20250812065214-4856ac3adc3c // indirect
	github.com/Kodeworks/golang-image-ico v0.0.0-20141118225523-73f0f4cfade9 // indirect
	github.com/akutz/memconn v0.1.0 // indirect
	github.com/alexbrainman/sspi v0.0.0-20231016080023-1a75b4708caa // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/avast/apkparser v0.0.0-20250626104540-d53391f4d69d // indirect
	github.com/avast/apkverifier v0.0.0-20250626104651-727e33396aec // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.58 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.5 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/coder/websocket v1.8.12 // indirect
	github.com/dblohm7/wingoes v0.0.0-20240119213807-a09d6be7affa // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/gaissmai/bart v0.18.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20250813024750-ebf49471dced // indirect
	github.com/godbus/dbus/v5 v5.1.1-0.20230522191255-76236955d466 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/jsimonetti/rtnetlink v1.4.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mdlayher/netlink v1.7.3-0.20250113171957-fbb4dce95f42 // indirect
	github.com/mdlayher/socket v0.5.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/tailscale/go-winio v0.0.0-20231025203758-c4f33415bf55 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go4.org/mem v0.0.0-20240501181205-ae6ca9944745 // indirect
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/oauth2 v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	golang.zx2c4.com/wireguard/windows v0.5.3 // indirect
)

replace (
	github.com/amnezia-vpn/amneziawg-go => github.com/BARGHEST-ngo/amnezia-wireguard-go v0.1.1-alpha.1
	github.com/mvt-project/androidqf_ward => github.com/BARGHEST-ngo/androidqf_mesh v0.1.0
	gvisor.dev/gvisor => gvisor.dev/gvisor v0.0.0-20250205023644-9414b50a5633
)

tool golang.org/x/tools/go/packages/gopackages
