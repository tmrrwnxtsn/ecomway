package model

type Method struct {
	ID             string
	DisplayedName  map[string]string
	ExternalSystem string
	ExternalMethod string
	Limits         map[string]Limits
	Commission     Commission
	IsFavorite     bool
}
