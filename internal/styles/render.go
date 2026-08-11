package styles

import (
	"html"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/mattn/go-runewidth"
)

var badgeTextWidth = &runewidth.Condition{
	EastAsianWidth:     false,
	StrictEmojiNeutral: false,
}

func renderStandard(options Options, label, message string, spec styleSpec) []byte {
	labelTextWidth := measureText(label, spec.FontSize, options.LetterSpacing)
	messageTextWidth := measureText(message, spec.FontSize, options.LetterSpacing)
	labelWidth := 0.0
	if label != "" {
		labelWidth = math.Ceil(labelTextWidth + spec.Padding*2)
	}
	messageWidth := math.Ceil(messageTextWidth + spec.Padding*2)
	if messageWidth < spec.Height {
		messageWidth = spec.Height
	}
	totalWidth := labelWidth + messageWidth
	aria := accessibleLabel(label, message)

	scale := float64(options.Size) / 100
	var builder strings.Builder
	builder.Grow(1400)
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(totalWidth * scale))
	builder.WriteString(`" height="`)
	builder.WriteString(number(spec.Height * scale))
	builder.WriteString(`" viewBox="0 0 `)
	builder.WriteString(number(totalWidth))
	builder.WriteByte(' ')
	builder.WriteString(number(spec.Height))
	builder.WriteString(`" role="img" aria-label="`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`"><title>`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`</title>`)
	writeDefinitions(&builder, spec, totalWidth, options)

	if spec.Radius > 0 {
		builder.WriteString(`<g clip-path="url(#badge-clip)">`)
	} else {
		builder.WriteString(`<g shape-rendering="crispEdges">`)
	}
	writeBackgrounds(&builder, spec, labelWidth, messageWidth, options)
	builder.WriteString(`</g>`)
	writeTexts(&builder, spec, label, message, labelWidth, messageWidth, labelTextWidth, messageTextWidth, options.LetterSpacing, options)
	builder.WriteString(`</svg>`)
	return []byte(builder.String())
}

func writeDefinitions(builder *strings.Builder, spec styleSpec, totalWidth float64, options Options) {
	builder.WriteString(`<defs>`)
	if spec.Radius > 0 {
		builder.WriteString(`<clipPath id="badge-clip"><rect width="`)
		builder.WriteString(number(totalWidth))
		builder.WriteString(`" height="`)
		builder.WriteString(number(spec.Height))
		builder.WriteString(`" rx="`)
		builder.WriteString(number(spec.Radius))
		builder.WriteString(`"/></clipPath>`)
	}
	if spec.Kind == "plastic" || spec.Kind == "glass" {
		builder.WriteString(`<linearGradient id="shine" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#fff" stop-opacity=".38"/><stop offset=".48" stop-color="#fff" stop-opacity=".08"/><stop offset=".52" stop-color="#000" stop-opacity=".04"/><stop offset="1" stop-color="#000" stop-opacity=".18"/></linearGradient>`)
	}
	if spec.Kind == "neon" {
		builder.WriteString(`<filter id="label-glow" x="-30%" y="-50%" width="160%" height="200%"><feDropShadow dx="0" dy="0" stdDeviation="1.4" flood-color="#`)
		builder.WriteString(options.LabelColor)
		builder.WriteString(`" flood-opacity=".8"/></filter><filter id="message-glow" x="-30%" y="-50%" width="160%" height="200%"><feDropShadow dx="0" dy="0" stdDeviation="1.4" flood-color="#`)
		builder.WriteString(options.MessageColor)
		builder.WriteString(`" flood-opacity=".8"/></filter>`)
	}
	builder.WriteString(`</defs>`)
}

func writeBackgrounds(builder *strings.Builder, spec styleSpec, labelWidth, messageWidth float64, options Options) {
	if labelWidth > 0 {
		writeRect(builder, 0, labelWidth, spec.Height, options.LabelColor, spec.Kind, "label-glow")
	}
	writeRect(builder, labelWidth, messageWidth, spec.Height, options.MessageColor, spec.Kind, "message-glow")

	switch spec.Kind {
	case "plastic", "glass":
		builder.WriteString(`<rect width="100%" height="100%" fill="url(#shine)"/>`)
		if spec.Kind == "glass" {
			builder.WriteString(`<rect x=".5" y=".5" width="calc(100% - 1px)" height="calc(100% - 1px)" rx="`)
			builder.WriteString(number(math.Max(0, spec.Radius-0.5)))
			builder.WriteString(`" fill="none" stroke="#fff" stroke-opacity=".32"/>`)
		}
	case "outline":
		builder.WriteString(`<rect x=".75" y=".75" width="calc(100% - 1.5px)" height="calc(100% - 1.5px)" rx="`)
		builder.WriteString(number(math.Max(0, spec.Radius-0.75)))
		builder.WriteString(`" fill="none" stroke="#fff" stroke-opacity=".55" stroke-width="1.5"/>`)
	}
}

func writeRect(builder *strings.Builder, x, width, height float64, color, kind, filterID string) {
	builder.WriteString(`<rect x="`)
	builder.WriteString(number(x))
	builder.WriteString(`" width="`)
	builder.WriteString(number(width))
	builder.WriteString(`" height="`)
	builder.WriteString(number(height))
	builder.WriteString(`" fill="#`)
	builder.WriteString(color)
	builder.WriteByte('"')
	if kind == "neon" {
		builder.WriteString(` filter="url(#`)
		builder.WriteString(filterID)
		builder.WriteString(`)"`)
	}
	builder.WriteString(`/>`)
}

func writeTexts(builder *strings.Builder, spec styleSpec, label, message string, labelWidth, messageWidth, labelTextWidth, messageTextWidth, letterSpacing float64, options Options) {
	y := spec.Height/2 + spec.FontSize*0.34
	builder.WriteString(`<g text-anchor="middle" font-family="Verdana, Geneva, DejaVu Sans, sans-serif" font-size="`)
	builder.WriteString(number(spec.FontSize))
	builder.WriteString(`" text-rendering="geometricPrecision">`)
	if label != "" {
		writeText(builder, label, labelWidth/2, y, labelTextWidth, options.LabelTextColor, letterSpacing, spec.BoldLabel)
	}
	writeText(builder, message, labelWidth+messageWidth/2, y, messageTextWidth, options.MessageTextColor, letterSpacing, spec.BoldMessage)
	builder.WriteString(`</g>`)
}

func writeText(builder *strings.Builder, text string, x, y, width float64, color string, letterSpacing float64, bold bool) {
	builder.WriteString(`<text x="`)
	builder.WriteString(number(x))
	builder.WriteString(`" y="`)
	builder.WriteString(number(y))
	builder.WriteString(`" textLength="`)
	builder.WriteString(number(width))
	builder.WriteString(`" lengthAdjust="spacing"`)
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

func measureText(text string, fontSize, letterSpacing float64) float64 {
	width := float64(badgeTextWidth.StringWidth(text)) * 0.62
	for _, char := range text {
		if char >= utf8.RuneSelf {
			continue
		}
		switch {
		case char == ' ', strings.ContainsRune("ilI1.,'`|!:;", char):
			width -= 0.28
		case strings.ContainsRune("mwMW@%&#", char):
			width += 0.30
		}
	}
	return math.Max(width*fontSize+letterSpacingWidth(text, letterSpacing), fontSize*0.5)
}

func letterSpacingWidth(text string, letterSpacing float64) float64 {
	if letterSpacing == 0 || text == "" {
		return 0
	}
	clusters := 0
	iterator := graphemes.FromString(text)
	for iterator.Next() {
		clusters++
	}
	return float64(max(0, clusters-1)) * letterSpacing
}

func accessibleLabel(label, message string) string {
	if label == "" {
		return message
	}
	return label + ": " + message
}

func writeFixedSVGStart(builder *strings.Builder, width, height float64, size int, aria string) {
	scale := float64(size) / 100
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(width * scale))
	builder.WriteString(`" height="`)
	builder.WriteString(number(height * scale))
	builder.WriteString(`" viewBox="0 0 `)
	builder.WriteString(number(width))
	builder.WriteByte(' ')
	builder.WriteString(number(height))
	builder.WriteString(`" role="img" aria-label="`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`" shape-rendering="crispEdges"><title>`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`</title><g shape-rendering="crispEdges">`)
}

func number(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
