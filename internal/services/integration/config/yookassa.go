package config

type YooKassaConfig struct {
	API     YooKassaAPIConfig `yaml:"api"`
	Methods MethodsConfig     `yaml:"methods"`
}

type YooKassaAPIConfig struct {
}
