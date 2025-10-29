// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

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
