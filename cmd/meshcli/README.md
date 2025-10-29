# Your Custom Tailscale CLI

This is a skeleton implementation of a custom Tailscale CLI that provides all the functionality of the official `tailscale` command with your own customizations.

## Architecture

The CLI is structured similarly to the official Tailscale CLI:

- `main.go` - Entry point that calls the CLI package
- `cli/cli.go` - Main CLI framework and command routing
- `cli/status.go` - Fully implemented status command (example)
- `cli/commands.go` - Skeleton implementations for all other commands

## Communication with tailscaled

The CLI communicates with the `tailscaled` daemon using the same LocalAPI that the official client uses:

```go
// The local client connects to tailscaled via Unix socket
var localClient = local.Client{
    Socket: paths.DefaultTailscaledSocket(),
}

// Example: Get status from daemon
st, err := localClient.Status(ctx)
```

## Available Commands

All major Tailscale CLI commands are included as skeletons:

### Core Connection
- `up` - Connect to Tailscale
- `down` - Disconnect from Tailscale  
- `login` - Log in to a Tailscale account
- `logout` - Log out and expire node key

### Status & Information
- `status` - Show daemon and connection status ✅ **IMPLEMENTED**
- `version` - Show version information ✅ **IMPLEMENTED**
- `whois` - Show machine/user for a Tailscale IP
- `ip` - Show Tailscale IP addresses

### Configuration
- `set` - Change preferences
- `configure` - Configure host features

### Networking
- `ping` - Ping a host on Tailscale network
- `netcheck` - Analyze local network conditions
- `dns` - Query DNS

### File Sharing
- `file` - Send/receive files via Taildrop
- `drive` - Share directories

### Services  
- `serve` - Serve content and local servers
- `funnel` - Serve content to the internet

### SSH
- `ssh` - SSH to a Tailscale machine

### Security
- `lock` - Manage tailnet lock
- `cert` - Get TLS certificates

### Debugging
- `debug` - Debug commands
- `bugreport` - Generate bug report
- `metrics` - Show metrics

### Other Utilities
- `switch` - Switch Tailscale accounts
- `exit-node` - Manage exit nodes
- `update` - Update Tailscale
- `web` - Run web interface
- `licenses` - Show license information
- `id-token` - Fetch OIDC tokens
- `systray` - Run in system tray
- `app-connector` - Manage app connector routes

## Implementation Guide

### 1. Study the Official Implementation

Look at the official Tailscale CLI commands in `cmd/tailscale/cli/` to understand:
- How they use the LocalAPI client
- Command-line argument parsing
- Error handling patterns
- Output formatting

### 2. Implement Commands One by One

Start with the most important commands for your use case:

```go
func upCmd() *ffcli.Command {
    return &ffcli.Command{
        Name: "up",
        Exec: func(ctx context.Context, args []string) error {
            // Parse flags and arguments
            // Call localClient.Start() with options
            // Handle errors and output
            return nil
        },
        FlagSet: newFlagSet("up"),
    }
}
```

### 3. Key LocalAPI Methods

The most important methods you'll use:

```go
// Connection management
localClient.Start(ctx, opts)           // Connect/configure
localClient.StartLoginInteractive(ctx) // Start login
localClient.Logout(ctx)                // Disconnect

// Status and information  
localClient.Status(ctx)                // Get full status
localClient.WhoIs(ctx, ip)            // Get user/machine info

// Configuration
localClient.EditPrefs(ctx, prefs)      // Change settings

// File operations
localClient.PushFile(ctx, target, size, name, reader)
localClient.WaitingFiles(ctx)

// Networking
localClient.Ping(ctx, ip, opts)        // Ping peer
localClient.NetworkCheck(ctx)          // Network analysis
```

### 4. Error Handling

Follow the pattern from the official CLI:

```go
if err != nil {
    if local.IsAccessDeniedError(err) && os.Getuid() != 0 {
        return fmt.Errorf("%v\n\nUse 'sudo mycli %s'", err, command)
    }
    return fmt.Errorf("operation failed: %w", err)
}
```

### 5. Testing

Test your CLI against a running `tailscaled`:

```bash
# Build your CLI
./build_mycli.sh

# Test basic commands
./mycli status
./mycli version
./mycli --help
```

## Building

```bash
# Make the build script executable
chmod +x build_mycli.sh

# Build your custom CLI
./build_mycli.sh
```

## Next Steps

1. **Implement core commands first**: `up`, `down`, `login`, `logout`
2. **Add your custom features**: Extend existing commands or add new ones
3. **Improve error handling**: Add better error messages and recovery
4. **Add tests**: Create unit tests for your command implementations
5. **Package**: Create installation packages for your custom CLI

## Examples

### Adding a Custom Command

```go
func myCustomCmd() *ffcli.Command {
    return &ffcli.Command{
        Name:      "my-feature",
        ShortHelp: "My custom Tailscale feature",
        Exec: func(ctx context.Context, args []string) error {
            // Your custom logic here
            st, err := localClient.Status(ctx)
            if err != nil {
                return err
            }
            
            // Do something with the status
            printf("Custom feature result: %v\n", st.BackendState)
            return nil
        },
        FlagSet: newFlagSet("my-feature"),
    }
}
```

# Start the daemon (requires sudo for TUN interface)
sudo ./tailscaled-nodrive --socket=/tmp/tailscale/tailscaled.sock --state=/tmp/tailscale/tailscaled.state --statedir=/tmp/tailscale --verbose=1

# In another terminal, use your CLI
./mesh status
./mesh ping 100.64.0.5
./mesh ssh 100.64.0.5  # If SSH is enabled