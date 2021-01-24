package translate

type Intf interface {
	Translate(from, to, q string) string
}
