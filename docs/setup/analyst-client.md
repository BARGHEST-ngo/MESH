# Analyst client setup

Now let's set up the MESH analyst client on your acquisition workstation. This will allow you to connect to the mesh network and conduct forensic collections.

## Step 1: Clone and build

Clone the MESH repository and build the analyst client if you haven't already:

```bash
# Clone the repository
git clone https://github.com/BARGHEST-ngo/mesh.git
cd mesh

# Build the MESH Docker images
task build
```

!!! tip "Build Time"
    The build process may take 5-10 minutes depending on your system.

## Step 2: Start the MESH analyst container

The MESH analyst container runs in the background and maintains your connection to the mesh network. When first run, you will be prompted to configure the analyst client with the control plane URL and the pre-authentication key you created in the control plane web UI.

```bash
task analyst
```

This will start the analyst container and open an interactive shell in the container. The analyst client will be running in the background and you can use the `meshcli` command to manage your connection to the mesh network, among other things.

## Step 3: Verify Connection

Check that your analyst client is connected to the mesh:

```bash
# Check connection status
meshcli status

# Check your MESH IP
meshcli ip
```

## Troubleshooting

If you encounter issues during setup, see the [Troubleshooting guide](../reference/troubleshooting.md) for common problems and solutions.

## Next steps

Your analyst client is now connected to the mesh. The next step is to install the endpoint client on an Android device.

---

← [Previous: Control plane Setup](control-plane.md) | [Next: Endpoint client Setup](endpoint-client.md) →
