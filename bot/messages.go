package bot

func Talk(msg string) string {
	switch msg {
	case "Hello":
		return "World"
	case "Live":
		return "Chat?"
	case "Grzegorz":
		return "Gajda"
	default:
		return "Sorry, I don't understand you"
	}
}
