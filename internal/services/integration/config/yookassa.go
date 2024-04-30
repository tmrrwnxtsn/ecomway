package config

const (
	yooKassaPaymentsSecretKeyEnvKey = "YOOKASSA_PAYMENTS_SECRET_KEY"
	yooKassaPayoutsSecretKeyEnvKey  = "YOOKASSA_PAYOUTS_SECRET_KEY"
)

type YooKassaConfig struct {
	API      YooKassaAPIConfig                `yaml:"api"`
	Channels map[string]YooKassaChannelConfig `yaml:"channels"`
	Methods  MethodsConfig                    `yaml:"methods"`
}

type YooKassaAPIConfig struct {
	BaseURL  string                    `yaml:"base_url"`
	Payments YooKassaAPIPaymentsConfig `yaml:"payments"`
	Payouts  YooKassaAPIPayoutsConfig  `yaml:"payouts"`
}

type YooKassaAPIPaymentsConfig struct {
	ShopID    string `yaml:"shop_id"`
	SecretKey string
}

type YooKassaAPIPayoutsConfig struct {
	AgentID   string `yaml:"agent_id"`
	SecretKey string
}

type YooKassaChannelConfig struct {
	Code                      string `yaml:"code"`
	PaymentMethodType         string `yaml:"payment_method_type"`
	PaymentTimeoutToFailedMin int    `yaml:"payment_timeout_to_failed_min"`
}
