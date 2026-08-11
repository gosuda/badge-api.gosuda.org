package styles

import (
	_ "embed"
	"html"
	"math"
	"strings"
)

const (
	classicWidth      = 88
	classicHeight     = 31
	clickHereHeight   = 31
	clickHereMinWidth = 88
	clickHereMaxWidth = 300
)

//go:embed assets/click-here-88x31.svg
var clickHereReferenceSVG string

var clickHereReferenceArtwork = referenceArtwork(clickHereReferenceSVG)
var clickHereStyle = styleSpec{Height: clickHereHeight, Kind: "click-here", Render: renderClickHere}

func renderClickHere(options Options, label, message string) []byte {
	phrase := strings.TrimSpace(strings.TrimSpace(label) + " " + strings.TrimSpace(message))

	// Use reference for default "click here" text
	if options.LetterSpacing == 0 && strings.EqualFold(phrase, "click here") {
		var builder strings.Builder
		builder.Grow(len(clickHereReferenceArtwork) + 1800)
		writeFixedSVGStart(&builder, clickHereMinWidth, clickHereHeight, options.Size, accessibleLabel(label, message))
		builder.WriteString(clickHereReferenceArtwork)
		builder.WriteString(`</g></svg>`)
		return []byte(builder.String())
	}

	// Variable width for custom text
	textWidth := int(math.Ceil(measureText(strings.ToUpper(phrase), 7, options.LetterSpacing))) + 12
	totalWidth := textWidth + 10 // 5px left cap + text + 5px right cap
	totalWidth = max(clickHereMinWidth, min(totalWidth, clickHereMaxWidth))
	textBoxWidth := totalWidth - 10

	var builder strings.Builder
	builder.Grow(3000)

	aria := accessibleLabel(label, message)
	scale := float64(options.Size) / 100
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(float64(totalWidth) * scale))
	builder.WriteString(`" height="`)
	builder.WriteString(number(clickHereHeight * scale))
	builder.WriteString(`" viewBox="0 0 `)
	builder.WriteString(number(float64(totalWidth)))
	builder.WriteByte(' ')
	builder.WriteString(number(clickHereHeight))
	builder.WriteString(`" role="img" aria-label="`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`"><title>`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`</title><g shape-rendering="crispEdges">`)

	// Background and frame
	writeClickHereFrame(&builder, totalWidth)

	// Text overlay
	builder.WriteString(`<rect x="3" y="3" width="`)
	builder.WriteString(number(float64(textBoxWidth)))
	builder.WriteString(`" height="24" fill="#cccccc"/>`)
	writeAdaptiveText(&builder, phrase, adaptiveTextBox{X: 5, Y: 5, Width: float64(textBoxWidth - 4), Height: 20, MaxFontSize: 11}, "3a0603", options.LetterSpacing, true, true)

	// Red arrow (right side)
	arrowX := totalWidth - 9
	builder.WriteString(`<path fill="#e74436" d="M`)
	builder.WriteString(number(float64(arrowX)))
	builder.WriteString(` 17h2v1h-2zM`)
	builder.WriteString(number(float64(arrowX)))
	builder.WriteString(` 19h2v1h-2z"/>`)

	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}

func writeClickHereFrame(builder *strings.Builder, width int) {
	// Base fill
	builder.WriteString(`<rect width="`)
	builder.WriteString(number(float64(width)))
	builder.WriteString(`" height="31" fill="#cccccc"/>`)

	// White top border
	builder.WriteString(`<path fill="#ffffff" d="M0 0h`)
	builder.WriteString(number(float64(width - 1)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width - 1)))
	builder.WriteString(`zM0 1h1v29h-1z"/>`)

	// Black right and bottom borders
	builder.WriteString(`<path fill="#000000" d="M`)
	builder.WriteString(number(float64(width - 1)))
	builder.WriteString(` 2h1v28h-1zM2 30h`)
	builder.WriteString(number(float64(width - 2)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width - 2)))
	builder.WriteString(`z"/>`)

	// Gray inner borders
	builder.WriteString(`<path fill="#2a2a2a" d="M`)
	builder.WriteString(number(float64(width - 2)))
	builder.WriteString(` 2h1v27h-1z"/>`)
	builder.WriteString(`<path fill="#222222" d="M2 29h`)
	builder.WriteString(number(float64(width - 4)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width - 4)))
	builder.WriteString(`z"/>`)
}
