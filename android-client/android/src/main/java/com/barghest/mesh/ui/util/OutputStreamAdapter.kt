// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

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
