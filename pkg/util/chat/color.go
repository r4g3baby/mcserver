package chat

import (
	"fmt"
	"math"
	"strconv"
)

type Color struct {
	Code string
	Name string
	Hex  string
}

var (
	// ColorChar is the special character which prefixes all chat color codes.
	ColorChar = "\u00A7"

	Black       = Color{"0", "black", "000000"}
	DarkBlue    = Color{"1", "dark_blue", "0000aa"}
	DarkGreen   = Color{"2", "dark_green", "00aa00"}
	DarkAqua    = Color{"3", "dark_aqua", "00aaaa"}
	DarkRed     = Color{"4", "dark_red", "aa0000"}
	DarkPurple  = Color{"5", "dark_purple", "aa00aa"}
	Gold        = Color{"6", "gold", "ffaa00"}
	Gray        = Color{"7", "gray", "aaaaaa"}
	DarkGray    = Color{"8", "dark_gray", "555555"}
	Blue        = Color{"9", "blue", "5555ff"}
	Green       = Color{"a", "green", "55ff55"}
	Aqua        = Color{"b", "aqua", "55ffff"}
	Red         = Color{"c", "red", "ff5555"}
	LightPurple = Color{"d", "light_purple", "ff55ff"}
	Yellow      = Color{"e", "yellow", "ffff55"}
	White       = Color{"f", "white", "ffffff"}

	Obfuscated    = Color{"k", "obfuscated", ""}
	Bold          = Color{"l", "bold", ""}
	Strikethrough = Color{"m", "strikethrough", ""}
	Underline     = Color{"n", "underline", ""}
	Italic        = Color{"o", "italic", ""}
	Reset         = Color{"r", "reset", ""}

	Colors = []Color{
		Black, DarkBlue, DarkGreen, DarkAqua, DarkRed, DarkPurple, Gold, Gray,
		DarkGray, Blue, Green, Aqua, Red, LightPurple, Yellow, White,
	}
)

func (color *Color) String() string {
	if color.Hex != "" {
		return fmt.Sprint("#", color.Hex)
	} else if color.Name != "" {
		return color.Name
	} else {
		return fmt.Sprint(ColorChar, color.Code)
	}
}

func (color *Color) RGB() (r int64, g int64, b int64) {
	if color.Hex != "" {
		r, _ = strconv.ParseInt(color.Hex[:2], 16, 10)
		g, _ = strconv.ParseInt(color.Hex[2:4], 16, 18)
		b, _ = strconv.ParseInt(color.Hex[4:], 16, 10)
	}
	return r, g, b
}

func (color *Color) Distance(other Color) float64 {
	cR, cG, cB := color.RGB()
	oR, oG, oB := other.RGB()
	return math.Sqrt(float64(sq(cR-oR) + sq(cG-oG) + sq(cB-oB)))
}

func FindByName(name string) Color {
	for _, potential := range Colors {
		if potential.Name == name {
			return potential
		}
	}
	return Black
}

func FindNearest(color Color) Color {
	match := Black
	matchDist := math.MaxFloat64
	for _, potential := range Colors {
		distance := color.Distance(potential)
		if distance < matchDist {
			match = potential
			matchDist = distance
		}
		if distance == 0 {
			break
		}
	}
	return match
}

func sq(v int64) int64 {
	return v * v
}
