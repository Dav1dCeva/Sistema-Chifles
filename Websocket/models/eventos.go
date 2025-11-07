package models

type EventMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}
