# MESH Documentation

This directory contains the documentation for MESH (Mesh Forensics Platform), built with [MkDocs](https://www.mkdocs.org/) and the [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) theme.

## Overview

MESH is a comprehensive forensics platform designed for endpoint analysis, network monitoring, and security operations. This documentation provides guides for installation, getting started, and advanced configuration.

## Documentation Structure

- **Overview** - Introduction, architecture, features, platform support, and use cases
- **Getting Started** - Quick start guides for control plane, endpoint clients, and analyst clients
- **Configuration** - Detailed installation instructions for all components
- **User Guide** - Usage patterns and workflows
- **Advanced** - Advanced configuration topics including AmneziaWG and control plane details
- **Reference** - CLI reference and troubleshooting guides

## Local Development

### Prerequisites

- Python 3.x
- pipx

### Installation

```bash
pipx install mkdocs
pipx inject mkdocs mkdocs-material
```

### Running Locally

```bash
mkdocs serve
```

The documentation will be available at `http://localhost:8000`

### Building

To build the static site:

```bash
mkdocs build
```

Output will be generated in the `site/` directory.

## Live Documentation

The documentation is automatically deployed to [https://docs.meshforensics.org/](https://docs.meshforensics.org/) when changes are pushed to the main branch.

## Contributing

To contribute to the documentation:

1. Make changes to the markdown files in the `docs/` directory
2. Test locally with `mkdocs serve`
3. Submit a pull request
