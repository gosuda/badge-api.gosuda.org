package styles

import (
	_ "embed"
	"strings"
)

const (
	classicWidth  = 88
	classicHeight = 31
)

//go:embed assets/click-here-88x31.svg
var clickHereReferenceSVG string

var clickHereReferenceArtwork = referenceArtwork(clickHereReferenceSVG)
var clickHereStyle = styleSpec{Height: classicHeight, Kind: "click-here", Render: renderClickHere}

func renderClickHere(options Options, label, message string) []byte {
	var builder strings.Builder
	builder.Grow(len(clickHereReferenceArtwork) + 1800)
	writeFixedSVGStart(&builder, classicWidth, classicHeight, options.Size, accessibleLabel(label, message))
	builder.WriteString(clickHereReferenceArtwork)

	phrase := strings.TrimSpace(strings.TrimSpace(label) + " " + strings.TrimSpace(message))
	if !strings.EqualFold(phrase, "click here") {
		builder.WriteString(`<rect x="3" y="3" width="76" height="24" fill="#cccccc"/>`)
		writeAdaptiveText(&builder, phrase, adaptiveTextBox{X: 5, Y: 5, Width: 72, Height: 20, MaxFontSize: 11}, "3a0603", true, true)
	}
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}
