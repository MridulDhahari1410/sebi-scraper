package constants

// Flag constants.
const (
	ModeKey                    = "mode"
	ModeDefaultValue           = TestMode
	ModeUsage                  = "run mode of the application, can be test or release"
	EnvKey                     = "ENV"
	EnvDefaultValue            = "uat"
	EnvProd                    = "prod"
	PortKey                    = "port"
	PortDefaultValue           = 8080
	PortUsage                  = "application.yml port"
	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "resources/configs/uat"
	BaseConfigPathUsage        = "path to folder that stores your configurations"
	AWSRegionKey               = "AWS_REGION"
	AWSRegionDefaultValue      = "ap-south-1"
	AWSAccessKeyID             = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey         = "AWS_SECRET_ACCESS_KEY"
	InternalAccessTokenKey     = "INTERNAL_ACCESS_TOKEN"
	BaseURLKey                 = "BASE_URL"
	AWSBucketKey               = "AWS_BUCKET"
)
