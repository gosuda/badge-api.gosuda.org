package main

import (
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type badgeOptions struct {
	Label            string
	Message          string
	Style            string
	LabelColor       string
	MessageColor     string
	LabelTextColor   string
	MessageTextColor string
}

type styleSpec struct {
	Height      float64
	Radius      float64
	FontSize    float64
	Padding     float64
	Uppercase   bool
	BoldLabel   bool
	BoldMessage bool
	Kind        string
}

var badgeStyles = map[string]styleSpec{
	"flat":        {Height: 20, Radius: 3, FontSize: 11, Padding: 7, BoldMessage: true, Kind: "flat"},
	"flat-square": {Height: 20, Radius: 0, FontSize: 11, Padding: 7, BoldMessage: true, Kind: "flat"},
	"plastic":     {Height: 18, Radius: 4, FontSize: 10, Padding: 7, BoldMessage: true, Kind: "plastic"},
	"round":       {Height: 24, Radius: 12, FontSize: 11, Padding: 10, BoldLabel: true, BoldMessage: true, Kind: "flat"},
	"outline":     {Height: 22, Radius: 5, FontSize: 11, Padding: 9, BoldLabel: true, BoldMessage: true, Kind: "outline"},
	"neon":        {Height: 24, Radius: 5, FontSize: 11, Padding: 9, BoldLabel: true, BoldMessage: true, Kind: "neon"},
	"glass":       {Height: 24, Radius: 7, FontSize: 11, Padding: 9, BoldMessage: true, Kind: "glass"},
	"flatbar":     {Height: 28, Radius: 0, FontSize: 10, Padding: 12, Uppercase: true, BoldMessage: true, Kind: "flatbar"},
}

var namedColors = map[string]string{
	"brightgreen":   "44cc11",
	"green":         "97ca00",
	"yellowgreen":   "a4a61d",
	"yellow":        "dfb317",
	"orange":        "fe7d37",
	"red":           "e05d44",
	"blue":          "007ec6",
	"grey":          "555555",
	"gray":          "555555",
	"lightgrey":     "9f9f9f",
	"lightgray":     "9f9f9f",
	"success":       "2f855a",
	"important":     "d97706",
	"critical":      "c53030",
	"informational": "2563eb",
	"inactive":      "718096",
}

func availableStyles() []string {
	return []string{"flat", "flat-square", "plastic", "round", "outline", "neon", "glass", "flatbar"}
}

func normalizeColor(value, fallback string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "#")
	if named, ok := namedColors[value]; ok {
		return named, nil
	}
	if value == "" {
		return fallback, nil
	}
	if len(value) == 3 {
		value = string([]byte{value[0], value[0], value[1], value[1], value[2], value[2]})
	}
	if len(value) != 6 {
		return "", fmt.Errorf("color must be a 3- or 6-digit hex value")
	}
	for _, char := range value {
		if !strings.ContainsRune("0123456789abcdef", char) {
			return "", fmt.Errorf("color contains a non-hex character")
		}
	}
	return value, nil
}

func validateBadgeOptions(options badgeOptions) (badgeOptions, error) {
	options.Label = strings.TrimSpace(options.Label)
	options.Message = strings.TrimSpace(options.Message)
	options.Style = strings.ToLower(strings.TrimSpace(options.Style))
	if options.Style == "" {
		options.Style = "flat"
	}
	if _, ok := badgeStyles[options.Style]; !ok {
		return badgeOptions{}, fmt.Errorf("unknown style %q; supported styles: %s", options.Style, strings.Join(availableStyles(), ", "))
	}
	if options.Message == "" {
		return badgeOptions{}, fmt.Errorf("message is required")
	}
	if utf8.RuneCountInString(options.Label) > 64 {
		return badgeOptions{}, fmt.Errorf("label must be 64 characters or fewer")
	}
	if utf8.RuneCountInString(options.Message) > 128 {
		return badgeOptions{}, fmt.Errorf("message must be 128 characters or fewer")
	}

	var err error
	if options.LabelColor, err = normalizeColor(options.LabelColor, "555555"); err != nil {
		return badgeOptions{}, fmt.Errorf("labelColor: %w", err)
	}
	if options.MessageColor, err = normalizeColor(options.MessageColor, "44cc11"); err != nil {
		return badgeOptions{}, fmt.Errorf("color: %w", err)
	}
	if options.LabelTextColor, err = normalizeColor(options.LabelTextColor, "ffffff"); err != nil {
		return badgeOptions{}, fmt.Errorf("labelTextColor: %w", err)
	}
	if options.MessageTextColor, err = normalizeColor(options.MessageTextColor, "ffffff"); err != nil {
		return badgeOptions{}, fmt.Errorf("textColor: %w", err)
	}
	return options, nil
}

func renderBadge(options badgeOptions) []byte {
	spec := badgeStyles[options.Style]
	label := options.Label
	message := options.Message
	if spec.Uppercase {
		label = strings.ToUpper(label)
		message = strings.ToUpper(message)
	}

	labelTextWidth := measureText(label, spec.FontSize)
	messageTextWidth := measureText(message, spec.FontSize)
	labelWidth := 0.0
	if label != "" {
		labelWidth = math.Ceil(labelTextWidth + spec.Padding*2)
	}
	messageWidth := math.Ceil(messageTextWidth + spec.Padding*2)
	if messageWidth < spec.Height {
		messageWidth = spec.Height
	}
	totalWidth := labelWidth + messageWidth
	aria := message
	if label != "" {
		aria = label + ": " + message
	}

	var builder strings.Builder
	builder.Grow(1400)
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(totalWidth))
	builder.WriteString(`" height="`)
	builder.WriteString(number(spec.Height))
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
	writeTexts(&builder, spec, label, message, labelWidth, messageWidth, labelTextWidth, messageTextWidth, options)
	builder.WriteString(`</svg>`)
	return []byte(builder.String())
}

func writeDefinitions(builder *strings.Builder, spec styleSpec, totalWidth float64, options badgeOptions) {
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

func writeBackgrounds(builder *strings.Builder, spec styleSpec, labelWidth, messageWidth float64, options badgeOptions) {
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

func writeTexts(builder *strings.Builder, spec styleSpec, label, message string, labelWidth, messageWidth, labelTextWidth, messageTextWidth float64, options badgeOptions) {
	y := spec.Height/2 + spec.FontSize*0.34
	builder.WriteString(`<g text-anchor="middle" font-family="Verdana, Geneva, DejaVu Sans, sans-serif" font-size="`)
	builder.WriteString(number(spec.FontSize))
	builder.WriteString(`" text-rendering="geometricPrecision">`)
	if label != "" {
		writeText(builder, label, labelWidth/2, y, labelTextWidth, options.LabelTextColor, spec.BoldLabel)
	}
	writeText(builder, message, labelWidth+messageWidth/2, y, messageTextWidth, options.MessageTextColor, spec.BoldMessage)
	builder.WriteString(`</g>`)
}

func writeText(builder *strings.Builder, text string, x, y, width float64, color string, bold bool) {
	builder.WriteString(`<text x="`)
	builder.WriteString(number(x))
	builder.WriteString(`" y="`)
	builder.WriteString(number(y))
	builder.WriteString(`" textLength="`)
	builder.WriteString(number(width))
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

func measureText(text string, fontSize float64) float64 {
	var width float64
	for _, char := range text {
		switch {
		case char == ' ':
			width += 0.34
		case strings.ContainsRune("ilI1.,'`|!:;", char):
			width += 0.34
		case strings.ContainsRune("mwMW@%&#", char):
			width += 0.92
		case char >= utf8.RuneSelf:
			width += 1.0
		default:
			width += 0.62
		}
	}
	return math.Max(width*fontSize, fontSize*0.5)
}

func number(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
