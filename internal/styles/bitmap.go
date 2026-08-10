package styles

import "strings"

var bitmapGlyphs = map[rune][5]uint8{
	' ': {},
	'A': {0b010, 0b101, 0b111, 0b101, 0b101}, 'B': {0b110, 0b101, 0b110, 0b101, 0b110},
	'C': {0b011, 0b100, 0b100, 0b100, 0b011}, 'D': {0b110, 0b101, 0b101, 0b101, 0b110},
	'E': {0b111, 0b100, 0b110, 0b100, 0b111}, 'F': {0b111, 0b100, 0b110, 0b100, 0b100},
	'G': {0b011, 0b100, 0b101, 0b101, 0b011}, 'H': {0b101, 0b101, 0b111, 0b101, 0b101},
	'I': {0b111, 0b010, 0b010, 0b010, 0b111}, 'J': {0b001, 0b001, 0b001, 0b101, 0b010},
	'K': {0b101, 0b101, 0b110, 0b101, 0b101}, 'L': {0b100, 0b100, 0b100, 0b100, 0b111},
	'M': {0b101, 0b111, 0b111, 0b101, 0b101}, 'N': {0b101, 0b111, 0b111, 0b111, 0b101},
	'O': {0b010, 0b101, 0b101, 0b101, 0b010}, 'P': {0b110, 0b101, 0b110, 0b100, 0b100},
	'Q': {0b010, 0b101, 0b101, 0b111, 0b011}, 'R': {0b110, 0b101, 0b110, 0b101, 0b101},
	'S': {0b011, 0b100, 0b010, 0b001, 0b110}, 'T': {0b111, 0b010, 0b010, 0b010, 0b010},
	'U': {0b101, 0b101, 0b101, 0b101, 0b111}, 'V': {0b101, 0b101, 0b101, 0b101, 0b010},
	'W': {0b101, 0b101, 0b111, 0b111, 0b101}, 'X': {0b101, 0b101, 0b010, 0b101, 0b101},
	'Y': {0b101, 0b101, 0b010, 0b010, 0b010}, 'Z': {0b111, 0b001, 0b010, 0b100, 0b111},
	'0': {0b111, 0b101, 0b101, 0b101, 0b111}, '1': {0b010, 0b110, 0b010, 0b010, 0b111},
	'2': {0b110, 0b001, 0b010, 0b100, 0b111}, '3': {0b110, 0b001, 0b010, 0b001, 0b110},
	'4': {0b101, 0b101, 0b111, 0b001, 0b001}, '5': {0b111, 0b100, 0b110, 0b001, 0b110},
	'6': {0b011, 0b100, 0b111, 0b101, 0b111}, '7': {0b111, 0b001, 0b010, 0b010, 0b010},
	'8': {0b111, 0b101, 0b111, 0b101, 0b111}, '9': {0b111, 0b101, 0b111, 0b001, 0b110},
	'!': {0b010, 0b010, 0b010, 0b000, 0b010}, '?': {0b110, 0b001, 0b010, 0b000, 0b010},
	'.': {0b000, 0b000, 0b000, 0b000, 0b010}, ',': {0b000, 0b000, 0b000, 0b010, 0b100},
	':': {0b000, 0b010, 0b000, 0b010, 0b000}, ';': {0b000, 0b010, 0b000, 0b010, 0b100},
	'-': {0b000, 0b000, 0b111, 0b000, 0b000}, '_': {0b000, 0b000, 0b000, 0b000, 0b111},
	'+': {0b000, 0b010, 0b111, 0b010, 0b000}, '=': {0b000, 0b111, 0b000, 0b111, 0b000},
	'[': {0b110, 0b100, 0b100, 0b100, 0b110}, ']': {0b011, 0b001, 0b001, 0b001, 0b011},
	'(': {0b010, 0b100, 0b100, 0b100, 0b010}, ')': {0b010, 0b001, 0b001, 0b001, 0b010},
	'<': {0b001, 0b010, 0b100, 0b010, 0b001}, '>': {0b100, 0b010, 0b001, 0b010, 0b100},
	'/': {0b001, 0b001, 0b010, 0b100, 0b100}, '\\': {0b100, 0b100, 0b010, 0b001, 0b001},
	'|': {0b010, 0b010, 0b010, 0b010, 0b010}, '&': {0b010, 0b101, 0b010, 0b101, 0b011},
	'#': {0b101, 0b111, 0b101, 0b111, 0b101}, '%': {0b101, 0b001, 0b010, 0b100, 0b101},
	'@': {0b111, 0b101, 0b111, 0b100, 0b011}, '…': {0b000, 0b000, 0b000, 0b000, 0b111},
}

func bitmapTextWidth(text []rune, scaleX, spacing int) int {
	if len(text) == 0 {
		return 0
	}
	return len(text)*(3*scaleX+spacing) - spacing
}

func fitBitmapText(text []rune, maxWidth, scaleX, spacing int) []rune {
	step := 3*scaleX + spacing
	maxGlyphs := (maxWidth + spacing) / step
	if maxGlyphs <= 0 {
		return nil
	}
	if len(text) <= maxGlyphs {
		return text
	}
	if maxGlyphs == 1 {
		return []rune{'…'}
	}
	fitted := make([]rune, maxGlyphs)
	copy(fitted, text[:maxGlyphs-1])
	fitted[maxGlyphs-1] = '…'
	return fitted
}

func writeBitmapText(builder *strings.Builder, text string, areaX, y, areaWidth, scaleX, scaleY, spacing int, color string) {
	glyphs := fitBitmapText([]rune(strings.ToUpper(text)), areaWidth, scaleX, spacing)
	if len(glyphs) == 0 {
		return
	}
	width := bitmapTextWidth(glyphs, scaleX, spacing)
	writeBitmapPath(builder, glyphs, areaX+(areaWidth-width)/2, y, scaleX, scaleY, spacing, color)
}

func writeBitmapPath(builder *strings.Builder, glyphs []rune, x, y, scaleX, scaleY, spacing int, color string) {
	builder.WriteString(`<path fill="#`)
	builder.WriteString(color)
	builder.WriteString(`" d="`)
	for _, char := range glyphs {
		glyph, ok := bitmapGlyphs[char]
		if !ok {
			glyph = bitmapGlyphs['?']
		}
		writeBitmapGlyphPath(builder, glyph, x, y, scaleX, scaleY)
		x += 3*scaleX + spacing
	}
	builder.WriteString(`"/>`)
}

func writeBitmapGlyphPath(builder *strings.Builder, glyph [5]uint8, x, y, scaleX, scaleY int) {
	for row, bits := range glyph {
		runStart := -1
		for column := 0; column <= 3; column++ {
			filled := column < 3 && bits&(1<<uint(2-column)) != 0
			if filled && runStart < 0 {
				runStart = column
			}
			if !filled && runStart >= 0 {
				runWidth := (column - runStart) * scaleX
				builder.WriteByte('M')
				builder.WriteString(number(float64(x + runStart*scaleX)))
				builder.WriteByte(' ')
				builder.WriteString(number(float64(y + row*scaleY)))
				builder.WriteByte('h')
				builder.WriteString(number(float64(runWidth)))
				builder.WriteByte('v')
				builder.WriteString(number(float64(scaleY)))
				builder.WriteString(`h-`)
				builder.WriteString(number(float64(runWidth)))
				builder.WriteByte('z')
				runStart = -1
			}
		}
	}
}
