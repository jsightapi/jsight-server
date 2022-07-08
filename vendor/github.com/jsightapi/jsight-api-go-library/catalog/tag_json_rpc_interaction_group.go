package catalog

type TagJsonRpcInteractionGroup struct {
	Protocol     Protocol        `json:"protocol"`
	Interactions []InteractionId `json:"interactions"`
}

func newTagJsonRpcInteractionGroup() *TagJsonRpcInteractionGroup {
	return &TagJsonRpcInteractionGroup{
		Protocol:     JsonRpc,
		Interactions: make([]InteractionId, 0, 5),
	}
}

func (l *TagJsonRpcInteractionGroup) append(i InteractionId) {
	l.Interactions = append(l.Interactions, i)
}
