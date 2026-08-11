package styles

type Options struct {
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

type styleSpec struct {
	Height        float64
	Radius        float64
	FontSize      float64
	LetterSpacing float64
	Padding       float64
	Uppercase     bool
	BoldLabel     bool
	BoldMessage   bool
	Kind          string
	Render        func(Options, string, string) []byte
}

type namedStyle struct {
	Name string
	Spec styleSpec
}
