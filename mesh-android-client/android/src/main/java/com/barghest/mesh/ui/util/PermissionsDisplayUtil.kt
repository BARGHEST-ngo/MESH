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

import android.net.Uri

/** Converts a SAF URI string to a more human-friendly folder display name. */
fun friendlyDirName(uriStr: String): String {
  val uri = Uri.parse(uriStr)
  val segment = uri.lastPathSegment ?: return uriStr

  return when {
    segment.startsWith("primary:") -> "Internal storage › " + segment.removePrefix("primary:")
    segment.contains(":") -> {
      val folder = segment.substringAfter(":")
      "SD card › $folder"
    }
    else -> segment
  }
}
