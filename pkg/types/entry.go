package types

// Entry is the common interface implemented by downlink and uplink combo records.
type Entry interface {
	Name() string
	Bands() []Band
	String() string
}
