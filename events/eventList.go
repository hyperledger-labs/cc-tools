package events

// eventList is the list which should contain all defined events
var eventList = []Event{}

// EventList returns a copy of the eventList variable.
func EventList() []Event {
	listCopy := make([]Event, len(eventList))
	copy(listCopy, eventList)
	return listCopy
}

// InitEventList appends custom events to eventList to avoid initialization loop.
func InitEventList(l []Event) {
	eventList = l
}

// FetchEvent returns a pointer to the event object or nil if event is not found.
func FetchEvent(eventTag string) *Event {
	for _, event := range eventList {
		if event.Tag == eventTag {
			return &event
		}
	}
	return nil
}
