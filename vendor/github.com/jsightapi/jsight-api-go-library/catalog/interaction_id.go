package catalog

type InteractionID interface {
	Protocol() Protocol
	Path() Path
	String() string
	MarshalText() ([]byte, error)
}
