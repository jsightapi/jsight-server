package catalog

type TagHTTPInteractionGroup struct {
	Protocol     Protocol        `json:"protocol"`
	Interactions []InteractionID `json:"interactions"`
}

func newTagHTTPInteractionGroup() *TagHTTPInteractionGroup {
	return &TagHTTPInteractionGroup{
		Protocol:     HTTP,
		Interactions: make([]InteractionID, 0, 5),
	}
}

func (l *TagHTTPInteractionGroup) append(i InteractionID) {
	l.Interactions = append(l.Interactions, i)
}
