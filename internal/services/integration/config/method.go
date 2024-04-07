package config

type MethodsConfig struct {
	Payment []MethodConfig `yaml:"payment"`
	Payout  []MethodConfig `yaml:"payout"`
}

type MethodConfig struct {
	ID             string                        `yaml:"id"`
	DisplayedName  map[string]string             `yaml:"name"`
	ExternalSystem string                        `yaml:"external_system"`
	ExternalMethod string                        `yaml:"external_method"`
	Limits         map[string]MethodLimitsConfig `yaml:"limits"`
	Commission     MethodCommissionConfig        `yaml:"commission"`
}

type MethodLimitsConfig struct {
	MinAmount float64 `yaml:"min_amount"`
	MaxAmount float64 `yaml:"max_amount"`
}

type MethodCommissionConfig struct {
	Type     string            `yaml:"type"`
	Currency *string           `yaml:"currency"`
	Percent  *float64          `yaml:"percent"`
	Absolute *float64          `yaml:"absolute"`
	Message  map[string]string `yaml:"message"`
}
