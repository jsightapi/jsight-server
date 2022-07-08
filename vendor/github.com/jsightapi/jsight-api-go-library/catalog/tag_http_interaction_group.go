package catalog

type TagHttpInteractionGroup struct {
	Protocol     Protocol        `json:"protocol"`
	Interactions []InteractionId `json:"interactions"`
}

func newTagHttpInteractionGroup() *TagHttpInteractionGroup {
	return &TagHttpInteractionGroup{
		Protocol:     HTTP,
		Interactions: make([]InteractionId, 0, 5),
	}
}

func (l *TagHttpInteractionGroup) append(i InteractionId) {
	l.Interactions = append(l.Interactions, i)
}
