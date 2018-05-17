package messagebird

type Balance struct {
	Payment string
	Type    string
	Amount  float32
	Errors  []Error
}
