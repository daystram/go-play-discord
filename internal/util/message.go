package util

const MessageMaxLength = 900

func TrimMessage(msg string) string {
	if len(msg) > MessageMaxLength-3 {
		return msg[:MessageMaxLength] + "..."
	}

	return msg
}
