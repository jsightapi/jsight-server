package catalog

type Interaction interface {
	Path() Path
	// SetPathVariables(*PathVariables)
}
