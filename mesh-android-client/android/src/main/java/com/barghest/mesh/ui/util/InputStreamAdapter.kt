// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.util

import java.io.InputStream

class InputStreamAdapter(private val inputStream: InputStream) : libtailscale.InputStream {
  override fun read(): ByteArray? {
    val b = ByteArray(4096)
    val i = inputStream.read(b)
    if (i == -1) {
      return null
    }
    return b.sliceArray(0 ..< i)
  }

  override fun close() {
    inputStream.close()
  }
}
