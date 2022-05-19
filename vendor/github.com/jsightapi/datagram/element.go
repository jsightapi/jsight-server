package datagram

const (
	ellipsis     = "..."
	ellipsisSize = 3
)

type element struct {
	key         string
	value       string
	truncatable bool
}

// truncate tries to shorten the string by the specified number of characters.
// An ellipsis is added at the end of the shortened line.
func (e *element) truncate(n int) {
	if len(e.value) > ellipsisSize {
		ll := len(e.value) - n
		if ll > ellipsisSize {
			e.value = e.value[:ll-ellipsisSize] + ellipsis
		} else {
			e.value = ellipsis
		}
	}
	e.truncatable = false
}
