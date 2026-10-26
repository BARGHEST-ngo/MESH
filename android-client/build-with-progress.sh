#!/bin/bash
set -e

echo "========================================="
echo "Building MESH Android Client"
echo "========================================="
echo "Building debug APK..."
make mesh-debug.apk
echo ""
echo "========================================="
echo "✅ Build complete!"
echo "========================================="
echo "APK location: mesh-debug.apk"
ls -lh mesh-debug.apk

