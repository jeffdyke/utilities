package cloudwatch

import (
	"encoding/json"
	"log"
)
const (
	SuricataFilter = `{ $.event_type = alert && $.alert.action = allowed && $.alert.signature_id!= 2013504 && $.alert.signature_id!= 2221002 && $.http.http_method!= PROXY}`
)
type IndexedSuricataAlert struct {
	Alert SuricataAlert `json:"alert"`
	Count uint32 `json:"count"`
}
type IndexedAlert = map[uint32]IndexedSuricataAlert

func SuricataEvents(startEnd StartEndFilter) IndexedAlert {
	var configs []LogConfig
	configs = append(configs, LogConfig{
		LogGroup:  "StagingSuricataIPS",
		LogPrefix: "staging",
	})
	configs = append(configs, LogConfig{
		LogGroup:  "ProductionSuricataIPS",
		LogPrefix: "prod",
	})

	flSlice := FilterList(configs, startEnd)
	var events []SuricataEvent
	for _, filter := range flSlice {
		events = append(events, FindEvents(filter)...)
	}
	return Aggregate(events)

}


func Aggregate(events []SuricataEvent) IndexedAlert {
	var agg = make(IndexedAlert)

	for _, event := range events {
		val, ok := agg[event.Alert.SignatureId]
		if ok {
			val.Count++
		} else {
			agg[event.Alert.SignatureId] = IndexedSuricataAlert{Count: 1, Alert:event.Alert}
		}
	}
	return agg
}

func FindEvents(f Filter) []SuricataEvent {
	filtered := f.FilterLogs()
	var swEvents []SuricataEvent
	for _, event := range filtered {
		var sEvent SuricataEvent
		data := []byte(*event.Message)
		err := json.Unmarshal(data, &sEvent)
		if err != nil {
			log.Fatalf("Failed to unmarshal %v\n", err)
		}
		swEvents = append(swEvents, sEvent)
	}
	return swEvents
}

