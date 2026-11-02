# Acquisitions and evidence handling

The analyst container collects Android acquisitions with `mesh adbcollect` (AndroidQF/WARD) and can analyse them offline with MVT. This page covers where acquisitions are stored, how to get them off the container, and how to analyse them.

## Where acquisitions are stored

The analyst container persists one path for evidence: `/home/mesh/acquisitions`. It is a **bind mount** from the host, so anything written there appears on the host immediately and survives the container being stopped, recreated, or destroyed.

| | Path |
|---|---|
| **Inside the container** | `/home/mesh/acquisitions` |
| **On the host** | `${ACQUISITION_DIR:-$HOME/mesh-acquisitions}` |

The host location defaults to `mesh-acquisitions` in the host user's home directory. Override it by setting `ACQUISITION_DIR` in your `.env` before starting the container — for example, to a dedicated encrypted evidence volume:

```bash
# .env
ACQUISITION_DIR=/mnt/evidence/mesh-acquisitions
```

!!! warning "Everything else in the container is ephemeral"
    Only `/home/mesh/acquisitions` is bind-mounted. Files written anywhere else — `/tmp`, `/home/mesh`, `/var/captures` — live in the container's writable layer and are **lost** when the container is recreated. Always keep evidence under the acquisitions path.

!!! note "Ownership"
    The container runs as `mesh` (uid/gid 1000) and writes acquisitions as that user, so the host directory must be owned by `1000:1000`. The entrypoint takes ownership of the mount on start if it is not already writable.

## Collecting an acquisition

The container's working directory is the acquisitions mount, so a bare `mesh adbcollect` already writes into it:

```bash
# Writes to /home/mesh/acquisitions/<uuid> — on the host, persisted
mesh adbcollect
```

For a named case directory, pass `--output` explicitly. Keep the path under `/home/mesh/acquisitions` so it lands on the host:

```bash
mesh adbcollect --output /home/mesh/acquisitions/case-2026-08-27
```

!!! danger "Do not point `--output` outside the mount"
    `mesh adbcollect --output /tmp/case1` writes to the container's ephemeral layer and the acquisition is **lost** on recreate. Any explicit `--output` must be under `/home/mesh/acquisitions`.

## Getting acquisitions off the container

Because `/home/mesh/acquisitions` is bind-mounted, **the files are already on the host** the moment they are written — there is usually nothing to "copy out". Retrieval is just a matter of reading them from the host path.

### Primary method — read the bind mount on the host

From the host (not inside the container):

```bash
# Default location
ls -la "${ACQUISITION_DIR:-$HOME/mesh-acquisitions}"

# A specific case
ls -la "${ACQUISITION_DIR:-$HOME/mesh-acquisitions}/case-2026-08-27"
```

Copy, archive, or move it wherever your evidence workflow requires:

```bash
ACQ="${ACQUISITION_DIR:-$HOME/mesh-acquisitions}"

# Hash before moving — establish integrity on the host
find "$ACQ/case-2026-08-27" -type f -exec sha256sum {} \; > case-2026-08-27.sha256

# Archive for transfer
tar -czf case-2026-08-27.tar.gz -C "$ACQ" case-2026-08-27
```

### Alternative — copy from a running container

If you did not use the bind mount (for example an older deployment, or an explicit `--output` elsewhere), copy the files out of the running container:

```bash
# Copy a case directory from the container to the host
docker compose cp analyst:/home/mesh/acquisitions/case-2026-08-27 ./

# Or stream it out as a tar, preserving ownership and permissions
docker compose exec -T analyst tar -C /home/mesh/acquisitions -cf - case-2026-08-27 \
  | tar -C ./ -xf -
```

!!! warning "Only while the container exists"
    `docker compose cp` works only while the container is running. Once it is recreated or destroyed, anything not on the bind mount is gone. The bind mount is the reliable path.

### Chain of custody

- Compute hashes on the **host** after retrieval (`sha256sum`), not inside the container.
- Store acquisitions on an encrypted filesystem and restrict access — the directory holds pulled APKs and a full device backup.
- Record the case identifier, collection time, and analyst alongside the acquisition.

## Analysing with MVT

The analyst image ships [MVT (Mobile Verification Toolkit)](https://github.com/mvt-project/mvt) for offline analysis of acquisitions. Run it inside the container, pointing it at a collected acquisition.

MVT's threat indicators are **not** bundled — download them before the first scan:

```bash
# Refresh public STIX2 indicators
mvt-android download-iocs
```

Then check an acquisition against them:

```bash
# Point MVT at the acquisition directory
mvt-android check-androidqf /home/mesh/acquisitions/case-2026-08-27

# Write MVT's own findings alongside the acquisition so they persist too
mvt-android check-androidqf /home/mesh/acquisitions/case-2026-08-27 \
  --output /home/mesh/acquisitions/case-2026-08-27/mvt
```

!!! note "Indicators are not persisted"
    Downloaded indicators live in `~/.local`, which is not bind-mounted, so `download-iocs` must be re-run after the container is recreated. Record which indicator set a case was scanned against for reproducibility.

## Next steps

- **[Packet capture and exit node](../advanced/exit-node-pcap.md)** — capture network traffic into the same acquisitions directory
- **[Troubleshooting](../reference/troubleshooting.md)** — common collection and connection issues
