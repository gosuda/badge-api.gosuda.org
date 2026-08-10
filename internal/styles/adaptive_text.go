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

func adaptiveFontSize(text string, box adaptiveTextBox) float64 {
	if text == "" || box.Width <= 0 || box.Height <= 0 || box.MaxFontSize <= 0 {
		return 0
	}
	widthAtOnePixel := measureText(text, 1)
	if widthAtOnePixel <= 0 {
		return 0
	}
	fontSize := math.Min(box.MaxFontSize, box.Height*0.82)
	fontSize = math.Min(fontSize, box.Width/widthAtOnePixel)
	return math.Max(fontSize, 0.5)
}

func writeAdaptiveText(builder *strings.Builder, text string, box adaptiveTextBox, color string, uppercase bool, bold bool) {
	if uppercase {
		text = strings.ToUpper(text)
	}
	fontSize := adaptiveFontSize(text, box)
	if fontSize == 0 {
		return
	}
	textWidth := math.Min(measureText(text, fontSize), box.Width)
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
	builder.WriteString(`" lengthAdjust="spacingAndGlyphs" fill="#`)
	builder.WriteString(color)
	builder.WriteByte('"')
	if bold {
		builder.WriteString(` font-weight="700"`)
	}
	builder.WriteByte('>')
	builder.WriteString(html.EscapeString(text))
	builder.WriteString(`</text>`)
}
