package conn

type Message struct {
	EventName string        `json:"event_name"`
	EventData []interface{} `json:"event_data"`
}
