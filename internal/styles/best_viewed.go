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
			writeBitmapText(&builder, message, 11, 10, 47, 2, 2, 1, "777777")
		} else {
			writeBitmapText(&builder, label, 11, 6, 47, 1, 1, 1, "0d0d0d")
			writeBitmapText(&builder, message, 11, 15, 47, 2, 2, 1, "777777")
		}
	}
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}
