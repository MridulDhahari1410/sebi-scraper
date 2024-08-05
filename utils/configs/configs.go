package configs

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"sebi-scrapper/constants"
	"sebi-scrapper/utils/flags"

	configs "github.com/angel-one/go-config-client"
)

// Client is what holds the configs together with the flags.
type Client struct {
	configs.Client
	env           map[string]string
	configOptions configs.Options
}

// for test mode initialisation, put a sync once.
var o sync.Once

// Client is the instance of the config client to be used by the application.
var client *Client

// InitTestModeConfigs is used to initialize the configs from local repo.
func InitTestModeConfigs(directory string, configNames ...string) error {
	var err error
	var c configs.Client
	options := configs.Options{
		Provider: configs.FileBased,
		Params: map[string]any{
			constants.ConfigDirectoryKey:        directory,
			constants.ConfigNamesKey:            configNames,
			constants.ConfigTypeKey:             "yaml",
			constants.ConfigSecretsDirectoryKey: directory,
			constants.ConfigEnvKey:              flags.Env(),
		},
	}
	o.Do(func() {
		c, err = configs.New(options)
		if err == nil {
			client = getClient(c)
			client.configOptions = options
		}
		time.Sleep(5 * time.Second)
	})
	return err
}

// InitReleaseModeConfigs is used to initialize the configs using AWS.
func InitReleaseModeConfigs(configNames ...string) error {
	options := configs.Options{
		Provider: configs.AWSAppConfig,
		Params: map[string]any{
			constants.ConfigIDKey:              constants.ApplicationName,
			constants.ConfigRegionKey:          flags.AWSRegion(),
			constants.ConfigAccessKeyID:        flags.AWSAccessKeyID(),
			constants.ConfigSecretKey:          flags.AWSSecretAccessKey(),
			constants.ConfigAppKey:             constants.ApplicationName,
			constants.ConfigEnvKey:             flags.Env(),
			constants.ConfigTypeKey:            "yaml",
			constants.ConfigNamesKey:           configNames,
			constants.ConfigCredentialsModeKey: configs.AppConfigSharedCredentialMode,
		},
	}
	c, err := configs.New(options)
	if err != nil {
		return err
	}
	client = getClient(c)
	client.configOptions = options
	return nil
}

// Get is used to get the instance of the client.
func Get() *Client {
	return client
}

// GetClientOptions is used to get the config options.
func GetClientOptions() configs.Options {
	return client.configOptions
}

// GetStringWithEnv is used to get the config by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *Client) GetStringWithEnv(config, key string) (string, error) {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return s, err
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s, nil
}

// GetStringWithEnvD is used to get the config with default value by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *Client) GetStringWithEnvD(config, key, defaultValue string) string {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return defaultValue
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}

func getEnvironment() map[string]string {
	env := os.Environ()
	result := make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) >= 2 {
			result[s[0]] = strings.Join(s[1:], "=")
		}
	}
	return result
}

func getClient(c configs.Client) *Client {
	return &Client{Client: c, env: getEnvironment()}
}
