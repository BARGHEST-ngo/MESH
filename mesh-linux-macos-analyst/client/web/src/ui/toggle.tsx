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
import React, { ChangeEvent } from "react"

type Props = {
  id?: string
  className?: string
  disabled?: boolean
  checked: boolean
  sizeVariant?: "small" | "medium" | "large"
  onChange: (checked: boolean) => void
}

export default function Toggle(props: Props) {
  const { className, id, disabled, checked, sizeVariant, onChange } = props

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    onChange(e.target.checked)
  }

  return (
    <input
      id={id}
      type="checkbox"
      className={cx(
        "toggle",
        {
          "toggle-large": sizeVariant === "large",
          "toggle-small": sizeVariant === "small",
        },
        className
      )}
      disabled={disabled}
      checked={checked}
      onChange={handleChange}
    />
  )
}

Toggle.defaultProps = {
  sizeVariant: "medium",
}
