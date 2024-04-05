package model

type Method struct {
	ID             string
	DisplayedName  map[string]string
	ExternalSystem string
	ExternalMethod string
	Limits         Limits
	Commission     Commission
}
