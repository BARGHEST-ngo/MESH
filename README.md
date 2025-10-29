<h1 align="center">MESH</h1>

<div align="center">
  <p>
    <a href="https://meshforensics.org/docs">
      <img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="Documentation" />
    </a>
    <a href="https://discord.com/invite/">
      <img src="https://img.shields.io/discord/1161119546170687619?logo=discord&style=flat-square" alt="Chat" />
    </a>
  </p>
</div>


<div align="center">
  <h3>
    <a href="https://meshforensics.org">
      MESH Docs
    </a>
    <span> | </span>
    <a href="https://Barghest.asia">
      BARGHEST
    </a>
  </h3>
</div>
<br/>

## What is MESH

MESH is a peer-to-peer encrypted overlay network that uses UDP-hole punching for secure live state mobile forensics & network monitoring when physical or direct network access is not possible.
It creates and tearsdown end-to-end mesh networks between a locked down analyst host and untrusted Android or iOS devices in a matter of seconds.

Build a forensics mesh and acquire the data in seconds, teardown down and restart without complicated configurations. 
This is remote forensic & network capture without hub-and-spoke topology. 
Remote investigations, that just work. 

Key functions:

- Builds P2P subnets for device analysis.
- End-to-end encrypts using WireGuard.
- Enabled a control plane to manage mesh networks between analysts, control mesh ACLs, and restrict access for compromised nodes on connection.
- Advertise LAN routes or an exit nodes for network monitoring or further forensics use cases. 
- Allows for ephemeral deployment and teardown on disconnection.
- Creates a virtual TUN interface (assigning an address from CGNAT ranges)
- Transfers forensic artifacts such as ADB bug reports and dumpsys data.
- Enables rapid creation, isolation, and teardown of remote investigation nodes.
