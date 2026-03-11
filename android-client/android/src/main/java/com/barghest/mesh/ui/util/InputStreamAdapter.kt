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
