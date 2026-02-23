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
import React, { HTMLAttributes } from "react"

type Props = {
  className?: string
  size: "sm" | "md"
} & HTMLAttributes<HTMLDivElement>

export default function Spinner(props: Props) {
  const { className, size, ...rest } = props

  return (
    <div
      className={cx(
        "spinner inline-block rounded-full align-middle",
        {
          "border-2 w-4 h-4": size === "sm",
          "border-4 w-8 h-8": size === "md",
        },
        className
      )}
      {...rest}
    />
  )
}

Spinner.defaultProps = {
  size: "md",
}
