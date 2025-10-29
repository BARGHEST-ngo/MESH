package internal

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

// StyleBuilder is a function that returns a styled string based on the supplied style.
func StyleBuilder(styles []string, text string, disable bool) (string, error) {
	if len(styles) == 0 {
		return "", fmt.Errorf("style string is empty")
	}

	var err error
	for _, style := range styles {
		text, err = applyStyle(text, style, disable)
		if err != nil {
			return "", err
		}
	}

	return text, nil
}

// applyStyle is a function that apply a style for string.
func applyStyle(text, style string, disable bool) (string, error) {
	if disable {
		return text, nil
	}

	if isHex(style) {
		err := checkHex(style)
		if err != nil {
			return "", err
		}
		text = hexForegroundColor(style, text)
		return text, nil
	}

	if isBackgroundHex(style) {
		rawHex := strings.TrimPrefix(style, "bg")
		err := checkHex(rawHex)
		if err != nil {
			return "", err
		}
		text = hexBackgroundColor(rawHex, text)
		return text, nil
	}

	clr, ok := Map[style]
	if !ok {
		customStyleFunc, ok := CustomMap[style]
		if !ok {
			return "", fmt.Errorf("unknown style '%s'", style)
		}
		text = customStyleFunc(text)
		return text, nil
	}

	return clr.Sprint(text), nil
}

func hexBackgroundColor(cl string, text string) string {
	return color.HEX(cl, true).Sprint(text)
}

func hexForegroundColor(cl string, text string) string {
	return color.HEX(cl).Sprint(text)
}

func isHex(val string) bool {
	return strings.HasPrefix(val, "#")
}

func isBackgroundHex(val string) bool {
	return strings.HasPrefix(val, "bg#")
}

func checkHex(val string) error {
	rgb := color.HexToRGB(val)
	if len(rgb) > 0 {
		return nil
	}

	return fmt.Errorf("invalid '%s' hex color", val)
}
