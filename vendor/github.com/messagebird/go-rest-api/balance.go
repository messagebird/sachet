package messagebird

// Balance describes your balance information.
type Balance struct {
	Payment string
	Type    string
	Amount  float32
	Errors  []Error
}
