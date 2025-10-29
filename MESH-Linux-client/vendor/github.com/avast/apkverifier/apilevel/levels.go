package apilevel

import (
	"fmt"
	"math"
)

// FIXME: Would break API :(
// type Level int32

// https://source.android.com/setup/start/build-numbers
const (
	V_AnyMin int32 = -1            // Minimum sdk version, if you don't care about the lower bound
	V_AnyMax int32 = math.MaxInt32 // Maximum sdk version, if you don't care about the upper bound

	V1_0_InitialRelease    int32 = 1
	V1_5_Cupcake           int32 = 3
	V1_6_Donut             int32 = 4
	V2_0_Eclair            int32 = 5
	V2_0_1_Eclair          int32 = 6
	V2_1_Eclair            int32 = 7
	V2_2_Froyo             int32 = 8
	V2_3_Gingerbread       int32 = 9
	V2_3_3_Gingerbread     int32 = 10
	V3_0_Honeycomb         int32 = 11
	V3_1_Honeycomb         int32 = 12
	V3_2_Honeycomb         int32 = 13
	V4_0_1_ICS             int32 = 14
	V4_0_3_ICS             int32 = 15
	V4_1_JellyBean         int32 = 16
	V4_2_JellyBean         int32 = 17
	V4_3_JellyBean         int32 = 18
	V4_4_KitKat            int32 = 19
	V5_0_Lollipop          int32 = 21
	V5_1_Lollipop          int32 = 22
	V6_0_Marshmallow       int32 = 23
	V7_0_Nougat            int32 = 24
	V7_1_Nougat            int32 = 25
	V8_0_Oreo              int32 = 26
	V8_1_Oreo              int32 = 27
	V9_0_Pie               int32 = 28
	V10_0_Ten              int32 = 29
	V11_0_Eleven           int32 = 30
	V12_0_S                int32 = 31
	V12_1_S_V2             int32 = 32
	V13_0_TIRAMISU         int32 = 33
	V14_0_UPSIDE_DOWN_CAKE int32 = 34
)

func SupportsSigV2(level int32) bool {
	return level >= V7_0_Nougat
}

func SupportsSigV3(level int32) bool {
	return level >= V9_0_Pie
}

func SupportsSigV31(level int32) bool {
	return level >= V13_0_TIRAMISU
}

func RequiresSandboxV2(level int32) bool {
	return level >= V8_0_Oreo
}

func SupportsStampVerification(level int32) bool {
	return level >= V11_0_Eleven
}

func String(level int32) string {
	switch level {
	case V1_5_Cupcake:
		return "V1_5_Cupcake"
	case V1_6_Donut:
		return "V1_6_Donut"
	case V2_0_Eclair:
		return "V2_0_Eclair"
	case V2_0_1_Eclair:
		return "V2_0_1_Eclair"
	case V2_1_Eclair:
		return "V2_1_Eclair"
	case V2_2_Froyo:
		return "V2_2_Froyo"
	case V2_3_Gingerbread:
		return "V2_3_Gingerbread"
	case V2_3_3_Gingerbread:
		return "V2_3_3_Gingerbread"
	case V3_0_Honeycomb:
		return "V3_0_Honeycomb"
	case V3_1_Honeycomb:
		return "V3_1_Honeycomb"
	case V3_2_Honeycomb:
		return "V3_2_Honeycomb"
	case V4_0_1_ICS:
		return "V4_0_1_ICS"
	case V4_0_3_ICS:
		return "V4_0_3_ICS"
	case V4_1_JellyBean:
		return "V4_1_JellyBean"
	case V4_2_JellyBean:
		return "V4_2_JellyBean"
	case V4_3_JellyBean:
		return "V4_3_JellyBean"
	case V4_4_KitKat:
		return "V4_4_KitKat"
	case V5_0_Lollipop:
		return "V5_0_Lollipop"
	case V5_1_Lollipop:
		return "V5_1_Lollipop"
	case V6_0_Marshmallow:
		return "V6_0_Marshmallow"
	case V7_0_Nougat:
		return "V7_0_Nougat"
	case V7_1_Nougat:
		return "V7_1_Nougat"
	case V8_0_Oreo:
		return "V8_0_Oreo"
	case V8_1_Oreo:
		return "V8_1_Oreo"
	case V9_0_Pie:
		return "V9_0_Pie"
	case V10_0_Ten:
		return "V10_0_Ten"
	case V11_0_Eleven:
		return "V11_0_Eleven"
	case V12_0_S:
		return "V12_0_S"
	case V12_1_S_V2:
		return "V12_1_S_V2"
	case V13_0_TIRAMISU:
		return "V13_0_TIRAMISU"
	case V14_0_UPSIDE_DOWN_CAKE:
		return "V14_0_UPSIDE_DOWN_CAKE"
	case V_AnyMin:
		return "-Infinity"
	case V_AnyMax:
		return "+Infinity"
	}
	return fmt.Sprintf("%d", level)
}
