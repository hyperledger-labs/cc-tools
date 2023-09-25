package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/hyperledger-labs/cc-tools/events"
)

func TestFetchEvent(t *testing.T) {
	event := *events.FetchEvent("createLibraryLog")
	expectedEvent := testEventTypeList[0]

	if !reflect.DeepEqual(event, expectedEvent) {
		log.Println("these should be deeply equal")
		log.Println(event)
		log.Println(expectedEvent)
		t.FailNow()
	}
}

func TestFetchEventList(t *testing.T) {
	eventList := events.EventList()
	expectedEventList := testEventTypeList

	if !reflect.DeepEqual(eventList, expectedEventList) {
		log.Println("these should be deeply equal")
		log.Println(eventList)
		log.Println(expectedEventList)
		t.FailNow()
	}
}
