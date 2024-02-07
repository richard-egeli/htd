package button

type ButtonType string

func (ButtonType) Submit() ButtonType { return ButtonType("submit") }
func (ButtonType) Button() ButtonType { return ButtonType("button") }
func (ButtonType) Reset() ButtonType  { return ButtonType("reset") }
func (ButtonType) Menu() ButtonType   { return ButtonType("menu") }

type ButtonData struct {
	Type ButtonType
	Text string
}
