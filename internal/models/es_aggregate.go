package models

type EsAggregate struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`
}
