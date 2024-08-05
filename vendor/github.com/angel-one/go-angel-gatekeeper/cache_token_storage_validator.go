package gatekeeper

import (
	"context"
	"errors"
	"fmt"
	cache "github.com/angel-one/go-cache-client"
	"runtime"
	"strings"
	"time"
)

// todo implement cache
type cacheBasedTokenStorageValidator struct {
	configProvider configProvider
	keyValidator   SessionValidator
	client         cache.Client
	logger         func(ctx context.Context, logMessage LogMessage)
}

func newCacheBasedTokenStorageValidator(provider configProvider, logger func(ctx context.Context, logMessage LogMessage)) (tokenStorageValidator, error) {

	cacheClient, err := getCacheClient(provider, logger)
	if err != nil {
		return nil, err
	}
	return &cacheBasedTokenStorageValidator{
		configProvider: provider,
		keyValidator:   nil,
		logger:         logger,
		client:         cacheClient,
	}, nil
}

func (receiver *cacheBasedTokenStorageValidator) validate(loginClaims *JWTTokenClaims) (bool, error) {
	key := receiver.keyValidator.GetSessionKey(loginClaims)
	ctx, cancel := context.WithTimeout(context.TODO(), getContextTimeout(receiver.configProvider))
	defer cancel()
	value, err := receiver.client.Get(ctx, key)
	return receiver.keyValidator.Validate(value, err, loginClaims), nil
}

func (receiver *cacheBasedTokenStorageValidator) setSessionValidator(validator SessionValidator) {
	receiver.keyValidator = validator
}

func getCacheClient(provider configProvider, logger func(ctx context.Context, logMessage LogMessage)) (cache.Client, error) {

	cacheType := provider.getStringConfigD(fmt.Sprintf("%s.%s", persistence, cacheType), defaultCacheType)

	switch cacheType {
	case inmemory:
		return nil, errors.New("in-memory cache is not supported")
	case redis:
		return getRedisClient(provider, logger)
	default:
		return getRedisClient(provider, logger)
	}
}

func getRedisClient(provider configProvider, logger func(ctx context.Context, logMessage LogMessage)) (cache.Client, error) {

	username := provider.getStringConfigD(fmt.Sprintf("%s.%s", persistence, cacheUsername), empty)
	password, err := provider.getStringSecret(fmt.Sprintf("%s.%s", persistence, cachePassword))
	if err != nil {
		return nil, err
	}

	clusterModeEnabled := provider.getBoolConfigD(fmt.Sprintf("%s.%s", persistence, clusterModeEnabled), false)

	var cacheLogger func(ctx context.Context, name string, latencyInMillis int64, err error)
	if logger != nil {
		cacheLogger = func(ctx context.Context, name string, latencyInMillis int64, err error) {
			logLevel := InfoLevel
			if err != nil {
				logLevel = WarnLevel
			}
			logger(ctx, LogMessage{
				LogLevel: logLevel,
				Message:  fmt.Sprintf("%s with latency %d", name, latencyInMillis),
			})
		}
	}

	if clusterModeEnabled {
		return getClusterClient(provider, username, password, cacheLogger)
	}
	return getNonClusterClient(provider, username, password, cacheLogger)
}

func getNonClusterClient(provider configProvider, username string, password string, cacheLogger func(ctx context.Context, name string, latencyInMillis int64, err error)) (cache.Client, error) {
	params := getRedisParams(provider, username, password)
	addressConfig, err := provider.getStringConfig(fmt.Sprintf("%s.%s", persistence, address))
	if err != nil {
		return nil, err
	}
	params[address] = addressConfig
	params[db] = provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, db), defaultDb)
	return cache.New(context.TODO(), cache.Options{
		Provider: cache.Redis,
		Params:   params,
		Logger:   cacheLogger,
	})
}

func getClusterClient(provider configProvider, username string, password string, cacheLogger func(ctx context.Context, name string, latencyInMillis int64, err error)) (cache.Client, error) {
	params := getRedisParams(provider, username, password)
	address, err := provider.getStringConfig(fmt.Sprintf("%s.%s", persistence, addresses))
	if err != nil {
		return nil, err
	}
	params[addresses] = strings.Split(address, ",")
	return cache.New(context.TODO(), cache.Options{
		Provider: cache.RedisCluster,
		Params:   params,
		Logger:   cacheLogger,
	})
}

func getRedisParams(provider configProvider, username string, password string) map[string]interface{} {
	return map[string]interface{}{
		cacheUsername: username,
		cachePassword: password,
		dialTimeout:   getDialTimeout(provider),
		readTimeout:   getReadTimeout(provider),
		poolTimeout:   getPoolTimeout(provider),
		// pool config
		poolSize:              provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, poolSize), int64(5*(runtime.GOMAXPROCS(0)+1))),
		minIdleConnections:    provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, minIdleConnections), defaultIdleConnections),
		maxConnectionAge:      getMaxConnectionAge(provider),
		idleConnectionTimeout: getIdleConnectionTimeout(provider),
		// retry config
		maxRetries:      provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, maxRetries), defaultRetries),
		minRetryBackoff: getMinRetryDuration(provider),
		maxRetryBackoff: getMaxRetryDuration(provider),
	}
}

func getMinRetryDuration(provider configProvider) interface{} {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, minRetryBackoff), defaultMinRetryBackOff)) * time.Millisecond
}

func getMaxRetryDuration(provider configProvider) interface{} {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, maxRetryBackoff), defaultMaxRetryBackOff)) * time.Millisecond
}

func getMaxConnectionAge(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, maxConnectionAge), defaultConnectionAge)) * time.Millisecond
}

func getDialTimeout(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, dialTimeout), defaultDialTimeout)) * time.Millisecond
}

func getReadTimeout(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, readTimeout), defaultReadTimeout)) * time.Millisecond
}

func getPoolTimeout(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, poolTimeout), defaultPoolTimeout)) * time.Millisecond
}

func getIdleConnectionTimeout(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, idleConnectionTimeout), defaultIdleConnectionTimeout)) * time.Millisecond
}

func getContextTimeout(provider configProvider) time.Duration {
	return time.Duration(provider.getIntConfigD(fmt.Sprintf("%s.%s", persistence, contextTimeout), defaultContextTimeout)) * time.Millisecond
}
