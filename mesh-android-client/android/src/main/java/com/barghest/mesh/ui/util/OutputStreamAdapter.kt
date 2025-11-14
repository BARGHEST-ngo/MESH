// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.util

import com.barghest.mesh.util.TSLog
import java.io.OutputStream

// This class adapts a Java OutputStream to the libtailscale.OutputStream interface.
class OutputStreamAdapter(private val outputStream: OutputStream) : libtailscale.OutputStream {
  // writes data to the outputStream in its entirety. Returns -1 on error.
  override fun write(data: ByteArray): Long {
    return try {
      outputStream.write(data)
      outputStream.flush()
      data.size.toLong()
    } catch (e: Exception) {
      TSLog.d("OutputStreamAdapter", "write exception: $e")
      -1L
    }
  }

  override fun close() {
    outputStream.close()
  }
}
