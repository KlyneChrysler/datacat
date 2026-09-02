package gate

// refusal is the shape of one refusing response.
type refusal struct {
	event   string
	status  int
	message string
}
