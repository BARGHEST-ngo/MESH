// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
