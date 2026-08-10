package styles

import (
	_ "embed"
	"html"
	"math"
	"strconv"
	"strings"
)

//go:embed assets/best-viewed-with-chrome-88x31.svg
var bestViewedReferenceSVG string

var bestViewedReferenceArtwork = referenceArtwork(bestViewedReferenceSVG)
var bestViewedStyle = styleSpec{Height: classicHeight, Kind: "best-viewed", Render: renderBestViewed}

func renderBestViewed(options Options, label, message string) []byte {
	isDefault := strings.EqualFold(strings.TrimSpace(label), "viewed with") &&
		strings.EqualFold(strings.TrimSpace(message), "chrome")

	// Use reference artwork for the default phrase
	if isDefault {
		var builder strings.Builder
		builder.Grow(len(bestViewedReferenceArtwork) + 1800)
		writeFixedSVGStart(&builder, classicWidth, classicHeight, options.Size, accessibleLabel(label, message))
		builder.WriteString(bestViewedReferenceArtwork)
		builder.WriteString(`</g></svg>`)
		return []byte(builder.String())
	}

	// Variable width: 10px red accent + text panel + 28px Chrome logo
	const (
		redPanelWidth      = 10
		logoWidth          = 28
		textPadding        = 4
		bestViewedMinWidth = 88
		bestViewedMaxWidth = 300
	)

	textNeeded := 0
	if label == "" {
		textNeeded = int(math.Ceil(measureText(strings.ToUpper(message), 10)))
	} else {
		labelNeeded := int(math.Ceil(measureText(strings.ToUpper(label), 5)))
		messageNeeded := int(math.Ceil(measureText(strings.ToUpper(message), 8)))
		textNeeded = max(labelNeeded, messageNeeded)
	}
	textPanelWidth := textNeeded + textPadding*2

	totalWidth := redPanelWidth + textPanelWidth + logoWidth
	totalWidth = max(bestViewedMinWidth, min(totalWidth, bestViewedMaxWidth))
	textPanelWidth = totalWidth - redPanelWidth - logoWidth

	var builder strings.Builder
	builder.Grow(3200)

	aria := accessibleLabel(label, message)
	scale := float64(options.Size) / 100
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(float64(totalWidth) * scale))
	builder.WriteString(`" height="`)
	builder.WriteString(number(classicHeight * scale))
	builder.WriteString(`" viewBox="0 0 `)
	builder.WriteString(number(float64(totalWidth)))
	builder.WriteByte(' ')
	builder.WriteString(number(classicHeight))
	builder.WriteString(`" role="img" aria-label="`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`"><title>`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`</title><g shape-rendering="crispEdges">`)

	writeBestViewedFrame(&builder, totalWidth)

	// Text panel background
	builder.WriteString(`<rect x="`)
	builder.WriteString(number(redPanelWidth))
	builder.WriteString(`" y="3" width="`)
	builder.WriteString(number(float64(textPanelWidth)))
	builder.WriteString(`" height="24" fill="#bcbcbc"/>`)

	textBoxX := float64(redPanelWidth + 1)
	textBoxWidth := float64(textPanelWidth - 2)
	if label == "" {
		writeAdaptiveText(&builder, message, adaptiveTextBox{X: textBoxX, Y: 5, Width: textBoxWidth, Height: 20, MaxFontSize: 15}, "777777", true, true)
	} else {
		writeAdaptiveText(&builder, label, adaptiveTextBox{X: textBoxX, Y: 4, Width: textBoxWidth, Height: 8, MaxFontSize: 7}, "0d0d0d", true, true)
		writeAdaptiveText(&builder, message, adaptiveTextBox{X: textBoxX, Y: 13, Width: textBoxWidth, Height: 13, MaxFontSize: 11}, "777777", true, true)
	}

	writeChromeLogo(&builder, totalWidth-logoWidth)

	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}

func writeBestViewedFrame(builder *strings.Builder, width int) {
	// Base fill
	builder.WriteString(`<rect width="`)
	builder.WriteString(number(float64(width)))
	builder.WriteString(`" height="31" fill="#bcbcbc"/>`)

	// White top and left border
	builder.WriteString(`<path fill="#ffffff" d="M0 0h`)
	builder.WriteString(number(float64(width - 2)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width - 2)))
	builder.WriteString(`zM0 1h1v28h-1z"/>`)

	// Black right and bottom border
	builder.WriteString(`<path fill="#000000" d="M`)
	builder.WriteString(number(float64(width - 1)))
	builder.WriteString(` 1h1v28h-1zM0 30h`)
	builder.WriteString(number(float64(width)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width)))
	builder.WriteString(`z"/>`)

	// Inner shadow line
	builder.WriteString(`<path fill="#a9a9a9" d="M10 28h`)
	builder.WriteString(number(float64(width - 12)))
	builder.WriteString(`v1h-`)
	builder.WriteString(number(float64(width - 12)))
	builder.WriteString(`z"/>`)

	// Red accent panel (left, matches reference stripe pattern)
	builder.WriteString(`<path fill="#d55546" d="M2 2h6v1h-6zM4 3h4v2h-4zM4 7h4v2h-4zM4 15h3v1h-3zM2 27h6v1h-6z"/>`)
	builder.WriteString(`<path fill="#da6a5d" d="M2 5h1v2h-1zM2 11h1v2h-1zM2 17h1v2h-1zM2 23h1v3h-1zM4 25h1v1h-1z"/>`)
	builder.WriteString(`<path fill="#fdf6f5" d="M3 5h1v2h-1zM3 11h1v2h-1zM3 17h1v2h-1zM3 23h1v3h-1z"/>`)
	builder.WriteString(`<path fill="#d75c4d" d="M4 11h1v2h-1zM4 17h1v2h-1zM4 23h1v2h-1z"/>`)
	builder.WriteString(`<path fill="#f9e8e6" d="M5 11h1v2h-1zM5 17h1v2h-1zM5 23h1v2h-1z"/>`)
	builder.WriteString(`<path fill="#d75c4e" d="M6 11h1v2h-1zM6 17h1v2h-1zM6 23h1v2h-1z"/>`)
	builder.WriteString(`<path fill="#fdf6f6" d="M7 11h1v2h-1zM7 17h1v2h-1zM7 23h1v3h-1z"/>`)
	builder.WriteString(`<path fill="#d96659" d="M2 7h1v1h-1zM2 13h1v1h-1zM2 19h1v1h-1zM2 26h1v1h-1z"/>`)
}

// writeChromeLogo draws the Chrome mark with its left edge at x.
// Reference paths are authored with the logo starting at x=60; every
// coordinate is shifted by the delta so the artwork stays on the pixel grid.
func writeChromeLogo(builder *strings.Builder, x int) {
	const referenceLogoX = 60
	shift := float64(x - referenceLogoX)

	writeShiftedPath(builder, "c46545", "M69 5h7v1h-7zM67 6h11v1h-11zM66 7h14v1h-14zM65 8h16v1h-16zM64 9h5v1h-5zM76 9h5v1h-5zM64 10h4v1h-4zM79 10h3v1h-3zM64 11h3v1h-3zM81 11h1v1h-1zM65 12h1v2h-1z", shift)
	writeShiftedPath(builder, "98c87f", "M62 12h2v3h-2zM62 15h3v1h-3zM62 16h4v3h-4zM63 19h4v1h-4zM63 20h5v1h-5zM64 21h5v1h-5zM65 22h10v1h-10zM66 23h8v1h-8zM67 24h6v1h-6zM69 25h3v1h-3z", shift)
	writeShiftedPath(builder, "64a8d2", "M70 11h3v1h-3zM70 12h5v1h-5zM69 13h7v1h-7zM68 14h8v3h-8zM69 17h7v1h-7zM70 18h5v1h-5zM70 19h3v1h-3z", shift)
	writeShiftedPath(builder, "f1d668", "M79 12h1v1h-1zM79 13h3v1h-3zM79 14h4v5h-4zM78 19h4v1h-4zM76 20h6v1h-6zM76 21h5v2h-5zM76 23h4v1h-4zM75 24h3v1h-3zM73 25h3v1h-3z", shift)
	writeShiftedPath(builder, "a5a5a5", "M69 3h7v1h-7z", shift)
	writeShiftedPath(builder, "9a9a9a", "M67 4h1v1h-1zM76 4h2v1h-2z", shift)
	writeShiftedPath(builder, "765246", "M67 5h1v1h-1zM76 5h2v1h-2z", shift)
	writeShiftedPath(builder, "654e46", "M78 6h1v1h-1zM81 9h1v1h-1z", shift)
	writeShiftedPath(builder, "676767", "M69 27h7v1h-7z", shift)
	writeShiftedPath(builder, "5c5c5c", "M67 26h1v1h-1zM76 26h2v1h-2z", shift)
	writeShiftedPath(builder, "656565", "M63 23h1v1h-1zM65 25h1v1h-1zM66 26h1v1h-1z", shift)
	writeShiftedPath(builder, "767676", "M68 27h1v1h-1z", shift)
	writeShiftedPath(builder, "ffffff", "M67 18h1v1h-1z", shift)
}

// writeShiftedPath re-emits a reference path with every horizontal
// coordinate translated by shift, keeping vertical geometry untouched.
func writeShiftedPath(builder *strings.Builder, color, path string, shift float64) {
	builder.WriteString(`<path fill="#`)
	builder.WriteString(color)
	builder.WriteString(`" d="`)
	for _, segment := range strings.Split(path, "M") {
		if segment == "" {
			continue
		}
		space := strings.IndexByte(segment, ' ')
		if space < 0 {
			continue
		}
		startX, err := strconv.ParseFloat(segment[:space], 64)
		if err != nil {
			continue
		}
		builder.WriteByte('M')
		builder.WriteString(number(startX + shift))
		builder.WriteString(segment[space:])
	}
	builder.WriteString(`"/>`)
}
