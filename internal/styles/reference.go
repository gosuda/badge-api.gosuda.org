package styles

import "strings"

func referenceArtwork(source string) string {
	start := strings.Index(source, "<rect ")
	end := strings.LastIndex(source, "</svg>")
	if start < 0 || end <= start {
		panic("invalid embedded SVG reference")
	}
	return source[start:end]
}
