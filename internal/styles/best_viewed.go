package styles

import (
	_ "embed"
	"strings"
)

//go:embed assets/best-viewed-with-chrome-88x31.svg
var bestViewedReferenceSVG string

var bestViewedReferenceArtwork = referenceArtwork(bestViewedReferenceSVG)
var bestViewedStyle = styleSpec{Height: classicHeight, Kind: "best-viewed", Render: renderBestViewed}

func renderBestViewed(options Options, label, message string) []byte {
	var builder strings.Builder
	builder.Grow(len(bestViewedReferenceArtwork) + 1800)
	writeFixedSVGStart(&builder, classicWidth, classicHeight, options.Size, accessibleLabel(label, message))
	builder.WriteString(bestViewedReferenceArtwork)

	if !strings.EqualFold(strings.TrimSpace(label), "viewed with") || !strings.EqualFold(strings.TrimSpace(message), "chrome") {
		builder.WriteString(`<rect x="10" y="3" width="50" height="24" fill="#bcbcbc"/>`)
		if label == "" {
			writeAdaptiveText(&builder, message, adaptiveTextBox{X: 11, Y: 5, Width: 47, Height: 20, MaxFontSize: 15}, "777777", true, true)
		} else {
			writeAdaptiveText(&builder, label, adaptiveTextBox{X: 11, Y: 4, Width: 47, Height: 8, MaxFontSize: 7}, "0d0d0d", true, true)
			writeAdaptiveText(&builder, message, adaptiveTextBox{X: 11, Y: 13, Width: 47, Height: 13, MaxFontSize: 11}, "777777", true, true)
		}
	}
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}
