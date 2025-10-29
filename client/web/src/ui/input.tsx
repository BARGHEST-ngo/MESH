// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

import cx from "classnames"
import React, { InputHTMLAttributes } from "react"

type Props = {
  className?: string
  inputClassName?: string
  error?: boolean
  suffix?: JSX.Element
} & InputHTMLAttributes<HTMLInputElement>

// Input is styled in a way that only works for text inputs.
const Input = React.forwardRef<HTMLInputElement, Props>((props, ref) => {
  const {
    className,
    inputClassName,
    error,
    prefix,
    suffix,
    disabled,
    ...rest
  } = props
  return (
    <div className={cx("relative", className)}>
      <input
        ref={ref}
        className={cx("input z-10", inputClassName, {
          "input-error": error,
        })}
        disabled={disabled}
        {...rest}
      />
      {suffix ? (
        <div className="bg-white top-1 bottom-1 right-1 rounded-r-md absolute flex items-center">
          {suffix}
        </div>
      ) : null}
    </div>
  )
})

export default Input
