package main

type Sensor struct {
	ID       int       `json:"id"`
	Text     string    `json:"Text"`
	Min      string    `json:"Min"`
	Value    string    `json:"Value"`
	Max      string    `json:"Max"`
	ImageUrl string    `json:"ImageUrl"`
	SensorId *string   `json:"SensorId"`
	Type     *string   `json:"Type"`
	Children []*Sensor `json:"Children"`
}
