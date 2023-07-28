package eventtypes

import (
	"fmt"
	"net/http"

	"github.com/goledgerdev/cc-tools/events"
)

// Description of a book
// var TestEvent = events.Event{
// 	Tag:         "testEvent",
// 	Label:       "Test Event",
// 	Description: "Test Event",
// 	Type:        events.EventTransaction,
// 	Transaction: "getNumberOfBooksFromLibrary",
// 	Receivers:   []string{"$org1MSP", "$org2MSP", "$orgMSP"},
// }

// var TestEvent = events.Event{
// 	Tag:         "testEvent",
// 	Label:       "Test Event",
// 	Description: "Test Event",
// 	Type:        events.EventLog,
// 	Receivers:   []string{"$org1MSP", "$org2MSP", "$orgMSP"},
// }

var TestEvent = events.Event{
	Tag:            "testEvent",
	Label:          "Test Event",
	Description:    "Test Event",
	Type:           events.EventCustom,
	CustomFunction: customFunc,
	Receivers:      []string{"$org1MSP", "$org2MSP", "$orgMSP"},
}

func customFunc(payload []byte) error {
	str := string(payload[:])

	link := fmt.Sprintf("http://3.215.180.194:8087/saveStr/%s", str)

	fmt.Printf("Sending request to %s\n", link)

	res, err := http.Get(link)
	fmt.Println(res)
	fmt.Println(err)
	if err != nil {
		return err
	}

	return nil
}
