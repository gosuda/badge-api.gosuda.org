package styles

import (
	"html"
	"math"
	"strings"
)

type adaptiveTextBox struct {
	X           float64
	Y           float64
	Width       float64
	Height      float64
	MaxFontSize float64
}

// minReadableFontSize is the smallest size at which Verdana glyphs stay
// recognizable. Badge width grows with the text so this floor is reached only
// by input past the style's maximum width.
const minReadableFontSize = 5.0

func adaptiveFontSize(text string, box adaptiveTextBox, letterSpacing float64) float64 {
	if text == "" || box.Width <= 0 || box.Height <= 0 || box.MaxFontSize <= 0 {
		return 0
	}
	widthAtOnePixel := measureText(text, 1, 0)
	if widthAtOnePixel <= 0 {
		return 0
	}
	availableGlyphWidth := box.Width - letterSpacingWidth(text, letterSpacing)
	fontSize := math.Min(box.MaxFontSize, box.Height*0.82)
	if availableGlyphWidth > 0 {
		fontSize = math.Min(fontSize, availableGlyphWidth/widthAtOnePixel)
	} else {
		fontSize = minReadableFontSize
	}
	return math.Max(fontSize, minReadableFontSize)
}

func writeAdaptiveText(builder *strings.Builder, text string, box adaptiveTextBox, color string, letterSpacing float64, uppercase bool, bold bool) {
	if uppercase {
		text = strings.ToUpper(text)
	}
	fontSize := adaptiveFontSize(text, box, letterSpacing)
	if fontSize == 0 {
		return
	}
	textWidth := math.Min(measureText(text, fontSize, letterSpacing), box.Width)
	x := box.X + box.Width/2
	y := box.Y + box.Height/2 + fontSize*0.34

	builder.WriteString(`<text x="`)
	builder.WriteString(number(x))
	builder.WriteString(`" y="`)
	builder.WriteString(number(y))
	builder.WriteString(`" text-anchor="middle" font-family="Verdana, Geneva, DejaVu Sans, sans-serif" font-size="`)
	builder.WriteString(number(fontSize))
	builder.WriteString(`" text-rendering="geometricPrecision" textLength="`)
	builder.WriteString(number(textWidth))
	builder.WriteString(`" lengthAdjust="spacingAndGlyphs"`)
	if letterSpacing != 0 {
		builder.WriteString(` letter-spacing="`)
		builder.WriteString(number(letterSpacing))
		builder.WriteByte('"')
	}
	builder.WriteString(` fill="#`)
	builder.WriteString(color)
	builder.WriteByte('"')
	if bold {
		builder.WriteString(` font-weight="700"`)
	}
	builder.WriteByte('>')
	builder.WriteString(html.EscapeString(text))
	builder.WriteString(`</text>`)
}
