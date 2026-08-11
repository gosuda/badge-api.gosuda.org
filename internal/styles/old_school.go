package styles

import (
	"html"
	"math"
	"strings"
)

const (
	oldSchoolHeight      = 15
	oldSchoolPanelHeight = 11
	oldSchoolMinWidth    = 80
	oldSchoolMaxWidth    = 300
)

var oldSchoolStyle = styleSpec{Height: oldSchoolHeight, Uppercase: true, Kind: "old-school", Render: renderOldSchool}

func renderOldSchool(options Options, label, message string) []byte {
	labelPanelWidth, messagePanelWidth, totalWidth := calculateOldSchoolWidths(label, message, options.LetterSpacing)
	var builder strings.Builder
	builder.Grow(2200)

	aria := accessibleLabel(label, message)
	scale := float64(options.Size) / 100
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="`)
	builder.WriteString(number(float64(totalWidth) * scale))
	builder.WriteString(`" height="`)
	builder.WriteString(number(oldSchoolHeight * scale))
	builder.WriteString(`" viewBox="0 0 `)
	builder.WriteString(number(float64(totalWidth)))
	builder.WriteByte(' ')
	builder.WriteString(number(oldSchoolHeight))
	builder.WriteString(`" role="img" aria-label="`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`"><title>`)
	builder.WriteString(html.EscapeString(aria))
	builder.WriteString(`</title><g shape-rendering="crispEdges">`)

	// Frame and background
	builder.WriteString(`<rect width="`)
	builder.WriteString(number(float64(totalWidth)))
	builder.WriteString(`" height="15" fill="#a5a5a5"/><rect x="1" y="1" width="`)
	builder.WriteString(number(float64(totalWidth - 2)))
	builder.WriteString(`" height="13" fill="#fff"/>`)

	messageX := 2
	if labelPanelWidth > 0 {
		writeOldSchoolPanel(&builder, 2, labelPanelWidth, options.LabelColor)
		messageX += labelPanelWidth + 1
	}
	writeOldSchoolPanel(&builder, messageX, messagePanelWidth, options.MessageColor)

	if labelPanelWidth > 0 {
		writeAdaptiveText(&builder, label, adaptiveTextBox{X: 3, Y: 2, Width: float64(labelPanelWidth - 2), Height: 11, MaxFontSize: 7}, options.LabelTextColor, options.LetterSpacing, false, true)
	}
	writeAdaptiveText(&builder, message, adaptiveTextBox{X: float64(messageX + 1), Y: 2, Width: float64(messagePanelWidth - 2), Height: 11, MaxFontSize: 7}, options.MessageTextColor, options.LetterSpacing, false, true)
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}

func calculateOldSchoolWidths(label, message string, letterSpacing float64) (int, int, int) {
	const (
		borderAndGaps = 5 // 2px left border + 1px gap + 2px right border
		minPanelWidth = 9
	)

	if label == "" {
		messageNeeded := int(math.Ceil(measureText(strings.ToUpper(message), 7, letterSpacing))) + 4
		messagePanelWidth := max(messageNeeded, 20)
		totalWidth := messagePanelWidth + 4 // 2px border each side
		totalWidth = max(oldSchoolMinWidth, min(totalWidth, oldSchoolMaxWidth))
		messagePanelWidth = totalWidth - 4
		return 0, messagePanelWidth, totalWidth
	}

	labelNeeded := int(math.Ceil(measureText(strings.ToUpper(label), 7, letterSpacing))) + 4
	messageNeeded := int(math.Ceil(measureText(strings.ToUpper(message), 7, letterSpacing))) + 4

	// Calculate ideal total width
	idealTotal := labelNeeded + messageNeeded + borderAndGaps
	totalWidth := max(oldSchoolMinWidth, min(idealTotal, oldSchoolMaxWidth))

	// Allocate exactly (totalWidth - borderAndGaps) across the two panels
	availableWidth := totalWidth - borderAndGaps

	var labelPanelWidth, messagePanelWidth int
	if labelNeeded+messageNeeded <= availableWidth {
		// Fits: use measured widths and distribute spare evenly
		spare := availableWidth - (labelNeeded + messageNeeded)
		labelPanelWidth = labelNeeded + spare/2
		messagePanelWidth = availableWidth - labelPanelWidth
	} else {
		// Overflow: proportionally reduce both, respecting 9px minima
		totalNeeded := labelNeeded + messageNeeded
		labelRatio := float64(labelNeeded) / float64(totalNeeded)
		labelPanelWidth = max(minPanelWidth, int(math.Round(float64(availableWidth)*labelRatio)))
		messagePanelWidth = availableWidth - labelPanelWidth

		// Ensure message also gets at least minPanelWidth
		if messagePanelWidth < minPanelWidth {
			messagePanelWidth = minPanelWidth
			labelPanelWidth = availableWidth - messagePanelWidth
		}
	}

	return labelPanelWidth, messagePanelWidth, totalWidth
}

func writeOldSchoolPanel(builder *strings.Builder, x, width int, color string) {
	builder.WriteString(`<rect x="`)
	builder.WriteString(number(float64(x)))
	builder.WriteString(`" y="2" width="`)
	builder.WriteString(number(float64(width)))
	builder.WriteString(`" height="`)
	builder.WriteString(number(oldSchoolPanelHeight))
	builder.WriteString(`" fill="#`)
	builder.WriteString(color)
	builder.WriteString(`"/>`)
}
