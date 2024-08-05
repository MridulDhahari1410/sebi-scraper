package gatekeeper

import (
	"errors"
	configs "github.com/angel-one/go-config-client"
)

type configProvider interface {
	getStringConfig(configKey string) (string, error)
	getIntConfig(configKey string) (int64, error)
	getBoolConfig(configKey string) (bool, error)

	getStringConfigD(configKey string, defaultValue string) string
	getIntConfigD(configKey string, defaultValue int64) int64
	getBoolConfigD(configKey string, defaultValue bool) bool

	getStringSecret(secretKey string) (string, error)

	refreshConfigs() error
}

func getCommonConfigProvider(applicationConfigClient configs.Client, provider int, params map[string]interface{}) (configProvider, error) {
	env := params[configEnvKey].(string)
	if provider == AWSAppConfig {
		secretName := globalLoginCredentials + "-" + env
		globalLoginParams := map[string]interface{}{
			configIDKey:              applicationName,
			configRegionKey:          params[configRegionKey],
			configAppKey:             applicationName,
			configEnvKey:             env,
			configTypeKey:            yaml,
			secretTypeKey:            yaml,
			configNamesKey:           []string{configTypeApplication},
			configCredentialsModeKey: params[configCredentialsModeKey],
			secretNamesKey:           []string{secretName},
		}
		if accessKeyIdValue, ok := params[accessKeyId]; ok {
			globalLoginParams[accessKeyId] = accessKeyIdValue
		}
		if secretKeyValue, ok := params[secretKey]; ok {
			globalLoginParams[secretKey] = secretKeyValue
		}
		configClient, err := configs.New(configs.Options{
			Provider: provider,
			Params:   globalLoginParams,
		})
		if err != nil {
			if applicationConfigClient == nil {
				return nil, err
			}
			return getApplicationConfigBasedCommonConfigProvider(applicationConfigClient, params)
		}
		return newConfigClientBasedConfigProvider(configClient, configTypeApplication, secretName)
	} else if provider == FileBased {
		consulBasedConfigProvider, err := newConsulBasedConfigProvider(configurations, secrets)
		if err != nil {
			return getFileBasedGlobalConfigProvider(applicationConfigClient, params)
		}
		return consulBasedConfigProvider, err
	}
	return nil, errors.New("incorrect provider type")
}

func getApplicationConfigBasedCommonConfigProvider(applicationConfigClient configs.Client, params map[string]interface{}) (configProvider, error) {

	commonConfigName := getGlobalLoginConfigName(params)
	commonSecretName := getGlobalLoginSecretName(params)
	return newConfigClientBasedConfigProvider(applicationConfigClient, commonConfigName, commonSecretName)
}

func getFileBasedGlobalConfigProvider(applicationConfigClient configs.Client, params map[string]interface{}) (configProvider, error) {
	commonConfigName := getGlobalLoginConfigName(params)
	commonSecretName := getGlobalLoginSecretName(params)
	if applicationConfigClient == nil {
		commonParams := map[string]interface{}{
			configsDirectoryKey: params[configsDirectoryKey],
			secretsDirectoryKey: params[secretsDirectoryKey],
			secretTypeKey:       yaml, //supporting yml for now we can read this also form options
			configTypeKey:       yaml,
			configNamesKey:      []string{commonConfigName},
			secretNamesKey:      []string{commonSecretName},
		}
		configClient, err := configs.New(configs.Options{
			Provider: configs.FileBased,
			Params:   commonParams,
		})
		if err != nil {
			return nil, err
		}
		applicationConfigClient = configClient
	}
	return newConfigClientBasedConfigProvider(applicationConfigClient, commonConfigName, commonSecretName)
}

func getGlobalLoginConfigName(params map[string]interface{}) string {
	val, ok := params[globalLoginConfigName]
	if ok {
		if valueString, isOk := val.(string); isOk {
			return valueString
		}
	}
	return defaultGlobalConfigName
}

func getGlobalLoginSecretName(params map[string]interface{}) string {
	val, ok := params[globalLoginSecretName]
	if ok {
		if valueString, isOk := val.(string); isOk {
			return valueString
		}
	}
	return defaultGlobalSecretName
}
