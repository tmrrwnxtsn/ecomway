package config

type YooKassaConfig struct {
	API      YooKassaAPIConfig                `yaml:"api"`
	Channels map[string]YooKassaChannelConfig `yaml:"channels"`
	Methods  MethodsConfig                    `yaml:"methods"`
}

type YooKassaAPIConfig struct {
	BaseURL   string `yaml:"base_url"`
	ShopID    string `yaml:"shop_id"`
	SecretKey string `yaml:"secret_key"`
}

type YooKassaChannelConfig struct {
	Code                      string `yaml:"code"`
	PaymentMethodType         string `yaml:"payment_method_type"`
	PaymentTimeoutToFailedMin int    `yaml:"payment_timeout_to_failed_min"`
}
