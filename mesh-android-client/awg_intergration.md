# AmneziaWG integration for MESH Android

- Integration of [AmneziaWG 2.0](https://github.com/amnezia-vpn/amneziawg-go/):
-- Replacing the standard WireGuard-Go implementation with the latest AmneziaWG-Go from the master branch.

The integration supports the following obfuscation parameters:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| **Jc** | int | Number of junk packets sent before handshake (0-128) | `Jc = 5` |
| **Jmin** | int | Minimum size of junk packets in bytes | `Jmin = 50` |
| **Jmax** | int | Maximum size of junk packets in bytes (≤1280) | `Jmax = 1000` |
| **S1** | int | Junk bytes added to handshake initiation (15-150) | `S1 = 30` |
| **S2** | int | Junk bytes added to handshake response (15-150) | `S2 = 40` |
| **S3** | int | Junk bytes added to cookie reply packets | `S3 = 25` |
| **S4** | int | Junk bytes added to transport packets | `S4 = 20` |
| **H1** | uint32 or range | Custom message type for handshake initiation | `H1 = 100` or `H1 = 100-200` |
| **H2** | uint32 or range | Custom message type for handshake response | `H2 = 200` or `H2 = 200-300` |
| **H3** | uint32 or range | Custom message type for cookie reply | `H3 = 300` or `H3 = 300-400` |
| **H4** | uint32 or range | Custom message type for transport data | `H4 = 400` or `H4 = 400-500` |

**Note:** AWG 2.0 removed the J1-J4 (junk packet generators) and Itime (interval time) parameters that were present in AWG 1.5. These were rarely used and not part of the final specification. [See more information here](https://github.com/amnezia-vpn/amneziawg-go/pull/91)

## Changes made (most of these are just to get working and still need further code review)

### 1. Dependency Management (`go.mod`)

**File**: `mesh-android/go.mod`

- Added replace directive to substitute `github.com/tailscale/wireguard-go` with local AmneziaWG fork:
  ```go
  replace (
      github.com/tailscale/wireguard-go => ./third_party/amneziawg-go-fork
      gvisor.dev/gvisor => gvisor.dev/gvisor v0.0.0-20250205023644-9414b50a5633
  )
  ```

### 2. AmneziaWG Fork Setup

**Directory**: `mesh-android/third_party/amneziawg-go-fork/`

- Copied AmneziaWG fork from `mesh-linux-macos-client/third_party/amneziawg-go-fork/`
- Modified module path from `github.com/amnezia-vpn/amneziawg-go` to `github.com/tailscale/wireguard-go`
- Updated all internal imports using find-and-replace:
  ```bash
  find third_party/amneziawg-go-fork -name "*.go" -type f -exec sed -i 's|github.com/amnezia-vpn/amneziawg-go|github.com/tailscale/wireguard-go|g' {} +
  ```

### 3. API Compatibility Shims

#### 3.1 Checksum Functions (`tun/checksum.go`)

Added exported wrapper functions for compatibility with `tailscale.com/wgengine/netstack/gro`:

```go
// Checksum computes the Internet checksum (RFC 1071) of the provided bytes.
// Accepts uint16 for compatibility with older tailscale.com versions.
func Checksum(b []byte, initial uint16) uint16 {
    return checksum(b, uint64(initial))
}

// PseudoHeaderChecksum calculates the pseudo-header checksum for TCP/UDP.
func PseudoHeaderChecksum(protocol uint8, srcAddr, dstAddr []byte, totalLen uint16) uint16 {
    sum := pseudoHeaderChecksumNoFold(protocol, srcAddr, dstAddr, totalLen)
    sum = (sum >> 16) + (sum & 0xffff)
    sum = (sum >> 16) + (sum & 0xffff)
    sum = (sum >> 16) + (sum & 0xffff)
    sum = (sum >> 16) + (sum & 0xffff)
    return uint16(sum)
}
```

#### 3.2 conn.Bind Interface Updates (`conn/conn.go`)

Updated `Send` method signature to match Tailscale's wireguard-go fork:

**Before**:
```go
Send(bufs [][]byte, ep Endpoint) error
```

**After**:
```go
Send(bufs [][]byte, ep Endpoint, offset int) error
```

Added new interface types:
- `InitiationAwareEndpoint` - For JIT peer configuration before handshake decryption
- `PeerAwareEndpoint` - For Cryptokey Routing identification awareness

#### 3.3 Send Method Implementations

Updated all `Send` implementations to include `offset int` parameter:

**Files modified**:
- `conn/bind_std.go` - Standard network bind
- `conn/bind_windows.go` - Windows-specific bind
- `conn/bindtest/bindtest.go` - Test bind

**Example**:
```go
func (s *StdNetBind) Send(bufs [][]byte, endpoint Endpoint, offset int) error {
    // Note: offset parameter is ignored in this implementation for compatibility
    // with tailscale.com's wireguard-go fork. AmneziaWG doesn't use offset.
    // ... existing implementation ...
}
```

#### 3.4 Send Call Sites

Updated call sites to pass `offset` parameter:

**Files modified**:
- `device/peer.go:147` - Changed to `peer.device.net.bind.Send(buffers, endpoint, 0)`
- `device/send.go:271` - Changed to `device.net.bind.Send([][]byte{junkedHeader}, initiatingElem.endpoint, 0)`

## Build Process

### Successful Build Output

```bash
make libtailscale
```

**Result**:
- `libtailscale.aar` (56MB) - Stripped production library
- `libtailscale_unstripped.aar` (64MB) - Debug symbols included
- Contains native libraries for all Android architectures:
  - `arm64-v8a` (26MB) - 64-bit ARM (most modern devices)
  - `armeabi-v7a` (31MB) - 32-bit ARM
  - `x86_64` (35MB) - 64-bit x86 (emulators)
  - `x86` (32MB) - 32-bit x86 (older emulators)

## DEtails

1. **Module Path Conflict**: `gomobile` doesn't support replace directives that create dual module paths. So I had to changed the fork's module path to match the original.

2. **API**: Tailscale's wireguard-go fork has evolved beyond the upstream WireGuard-Go that AmneziaWG is based on:
   - Added `offset` parameter to `Send()` for additional encapsulation support
   - Added `InitiationAwareEndpoint` and `PeerAwareEndpoint` interfaces for advanced peer management

3. **Checksum Compatibility**: The `tailscale.com/wgengine/netstack/gro` package expects checksum functions in the `tun` package, but AmneziaWG only has internal (unexported) versions.

### Backward Compatibility

- When all AmneziaWG obfuscation parameters are zero/default, the protocol functions as standard WireGuard
- The `offset` parameter is ignored in AmneziaWG's Send implementations (always 0)
- No changes required to Android UI code - integration is transparent at the Go layer

## Next Steps

### Configuration Layer (Not Yet Implemented)

To fully support AmneziaWG obfuscation parameters, the following changes are needed:

1. **Extend wgengine/wgcfg** (in tailscale.com):
   - Add AmneziaWG parameter fields (Jc, Jmin, Jmax, S1, S2, H1-H4) to config structures
   - Update parser to accept AmneziaWG-specific fields
   - Update writer to serialize AmneziaWG parameters

2. **Configuration File Support**:
   - Implement config file reader for `/data/data/com.barghest.mesh/files/amneziawg.conf`
   - Add validation for obfuscation parameters
   - Apply parameters to WireGuard device on startup

3. **Android UI Integration** (optional):
   - Add UI controls for AmneziaWG parameters
   - Provide presets for common DPI evasion scenarios

## Files

### Core Integration
- `mesh-android/go.mod` - Dependency management
- `mesh-android/third_party/amneziawg-go-fork/go.mod` - Module path change
- All `*.go` files in `third_party/amneziawg-go-fork/` - Import path updates

### API Compatibility
- `third_party/amneziawg-go-fork/tun/checksum.go` - Exported checksum functions
- `third_party/amneziawg-go-fork/conn/conn.go` - Interface updates
- `third_party/amneziawg-go-fork/conn/bind_std.go` - Send signature update
- `third_party/amneziawg-go-fork/conn/bind_windows.go` - Send signature update
- `third_party/amneziawg-go-fork/conn/bindtest/bindtest.go` - Send signature update
- `third_party/amneziawg-go-fork/device/peer.go` - Send call site update
- `third_party/amneziawg-go-fork/device/send.go` - Send call site update
