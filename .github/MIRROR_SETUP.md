# Auto-Mirror Setup Guide

This repository automatically mirrors `mesh-linux-macos-analyst/` to the separate [BARGHEST-ngo/mesh-analyst-client](https://github.com/BARGHEST-ngo/mesh-analyst-client) repository for Go module compatibility.

## Setup Instructions

### 1. Create a GitHub Personal Access Token (PAT)

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a descriptive name: `MESH Mirror Bot`
4. Set expiration (recommend: No expiration for automation)
5. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `workflow` (Update GitHub Action workflows)
6. Click "Generate token"
7. **Copy the token immediately** (you won't see it again!)

### 2. Add Token to Repository Secrets

1. Go to this repository's Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `MIRROR_TOKEN`
4. Value: Paste the PAT you just created
5. Click "Add secret"

### 3. Test the Workflow

The workflow triggers automatically on:
- Push to `main` branch (when `mesh-linux-macos-analyst/` changes)
- Push to `dev/**` branches
- Creating tags (e.g., `v0.1.2-alpha.1`)

To manually test:
1. Go to Actions tab
2. Select "Mirror mesh-analyst-client" workflow
3. Click "Run workflow"
4. Select branch and click "Run workflow"

### 4. Verify Mirroring

After the workflow runs:
1. Check [BARGHEST-ngo/mesh-analyst-client](https://github.com/BARGHEST-ngo/mesh-analyst-client)
2. Verify the latest commit matches `mesh-linux-macos-analyst/`
3. If you pushed a tag, verify it appears in the target repo

## How It Works

The workflow uses `git subtree split` to:
1. Extract only the `mesh-linux-macos-analyst/` directory
2. Create a temporary branch with just that content
3. Force-push to `BARGHEST-ngo/mesh-analyst-client`
4. Mirror tags when created

This keeps the Go module path clean while maintaining a monorepo structure.

## Troubleshooting

### Workflow fails with "Permission denied"
- Verify `MIRROR_TOKEN` is set correctly in repository secrets
- Ensure the PAT has `repo` scope
- Check the token hasn't expired

### Changes not appearing in target repo
- Check the workflow run in the Actions tab for errors
- Verify the `paths` filter in the workflow matches your changes
- Ensure you pushed to `main` or a `dev/**` branch

### Tag not mirroring
- Ensure the tag was pushed: `git push origin v0.1.2-alpha.1`
- Check the workflow run for tag-specific errors
- Verify the tag format matches `v*` pattern

