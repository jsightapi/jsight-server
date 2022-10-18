package catalog

type TagInteractionGroup interface {
	append(i InteractionID)
}

func newTagInteractionGroup(p Protocol) TagInteractionGroup {
	switch p {
	case JsonRpc:
		return newTagJsonRpcInteractionGroup()
	default: // case http:
		return newTagHTTPInteractionGroup()
	}
}
