package types

type Entry interface {
	Name() string
	Bands() []Band
	String() string
}
