package styles

import (
	"strings"
	"testing"
)

func TestMeasureTextUsesUnicodeDisplayWidth(t *testing.T) {
	plain := measureText("e", 11)
	if got := measureText("e\u0301", 11); got != plain {
		t.Fatalf("combining text width = %v, want %v", got, plain)
	}
	wide := measureText("界", 11)
	if wide <= measureText("a", 11) {
		t.Fatalf("CJK width = %v, want wider than ASCII", wide)
	}
	if emoji := measureText("👩‍💻", 11); emoji != wide {
		t.Fatalf("emoji grapheme width = %v, want %v", emoji, wide)
	}
}

func TestReferenceExamplesPreserveExactArtwork(t *testing.T) {
	tests := []struct {
		name      string
		options   Options
		artwork   string
		overlay   string
		customize Options
	}{
		{
			name:      "click-here",
			options:   Options{Label: "click", Message: "here", Style: "click-here", Size: 100},
			artwork:   clickHereReferenceArtwork,
			overlay:   `<rect x="3" y="3" width="76" height="24" fill="#cccccc"/>`,
			customize: Options{Label: "ship", Message: "now", Style: "click-here", Size: 100},
		},
		{
			name:      "best-viewed",
			options:   Options{Label: "viewed with", Message: "chrome", Style: "best-viewed", Size: 100},
			artwork:   bestViewedReferenceArtwork,
			overlay:   `<rect x="10" y="3" width="50" height="24" fill="#bcbcbc"/>`,
			customize: Options{Label: "built with", Message: "Go", Style: "best-viewed", Size: 100},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reference := string(Render(test.options))
			if !strings.Contains(reference, test.artwork) {
				t.Fatal("reference artwork coordinates were not preserved verbatim")
			}
			if strings.Contains(reference, test.overlay) {
				t.Fatal("reference example unexpectedly paints over original artwork")
			}
			custom := string(Render(test.customize))
			if !strings.Contains(custom, test.artwork) || !strings.Contains(custom, test.overlay) {
				t.Fatal("custom renderer does not preserve artwork before replacing the text region")
			}
			if !strings.Contains(custom, `<text `) || !strings.Contains(custom, `text-rendering="geometricPrecision"`) {
				t.Fatal("custom renderer does not use adaptive vector text")
			}
		})
	}
}

func TestAdaptiveFontSizeShrinksForLongText(t *testing.T) {
	box := adaptiveTextBox{Width: 47, Height: 20, MaxFontSize: 11}
	short := adaptiveFontSize("GO", box)
	long := adaptiveFontSize("A VERY LONG STATUS MESSAGE", box)
	if short != box.MaxFontSize {
		t.Fatalf("short font size = %v, want maximum %v", short, box.MaxFontSize)
	}
	if long >= short {
		t.Fatalf("long font size = %v, want smaller than %v", long, short)
	}
	if long < 0.5 {
		t.Fatalf("long font size = %v, want at least 0.5", long)
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
