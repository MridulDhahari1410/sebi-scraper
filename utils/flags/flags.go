package flags

import (
	"os"

	"sebi-scrapper/constants"

	flag "github.com/spf13/pflag"
)

var (
	mode           = flag.String(constants.ModeKey, constants.ModeDefaultValue, constants.ModeUsage)
	port           = flag.Int(constants.PortKey, constants.PortDefaultValue, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.BaseConfigPathDefaultValue,
		constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

// Mode is the run mode, can be either test or release.
func Mode() string {
	return *mode
}

// Env is the runtime environment.
func Env() string {
	env := os.Getenv(constants.EnvKey)
	if env == "" {
		return constants.EnvDefaultValue
	}
	return env
}

// Port is the application.yml port number where the process will be started.
func Port() int {
	return *port
}

// BaseConfigPath is the path that holds the configuration files.
func BaseConfigPath() string {
	return *baseConfigPath
}

// AWSRegion is the region where the application is running.
func AWSRegion() string {
	region := os.Getenv(constants.AWSRegionKey)
	if region == "" {
		return constants.AWSRegionDefaultValue
	}
	return region
}

// AWSAccessKeyID is the access key id for aws.
func AWSAccessKeyID() string {
	return os.Getenv(constants.AWSAccessKeyID)
}

// AWSSecretAccessKey is the secret access key for aws.
func AWSSecretAccessKey() string {
	return os.Getenv(constants.AWSSecretAccessKey)
}

// AWSBucket is the bucket associated with this application.
func AWSBucket() string {
	return os.Getenv(constants.AWSBucketKey)
}

// InternalAccessToken is the token to access internal APIs.
func InternalAccessToken() string {
	return os.Getenv(constants.InternalAccessTokenKey)
}

// BaseURL is used as the base url of the service.
func BaseURL() string {
	return os.Getenv(constants.BaseURLKey)
}
