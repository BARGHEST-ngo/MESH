#!/bin/bash
set -e

echo "========================================="
echo "Building MESH Android Client with AWG"
echo "========================================="
echo ""

echo "[1/4] Building unstripped AAR (native library)..."
echo "      This takes 5-15 minutes - please be patient!"
echo ""
make build-unstripped-aar

echo ""
echo "[2/4] Extracting debug symbols..."
make debug-symbols

echo ""
echo "[3/4] Building stripped AAR..."
make $(make print-LIBTAILSCALE_AAR)

echo ""
echo "[4/4] Building debug APK..."
make mesh-debug.apk

echo ""
echo "========================================="
echo "✅ Build complete!"
echo "========================================="
echo "APK location: mesh-debug.apk"
ls -lh mesh-debug.apk

