// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
