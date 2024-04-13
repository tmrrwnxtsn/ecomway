package model

type CommissionType string

const (
	CommissionTypePercent  CommissionType = "percent"
	CommissionTypeFixed    CommissionType = "fixed"
	CommissionTypeCombined CommissionType = "combined"
	CommissionTypeText     CommissionType = "text"
)

type Commission struct {
	Type     CommissionType
	Currency string
	Percent  *float64
	Absolute *float64
	Message  map[string]string
}
