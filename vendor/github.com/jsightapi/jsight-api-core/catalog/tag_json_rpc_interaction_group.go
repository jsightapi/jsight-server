package catalog

type TagJsonRpcInteractionGroup struct {
	Protocol     Protocol        `json:"protocol"`
	Interactions []InteractionID `json:"interactions"`
}

func newTagJsonRpcInteractionGroup() *TagJsonRpcInteractionGroup {
	return &TagJsonRpcInteractionGroup{
		Protocol:     JsonRpc,
		Interactions: make([]InteractionID, 0, 5),
	}
}

func (l *TagJsonRpcInteractionGroup) append(i InteractionID) {
	l.Interactions = append(l.Interactions, i)
}
