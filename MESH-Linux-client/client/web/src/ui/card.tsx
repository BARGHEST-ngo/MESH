// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

import cx from "classnames"
import React from "react"

type Props = {
  children: React.ReactNode
  className?: string
  elevated?: boolean
  empty?: boolean
  noPadding?: boolean
}

/**
 * Card is a box with a border, rounded corners, and some padding. Use it to
 * group content into a single container and give it more importance. The
 * elevation prop gives it a box shadow, while the empty prop a light gray
 * background color.
 *
 *     <Card>{content}</Card>
 *     <Card elevated>{content}</Card>
 *     <Card empty><EmptyState description="You don't have any keys" /></Card>
 *
 */
export default function Card(props: Props) {
  const { children, className, elevated, empty, noPadding } = props
  return (
    <div
      className={cx("rounded-md border", className, {
        "shadow-soft": elevated,
        "bg-gray-0": empty,
        "bg-white": !empty,
        "p-6": !noPadding,
      })}
    >
      {children}
    </div>
  )
}
