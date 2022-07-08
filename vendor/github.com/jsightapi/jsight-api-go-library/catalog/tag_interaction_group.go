package catalog

type TagInteractionGroup interface {
	append(i InteractionId)
}

func newTagInteractionGroup(p Protocol) TagInteractionGroup {
	switch p {
	case JsonRpc:
		return newTagJsonRpcInteractionGroup()
	default: // case http:
		return newTagHttpInteractionGroup()
	}
}
