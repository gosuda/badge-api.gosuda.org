package styles

import "strings"

var styleCatalog = []namedStyle{
	{Name: "flat", Spec: flatStyle},
	{Name: "flat-square", Spec: flatSquareStyle},
	{Name: "plastic", Spec: plasticStyle},
	{Name: "round", Spec: roundStyle},
	{Name: "outline", Spec: outlineStyle},
	{Name: "neon", Spec: neonStyle},
	{Name: "glass", Spec: glassStyle},
	{Name: "flatbar", Spec: flatbarStyle},
	{Name: "old-school", Spec: oldSchoolStyle},
	{Name: "click-here", Spec: clickHereStyle},
	{Name: "best-viewed", Spec: bestViewedStyle},
}

var stylesByName = func() map[string]styleSpec {
	styles := make(map[string]styleSpec, len(styleCatalog))
	for _, style := range styleCatalog {
		styles[style.Name] = style.Spec
	}
	return styles
}()

func Available() []string {
	styles := make([]string, len(styleCatalog))
	for index, style := range styleCatalog {
		styles[index] = style.Name
	}
	return styles
}

func Exists(name string) bool {
	_, ok := stylesByName[name]
	return ok
}

func Render(options Options) []byte {
	spec, ok := stylesByName[options.Style]
	if !ok {
		return nil
	}
	label := options.Label
	message := options.Message
	if spec.Uppercase {
		label = strings.ToUpper(label)
		message = strings.ToUpper(message)
	}
	if spec.Render != nil {
		return spec.Render(options, label, message)
	}
	return renderStandard(options, label, message, spec)
}
