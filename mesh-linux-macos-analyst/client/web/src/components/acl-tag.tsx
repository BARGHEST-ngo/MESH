// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This file contains code originally from Tailscale (BSD-3-Clause)
// with modifications by BARGHEST. The modified version is licensed
// under AGPL-3.0-or-later. See LICENSE for details.

import cx from "classnames"
import React from "react"
import Badge from "src/ui/badge"

/**
 * ACLTag handles the display of an ACL tag.
 */
export default function ACLTag({
  tag,
  className,
}: {
  tag: string
  className?: string
}) {
  return (
    <Badge
      variant="status"
      color="outline"
      className={cx("flex text-xs items-center", className)}
    >
      <span className="font-medium">tag:</span>
      <span className="text-gray-500">{tag.replace("tag:", "")}</span>
    </Badge>
  )
}
