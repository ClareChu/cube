package utils

func GetEventType(eventType string, eventTypes []string) (event string) {
	j := -1
	for i, event := range eventTypes {
		if eventType == event {
			j = i
		}
	}
	if j <= 0 || j == len(eventTypes) {
		return
	} else {
		event = eventTypes[j+1]
		return
	}
}
