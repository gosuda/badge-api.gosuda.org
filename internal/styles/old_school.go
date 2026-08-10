package styles

import (
	"math"
	"strings"
)

const (
	oldSchoolWidth       = 80
	oldSchoolHeight      = 15
	oldSchoolPanelHeight = 11
)

var oldSchoolStyle = styleSpec{Height: oldSchoolHeight, Uppercase: true, Kind: "old-school", Render: renderOldSchool}

func renderOldSchool(options Options, label, message string) []byte {
	labelPanelWidth, messagePanelWidth := oldSchoolPanelWidths(label, message)
	var builder strings.Builder
	builder.Grow(2200)
	writeFixedSVGStart(&builder, oldSchoolWidth, oldSchoolHeight, options.Size, accessibleLabel(label, message))
	builder.WriteString(`<rect width="80" height="15" fill="#a5a5a5"/><rect x="1" y="1" width="78" height="13" fill="#fff"/>`)

	messageX := 2
	if labelPanelWidth > 0 {
		writeOldSchoolPanel(&builder, 2, labelPanelWidth, options.LabelColor)
		messageX += labelPanelWidth + 1
	}
	writeOldSchoolPanel(&builder, messageX, messagePanelWidth, options.MessageColor)
	if labelPanelWidth > 0 {
		writeAdaptiveText(&builder, label, adaptiveTextBox{X: 3, Y: 2, Width: float64(labelPanelWidth - 2), Height: 11, MaxFontSize: 7}, options.LabelTextColor, false, true)
	}
	writeAdaptiveText(&builder, message, adaptiveTextBox{X: float64(messageX + 1), Y: 2, Width: float64(messagePanelWidth - 2), Height: 11, MaxFontSize: 7}, options.MessageTextColor, false, true)
	builder.WriteString(`</g></svg>`)
	return []byte(builder.String())
}

func oldSchoolPanelWidths(label, message string) (int, int) {
	if label == "" {
		return 0, 76
	}
	const usableWidth = 75
	labelNeeded := int(math.Ceil(measureText(label, 7))) + 4
	messageNeeded := int(math.Ceil(measureText(message, 7))) + 4
	totalNeeded := labelNeeded + messageNeeded
	labelWidth := 0
	if totalNeeded <= usableWidth {
		labelWidth = labelNeeded + (usableWidth-totalNeeded)/3
	} else {
		labelWidth = int(math.Round(float64(usableWidth) * float64(labelNeeded) / float64(totalNeeded)))
	}
	labelWidth = max(9, min(labelWidth, usableWidth-9))
	return labelWidth, usableWidth - labelWidth
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
