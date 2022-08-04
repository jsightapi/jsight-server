package catalog

type Interaction interface {
	Path() Path
	appendTagName(TagName)
}
