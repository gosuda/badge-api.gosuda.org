package styles

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestMeasureTextUsesUnicodeDisplayWidth(t *testing.T) {
	plain := measureText("e", 11, 0)
	if got := measureText("e\u0301", 11, 0); got != plain {
		t.Fatalf("combining text width = %v, want %v", got, plain)
	}
	wide := measureText("界", 11, 0)
	if wide <= measureText("a", 11, 0) {
		t.Fatalf("CJK width = %v, want wider than ASCII", wide)
	}
	if emoji := measureText("👩‍💻", 11, 0); emoji != wide {
		t.Fatalf("emoji grapheme width = %v, want %v", emoji, wide)
	}
}

func TestMeasureTextAppliesSpacingBetweenGraphemeClusters(t *testing.T) {
	const spacing = 0.5
	base := measureText("e\u0301x", 11, 0)
	spaced := measureText("e\u0301x", 11, spacing)
	if got := spaced - base; math.Abs(got-spacing) > 1e-9 {
		t.Fatalf("spacing width adjustment = %v, want %v", got, spacing)
	}
}
func TestEveryStyleAppliesLetterSpacing(t *testing.T) {
	for _, style := range Available() {
		t.Run(style, func(t *testing.T) {
			svg := string(Render(Options{
				Label:         "build",
				Message:       "ready",
				Style:         style,
				Size:          100,
				LetterSpacing: 0.5,
			}))
			if !strings.Contains(svg, `letter-spacing="`) {
				t.Fatal("SVG does not apply requested letter spacing")
			}
		})
	}
}

func TestFlatbarUsesRelaxedDefaultTracking(t *testing.T) {
	svg := string(Render(Options{Label: "mood", Message: "loud", Style: "flatbar", Size: 100}))
	if !strings.Contains(svg, `letter-spacing="0.45"`) {
		t.Fatal("flatbar SVG does not apply its relaxed default tracking")
	}
}

func TestReferenceExamplesPreserveExactArtwork(t *testing.T) {
	tests := []struct {
		name      string
		reference Options
		artwork   string
		custom    Options
	}{
		{
			name:      "click-here",
			reference: Options{Label: "click", Message: "here", Style: "click-here", Size: 100},
			artwork:   clickHereReferenceArtwork,
			custom:    Options{Label: "ship", Message: "now", Style: "click-here", Size: 100},
		},
		{
			name:      "best-viewed",
			reference: Options{Label: "viewed with", Message: "chrome", Style: "best-viewed", Size: 100},
			artwork:   bestViewedReferenceArtwork,
			custom:    Options{Label: "built with", Message: "Go", Style: "best-viewed", Size: 100},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reference case should embed exact artwork
			refSVG := string(Render(test.reference))
			if !strings.Contains(refSVG, test.artwork) {
				t.Fatal("reference artwork coordinates were not preserved verbatim")
			}
			spacedReference := test.reference
			spacedReference.LetterSpacing = 0.5
			spacedSVG := string(Render(spacedReference))
			if strings.Contains(spacedSVG, test.artwork) || !strings.Contains(spacedSVG, `letter-spacing="0.5"`) {
				t.Fatal("letter spacing should replace fixed artwork with adjustable vector text")
			}

			// Custom text generates procedural frames and adaptive text
			customSVG := string(Render(test.custom))
			if !strings.Contains(customSVG, `<text `) {
				t.Fatal("custom renderer does not use adaptive vector text")
			}
			if !strings.Contains(customSVG, `text-rendering="geometricPrecision"`) {
				t.Fatal("custom renderer missing geometricPrecision")
			}
			// Variable-width designs generate frames procedurally, not by embedding reference artwork
			if strings.Contains(customSVG, test.artwork) {
				t.Fatal("custom renderer should generate procedural frames, not embed reference artwork")
			}
		})
	}
}

func TestAdaptiveFontSizeShrinksForLongText(t *testing.T) {
	box := adaptiveTextBox{Width: 47, Height: 20, MaxFontSize: 11}
	short := adaptiveFontSize("GO", box, 0)
	long := adaptiveFontSize("A VERY LONG STATUS MESSAGE", box, 0)
	if short != box.MaxFontSize {
		t.Fatalf("short font size = %v, want maximum %v", short, box.MaxFontSize)
	}
	if long >= short {
		t.Fatalf("long font size = %v, want smaller than %v", long, short)
	}
	if long < 5.0 {
		t.Fatalf("long font size = %v, want at least 5.0", long)
	}
}

func TestAdaptiveCustomTextPreservesFullContent(t *testing.T) {
	message := "A very long status message"
	svg := string(Render(Options{Label: "release channel", Message: message, Style: "best-viewed", Size: 100}))
	for _, expected := range []string{"RELEASE CHANNEL", "A VERY LONG STATUS MESSAGE", `lengthAdjust="spacingAndGlyphs"`} {
		if !strings.Contains(svg, expected) {
			t.Fatalf("adaptive SVG does not contain %q", expected)
		}
	}
}

// TestRetroVariableWidthFormula locks the three width-calculation contracts:
// measurement uses the uppercased phrase at the style's target font size,
// widths stay within [min, max], and panels fill the available canvas exactly.
func TestRetroVariableWidthFormula(t *testing.T) {
	t.Run("old-school panels fill canvas", func(t *testing.T) {
		for _, c := range []struct{ label, msg string }{
			{"build", "passing"},
			{"", "ok"},
			{"a very long label with many words", "and an equally long message right here"},
		} {
			lp, mp, total := calculateOldSchoolWidths(c.label, c.msg, 0)
			want := total - 5 // borderAndGaps: 2px border + 1px gap + 2px border
			if c.label == "" {
				want = total - 4 // no gap when there is no label panel
			}
			if got := lp + mp; got != want {
				t.Errorf("label=%q msg=%q: panels %d+%d=%d, want %d (total=%d)",
					c.label, c.msg, lp, mp, got, want, total)
			}
			if total < oldSchoolMinWidth || total > oldSchoolMaxWidth {
				t.Errorf("label=%q msg=%q: total %d out of [%d, %d]",
					c.label, c.msg, total, oldSchoolMinWidth, oldSchoolMaxWidth)
			}
		}
	})

	t.Run("click-here measures uppercase at 7px", func(t *testing.T) {
		phrase := "ship now"
		upper := strings.ToUpper(phrase)
		needed := int(math.Ceil(measureText(upper, 7, 0))) + 12 + 10
		want := max(clickHereMinWidth, min(needed, clickHereMaxWidth))
		svg := string(Render(Options{Label: "", Message: phrase, Style: "click-here", Size: 100}))
		if got := viewBoxWidth(t, svg); got != float64(want) {
			t.Errorf("click-here width = %v, want %d for %q at 7px", got, want, upper)
		}
	})

	t.Run("retro widths clamp to max 300", func(t *testing.T) {
		long := strings.Repeat("w", 200)
		for _, style := range []string{"old-school", "click-here", "best-viewed"} {
			svg := string(Render(Options{Label: long, Message: long, Style: style, Size: 100}))
			if got := viewBoxWidth(t, svg); got != 300 {
				t.Errorf("%s: width = %v, want exactly 300 for extreme input", style, got)
			}
		}
	})
}

// viewBoxWidth returns the width component of the SVG viewBox, failing the
// test when the attribute is absent or unparseable.
func viewBoxWidth(t *testing.T, svg string) float64 {
	t.Helper()
	const prefix = `viewBox="0 0 `
	i := strings.Index(svg, prefix)
	if i < 0 {
		t.Fatalf("no viewBox in SVG")
	}
	rest := svg[i+len(prefix):]
	j := strings.IndexByte(rest, ' ')
	if j < 0 {
		t.Fatalf("malformed viewBox in SVG")
	}
	width, err := strconv.ParseFloat(rest[:j], 64)
	if err != nil {
		t.Fatalf("unparseable viewBox width %q: %v", rest[:j], err)
	}
	return width
}
