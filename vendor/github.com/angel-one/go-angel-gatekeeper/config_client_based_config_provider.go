package gatekeeper

import configs "github.com/angel-one/go-config-client"

type configClientBasedConfigProvider struct {
	configClient configs.Client
	baseConfig   string
	secretName   string
}

func (c *configClientBasedConfigProvider) refreshConfigs() error {
	return nil
}

func newConfigClientBasedConfigProvider(configClient configs.Client, baseConfigName string, secretName string) (configProvider, error) {
	return &configClientBasedConfigProvider{
		configClient: configClient,
		baseConfig:   baseConfigName,
		secretName:   secretName,
	}, nil
}

func (c *configClientBasedConfigProvider) getStringConfigD(configKey string, defaultValue string) string {
	return c.configClient.GetStringD(c.baseConfig, configKey, defaultValue)
}

func (c *configClientBasedConfigProvider) getIntConfigD(configKey string, defaultValue int64) int64 {
	return c.configClient.GetIntD(c.baseConfig, configKey, defaultValue)
}

func (c *configClientBasedConfigProvider) getBoolConfigD(configKey string, defaultValue bool) bool {
	return c.configClient.GetBoolD(c.baseConfig, configKey, defaultValue)
}

func (c *configClientBasedConfigProvider) getStringConfig(configKey string) (string, error) {
	return c.configClient.GetString(c.baseConfig, configKey)
}

func (c *configClientBasedConfigProvider) getIntConfig(configKey string) (int64, error) {
	return c.configClient.GetInt(c.baseConfig, configKey)
}

func (c *configClientBasedConfigProvider) getBoolConfig(configKey string) (bool, error) {
	return c.configClient.GetBool(c.baseConfig, configKey)
}

func (c *configClientBasedConfigProvider) getStringSecret(secretKey string) (string, error) {
	return c.configClient.GetStringSecret(c.secretName, secretKey)
}
