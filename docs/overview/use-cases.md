# Use cases

MESH is designed for forensic investigators, human rights defenders, and security researchers operating in challenging environments. Here are real-world scenarios where MESH provides critical capabilities.

## Remote mobile forensics

Investigate Android and iOS devices remotely without physical access. Collect artifacts, run spyware scans, and perform live analysis over an encrypted mesh network.

### Scenario

A human rights organisation needs to analyse the phone of an activist who suspects they've been targeted with spyware. The activist is in a different country and cannot safely mail their device.

### How MESH helps

1. **Deploy endpoint client** - Activist installs MESH app on their device
2. **Enable ADB** - Activist enables ADB-over-Wifi on their device
3. **Establish secure connection** - Device joins the mesh over encrypted tunnel
4. **Remote analysis** - Investigator runs AndroidQF/MVT/WARD to detect spyware
5. **Collect evidence** - Artifacts are securely transferred over the mesh
6. **Tear down** - Mesh connection is terminated

### Benefits

- **No physical access required** - Investigate devices anywhere in the world
- **Secure data transfer** - All forensic data is encrypted end-to-end
- **Real-time analysis** - Perform live forensics on running devices
- **Minimal user interaction** - Simple app installation, no technical expertise needed
- **Low configuration**

## Network monitoring and packet capture

Deploy forensic capabilities to remote locations for network monitoring and evidence collection.

### Scenario

An incident response team needs to monitor network traffic from a compromised device in a remote office to identify command-and-control (C2) communications.

### How MESH helps

1. **Deploy endpoint client** - Install MESH on the compromised device
2. **Enable exit node** - Route device traffic through analyst node
3. **Capture traffic** - Monitor and log all network connections
4. **Analyse patterns** - Identify malicious C2 traffic
5. **Collect evidence** - Preserve network forensics for investigation

### Benefits

- **Remote deployment** - No need to travel to remote locations
- **Real-time monitoring** - Live visibility into device network activity
- **Encrypted collection** - Evidence is protected during transfer
- **Rapid response** - Deploy in minutes, not hours or days

!!! example "Workflow"
    ```bash
    # Enable exit node on analyst node
    sudo ./meshcli up --advertise-exit-node

    # On endpoint, route all traffic through analyst
    # (configured in MESH app)
    
    # Enable port forwarding on analyst 

    # Capture traffic on analyst workstation
    tcpdump -i mesh0 -w capture.pcap
    ```

## Human rights investigations

Securely analyse devices of activists, journalists, and human rights defenders in hostile environments where network traffic is monitored or censored.

### Scenario

A journalist in an authoritarian country suspects their phone has been compromised. They need their device analysed but cannot risk exposing the investigation to state surveillance.

### How MESH Helps

1. **Censorship resistance** - AmneziaWG obfuscation bypasses DPI and VPN blocking
2. **Encrypted communication** - All forensic data is protected from interception
3. **Secure analysis** - Remote investigation without physical device transfer
4. **Evidence preservation** - Forensic artifacts collected securely
5. **Operational security** - No suspicious VPN traffic visible to censors

### Benefits

- **Bypasses censorship** - Works in countries with aggressive internet filtering
- **Protects investigators** - Encrypted connections hide forensic activity
- **Protects subjects** - No need to physically transport devices across borders
- **Maintains evidence integrity** - Secure chain of custody

!!! warning "Operational security"
    When operating in hostile environments, follow proper OPSEC:

    - Use AmneziaWG obfuscation or HTTPS DERP relays (see [configuration guide](../advanced/amneziawg.md))
    - Deploy control plane in a safe jurisdiction with reverse proxy
    - Use HTTPS relay fallback for maximum compatibility
    - Consider using Tor or other anonymity networks for control plane access

## Censored network operations

Operate forensic investigations in countries with aggressive internet censorship using AmneziaWG obfuscation and HTTPS fallback.

### Scenario

A security researcher needs to investigate devices in China, Russia, or Iran where VPN protocols are actively blocked by state-level firewalls.c

### How MESH Helps

1. **AmneziaWG obfuscation** - Makes VPN traffic appear as regular HTTPS
2. **HTTPS relay fallback** - Works when UDP is completely blocked
3. **DPI evasion** - Bypasses Deep packet inspection systems
4. **Reliable connectivity** - Maintains connection despite censorship attempts

### Tested Against

- **Great Firewall of China** - Successfully bypasses DPI and protocol blocking
- **Russian SORM** - Evades traffic analysis and VPN detection
- **Iranian filtering** - Works through national firewall infrastructure
- **Corporate DPI** - Bypasses enterprise-grade packet inspection

### Benefits

- **Reliable access** - Works where traditional VPNs fail
- **Indistinguishable traffic** - Appears as normal HTTPS web browsing
- **Automatic adaptation** - Falls back to relay when P2P is blocked
- **Proven effectiveness** - Tested in real-world censored environments

!!! tip "Configuration for censored networks"
    See the [AmneziaWG Configuration guide](../advanced/amneziawg.md) for detailed setup instructions optimized for censored environments.

### MESH provides in all cases

- **Speed** - Deploy in minutes, not hours
- **No persistent infrastructure** - Ephemeral mesh leaves no trace
- **Secure evidence handling** - Encrypted transfer maintains chain of custody
- **Minimal disruption** - No need to power down or transport devices

!!! example "Incident Response Timeline"
    - **T+0 min**: Incident detected
    - **T+5 min**: Control plane deployed
    - **T+10 min**: Analyst client connected
    - **T+15 min**: Endpoint client installed on compromised device
    - **T+20 min**: Forensic collection begins
    - **T+60 min**: Evidence collected, mesh torn down

---

**Next:** Check [Platform support](platform-support.md) to see which platforms are supported →
