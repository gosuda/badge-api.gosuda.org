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
		glyphs := fitBitmapText([]rune(strings.ToUpper(phrase)), 69, 2, 1)
		if len(glyphs) > 0 {
			textWidth := bitmapTextWidth(glyphs, 2, 1)
			writeBitmapPath(&builder, glyphs, 5+(72-textWidth)/2, 10, 2, 2, 1, "3a0603")
		}
	}
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}
