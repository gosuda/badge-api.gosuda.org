package main

import (
	"badge-api.gosuda.org/internal/styles"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type badgeOptions struct {
	Label            string
	Message          string
	Style            string
	LabelColor       string
	MessageColor     string
	LabelTextColor   string
	MessageTextColor string
	Size             int
	LetterSpacing    float64
}

const (
	defaultBadgeSize = 100
	minBadgeSize     = 50
	maxBadgeSize     = 300
	minLetterSpacing = -1.0
	maxLetterSpacing = 3.0
)

var namedColors = map[string]string{
	"brightgreen":   "44cc11",
	"green":         "97ca00",
	"yellowgreen":   "a4a61d",
	"yellow":        "dfb317",
	"orange":        "fe7d37",
	"red":           "e05d44",
	"blue":          "007ec6",
	"grey":          "555555",
	"gray":          "555555",
	"lightgrey":     "9f9f9f",
	"lightgray":     "9f9f9f",
	"success":       "2f855a",
	"important":     "d97706",
	"critical":      "c53030",
	"informational": "2563eb",
	"inactive":      "718096",
}

func availableStyles() []string {
	return styles.Available()
}

func normalizeColor(value, fallback string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "#")
	if named, ok := namedColors[value]; ok {
		return named, nil
	}
	if value == "" {
		return fallback, nil
	}
	if len(value) == 3 {
		value = string([]byte{value[0], value[0], value[1], value[1], value[2], value[2]})
	}
	if len(value) != 6 {
		return "", fmt.Errorf("color must be a 3- or 6-digit hex value")
	}
	for _, char := range value {
		if !strings.ContainsRune("0123456789abcdef", char) {
			return "", fmt.Errorf("color contains a non-hex character")
		}
	}
	return value, nil
}

func validateBadgeOptions(options badgeOptions) (badgeOptions, error) {
	options.Label = strings.TrimSpace(options.Label)
	options.Message = strings.TrimSpace(options.Message)
	options.Style = strings.ToLower(strings.TrimSpace(options.Style))
	if options.Style == "" {
		options.Style = "flat"
	}
	if !styles.Exists(options.Style) {
		return badgeOptions{}, fmt.Errorf("unknown style %q; supported styles: %s", options.Style, strings.Join(availableStyles(), ", "))
	}
	if options.Size == 0 {
		options.Size = defaultBadgeSize
	}
	if options.Size < minBadgeSize || options.Size > maxBadgeSize {
		return badgeOptions{}, fmt.Errorf("size must be between %d and %d percent", minBadgeSize, maxBadgeSize)
	}
	if math.IsNaN(options.LetterSpacing) || math.IsInf(options.LetterSpacing, 0) ||
		options.LetterSpacing < minLetterSpacing || options.LetterSpacing > maxLetterSpacing {
		return badgeOptions{}, fmt.Errorf("letterSpacing must be between %g and %g pixels", minLetterSpacing, maxLetterSpacing)
	}
	if options.Message == "" {
		return badgeOptions{}, fmt.Errorf("message is required")
	}
	if utf8.RuneCountInString(options.Label) > 64 {
		return badgeOptions{}, fmt.Errorf("label must be 64 characters or fewer")
	}
	if utf8.RuneCountInString(options.Message) > 128 {
		return badgeOptions{}, fmt.Errorf("message must be 128 characters or fewer")
	}

	var err error
	if options.LabelColor, err = normalizeColor(options.LabelColor, "555555"); err != nil {
		return badgeOptions{}, fmt.Errorf("labelColor: %w", err)
	}
	if options.MessageColor, err = normalizeColor(options.MessageColor, "44cc11"); err != nil {
		return badgeOptions{}, fmt.Errorf("color: %w", err)
	}
	if options.LabelTextColor, err = normalizeColor(options.LabelTextColor, "ffffff"); err != nil {
		return badgeOptions{}, fmt.Errorf("labelTextColor: %w", err)
	}
	if options.MessageTextColor, err = normalizeColor(options.MessageTextColor, "ffffff"); err != nil {
		return badgeOptions{}, fmt.Errorf("textColor: %w", err)
	}
	return options, nil
}

func renderBadge(options badgeOptions) []byte {
	return styles.Render(styles.Options{
		Label:            options.Label,
		Message:          options.Message,
		Style:            options.Style,
		LabelColor:       options.LabelColor,
		MessageColor:     options.MessageColor,
		LabelTextColor:   options.LabelTextColor,
		MessageTextColor: options.MessageTextColor,
		Size:             options.Size,
		LetterSpacing:    options.LetterSpacing,
	})
}
