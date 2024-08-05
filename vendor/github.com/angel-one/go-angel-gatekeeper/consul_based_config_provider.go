package gatekeeper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	httpclient "github.com/angel-one/go-http-client"
	"github.com/spf13/cast"
	yamlv3 "gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type consulBasedConfigProvider struct {
	httpClient     *httpclient.Client
	configurations map[string]interface{}
	secrets        map[string]interface{}
	close          chan bool
	refetch        bool
	lock           sync.RWMutex
}

func newConsulBasedConfigProvider(configurationsKey string, secretsKey string) (configProvider, error) {

	headers := ConsulConfiguration[headers].(map[string]interface{})
	headers[consulTokenHeaderKey] = os.Getenv(consulToken)
	requestConfig := httpclient.NewRequestConfig(consulConfiguration, ConsulConfiguration).
		SetHeaderParams(headers)
	httpClient := initHttpClient(requestConfig)
	configMap, configReadErr := getConfiguration(httpClient, configurationsKey)
	if configReadErr != nil {
		return nil, configReadErr
	}
	secretMap, secretReadErr := getSecret(httpClient, secretsKey)
	if secretReadErr != nil {
		return nil, secretReadErr
	}

	configProvider := &consulBasedConfigProvider{
		httpClient:     httpClient,
		configurations: configMap,
		secrets:        secretMap,
	}
	refreshInterval, err := configProvider.getIntConfig(refreshIntervalInSec)

	if err != nil || refreshInterval == -1 {
		refreshInterval = defaultRefreshInterval
	}
	go configProvider.updateConfig(refreshInterval)
	return configProvider, nil
}

func getConfiguration(httpClient *httpclient.Client, configurationsKey string) (map[string]interface{}, error) {
	request := httpclient.NewRequest(consulConfiguration).
		SetURL(fmt.Sprintf("%s/%s", os.Getenv(consulUrl), configurationsKey))
	configMap, readErr := getValueFromConsul(httpClient, request)
	return configMap, readErr
}

func getSecret(httpClient *httpclient.Client, secretsKey string) (map[string]interface{}, error) {
	request := httpclient.NewRequest(consulConfiguration).
		SetURL(fmt.Sprintf("%s/%s", os.Getenv(consulUrl), secretsKey))
	secretMap, readErr := getValueFromConsul(httpClient, request)
	return secretMap, readErr
}

func getValueFromConsul(httpClient *httpclient.Client, request *httpclient.Request) (map[string]interface{}, error) {

	response, httpErr := httpClient.Request(request)
	if httpErr != nil {
		return nil, httpErr
	}
	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()
	if response.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("consul error : %d", response.StatusCode))
		return nil, err
	}
	body, bodyReadErr := ioutil.ReadAll(response.Body)
	if bodyReadErr != nil {
		return nil, bodyReadErr
	}
	var consulData []map[string]interface{}
	unmarshalErr := json.Unmarshal(body, &consulData)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	configuration, decodeErr := base64.StdEncoding.DecodeString(consulData[0]["Value"].(string))
	if decodeErr != nil {
		return nil, decodeErr
	}
	configMap := make(map[string]interface{})
	configReadErr := yamlv3.Unmarshal(configuration, &configMap)
	if configReadErr != nil {
		return nil, configReadErr
	}
	return configMap, nil
}

func initHttpClient(configs ...*httpclient.RequestConfig) *httpclient.Client {
	return httpclient.ConfigureHTTPClient(configs...)
}

func (c *consulBasedConfigProvider) getStringConfig(configKey string) (string, error) {
	val, err := c.getConfig(configKey)
	if val, ok := val.(string); ok {
		return val, nil
	}
	return "", err
}

func (c *consulBasedConfigProvider) getIntConfig(configKey string) (int64, error) {
	val, err := c.getConfig(configKey)
	if err != nil {
		return -1, err
	}
	return cast.ToInt64(val), nil
}

func (c *consulBasedConfigProvider) getBoolConfig(configKey string) (bool, error) {
	val, err := c.getConfig(configKey)
	if val, ok := val.(bool); ok {
		return val, nil
	}
	return false, err
}

func (c *consulBasedConfigProvider) getStringConfigD(configKey string, defaultValue string) string {
	value, err := c.getStringConfig(configKey)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *consulBasedConfigProvider) getIntConfigD(configKey string, defaultValue int64) int64 {
	value, err := c.getIntConfig(configKey)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *consulBasedConfigProvider) getBoolConfigD(configKey string, defaultValue bool) bool {
	value, err := c.getBoolConfig(configKey)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *consulBasedConfigProvider) getStringSecret(secretKey string) (string, error) {
	return c.getSecretValue(secretKey)
}

func (c *consulBasedConfigProvider) getSecretValue(key string) (string, error) {
	if value, ok := c.secrets[key]; ok {
		if stringValue, ok := value.(string); ok {
			return stringValue, nil
		}
	}
	return "", errors.New("either secret is not present or not string type")
}

func (c *consulBasedConfigProvider) getConfig(key string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	keySplits := strings.Split(key, ".")
	val := c.configurations[keySplits[0]]
	return c.readKey("", keySplits[1:], val)
}

func (c *consulBasedConfigProvider) readKey(key string, keySplits []string, val interface{}) (interface{}, error) {

	if len(keySplits) == 0 {
		if key == "" || key == "." {
			return val, nil
		}
		return nil, errors.New("key not found")
	}

	if m, ok := val.(map[string]interface{}); ok {
		newKey := key + keySplits[0]
		if v, ok := m[newKey]; ok {
			return c.readKey("", keySplits[1:], v)
		}
		return c.readKey(newKey+".", keySplits[1:], m)
	}
	return nil, errors.New("key not found")
}

func (c *consulBasedConfigProvider) refreshConfigs() error {
	return c.overWriteConfig(configurations)
}

func (c *consulBasedConfigProvider) isfetchNeeded() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.refetch
}

func (c *consulBasedConfigProvider) updateConfig(refreshIntervalInSec int64) {
	ticker := time.NewTicker(time.Duration(time.Second * time.Duration(refreshIntervalInSec)))
	for {
		select {
		case <-ticker.C:
			c.lock.Lock()
			c.refetch = true
			c.lock.Unlock()
		case <-c.close:
			ticker.Stop()
			return
		}
	}
}

func (c *consulBasedConfigProvider) overWriteConfig(configurationsKey string) error {
	if !c.isfetchNeeded() {
		return nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	configMap, configErr := getConfiguration(c.httpClient, configurationsKey)
	if configErr != nil {
		return fmt.Errorf("error reading config, err :  %w", configErr)
	}
	c.configurations = configMap
	c.refetch = false
	return nil
}
