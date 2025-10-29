package internal

import (
	"github.com/gookit/color"
)

// CustomMap is the custom styles map
var CustomMap = map[string]func(string) string{}

// Map is the styles map
var Map = map[string]color.Color{
	"red":          color.Red,
	"cyan":         color.Cyan,
	"gray":         color.Gray,
	"blue":         color.Blue,
	"black":        color.Black,
	"green":        color.Green,
	"white":        color.White,
	"yellow":       color.Yellow,
	"magenta":      color.Magenta,
	"lightRed":     color.LightRed,
	"lightCyan":    color.LightCyan,
	"lightBlue":    color.LightBlue,
	"lightGreen":   color.LightGreen,
	"lightWhite":   color.LightWhite,
	"lightYellow":  color.LightYellow,
	"lightMagenta": color.LightMagenta,

	"bgRed":          color.BgRed,
	"bgGray":         color.BgGray,
	"bgCyan":         color.BgCyan,
	"bgBlue":         color.BgBlue,
	"bgWhite":        color.BgWhite,
	"bgBlack":        color.BgBlack,
	"bgGreen":        color.BgGreen,
	"bgYellow":       color.BgYellow,
	"bgMagenta":      color.BgMagenta,
	"bgDefault":      color.BgDefault,
	"bgDarkGray":     color.BgDarkGray,
	"bgLightRed":     color.BgLightRed,
	"bgLightCyan":    color.BgLightCyan,
	"bgLightBlue":    color.BgLightBlue,
	"bgLightWhite":   color.BgLightWhite,
	"bgLightGreen":   color.BgLightGreen,
	"bgLightYellow":  color.BgLightYellow,
	"bgLightMagenta": color.BgLightMagenta,

	"bold":      color.Bold,
	"italic":    color.OpItalic,
	"crossout":  color.OpStrikethrough,
	"underline": color.OpUnderscore,
	"blink":     color.OpBlink,
	"reverse":   color.OpReverse,
	"concealed": color.OpConcealed,
}
