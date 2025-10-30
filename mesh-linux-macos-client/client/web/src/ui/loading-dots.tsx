// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

import cx from "classnames"
import React, { HTMLAttributes } from "react"

type Props = HTMLAttributes<HTMLDivElement>

/**
 * LoadingDots provides a set of horizontal dots to indicate a loading state.
 * These dots are helpful in horizontal contexts (like buttons) where a spinner
 * doesn't fit as well.
 */
export default function LoadingDots(props: Props) {
  const { className, ...rest } = props
  return (
    <div className={cx(className, "loading-dots")} {...rest}>
      <span />
      <span />
      <span />
    </div>
  )
}
