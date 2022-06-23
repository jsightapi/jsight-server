package catalog

type InteractionId interface {
	Protocol() Protocol
	Path() Path
	String() string
	MarshalText() ([]byte, error)
}
