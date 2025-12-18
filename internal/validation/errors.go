package validation

const (
	NotBlankMessage = "Must not be blank."
)

type Result struct {
	GeneralErrors []string
	FieldErrors   map[string][]string
}
