package messages

func Talk(msg string) string {
	switch msg {
	case "Hello":
		return "World"
	default:
		return "Sorry, I don't understand you"
	}
}
