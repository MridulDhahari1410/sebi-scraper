package gatekeeper

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type tokenCacheModel struct {
	token     string
	cacheTime time.Time
}

type CachedS2STokenGenerator struct {
	DefaultTokenGenerator
	tokenCacheMap       map[string]tokenCacheModel
	bufferTimeInSeconds int64
	mu                  sync.Mutex
}

// dont create a token for caching with less than 5 seconds expiry
func NewS2STokenGeneratorWithCache(keyprovider KeyProvider, bufferInSecs int64) (TokenGenerator, error) {
	tokenGenerator := &DefaultTokenGenerator{
		keyProvider: keyprovider,
	}
	if bufferInSecs < defaultBufferCacheSeconds {
		return nil, errors.New(fmt.Sprintf("bufferInSecs is lesser than minimum buffer seconds : %d", defaultBufferCacheSeconds))
	}
	return &CachedS2STokenGenerator{DefaultTokenGenerator: *tokenGenerator, tokenCacheMap: make(map[string]tokenCacheModel),
		bufferTimeInSeconds: bufferInSecs, mu: sync.Mutex{}}, nil
}

func (d *CachedS2STokenGenerator) GenerateToken(customClaims AngelOneClaims, expiresInSec time.Duration) (string, error) {
	currentTime := time.Now()
	iat := currentTime.Add(-time.Second * issuedTimeTolerance)
	exp := currentTime.Add(expiresInSec)

	return d.GenerateTokenWithIssueTimeAndExpiryTime(customClaims, iat, exp)
}

func (d *CachedS2STokenGenerator) GenerateTokenWithIssueTimeAndExpiryTime(customClaims AngelOneClaims, iat time.Time, exp time.Time) (string, error) {
	if exp.Sub(time.Now()) < time.Duration(d.bufferTimeInSeconds)*time.Second {
		return empty, errors.New("expiryTime can not be less than buffer cache time")
	}

	// Check if the token is present in the cache and not expired
	if val, ok := d.tokenCacheMap[customClaims.KeyId]; ok && time.Now().Before(val.cacheTime) {
		return val.token, nil
	}

	if d.invalidUserType(customClaims) || customClaims.Issuer == empty || customClaims.KeyId == empty {
		return empty, errors.New("mandatory params are empty")
	}
	key, err := d.getPrivateKey(customClaims)
	if err != nil {
		return empty, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	//check again to avoid regeneration
	if val, ok := d.tokenCacheMap[customClaims.KeyId]; ok && time.Now().Before(val.cacheTime) {
		return val.token, nil
	}

	token, err := d.getToken(customClaims, exp, iat, key)
	if err != nil {
		return empty, err
	}

	//new map to avoid concurrency issues
	tokenCacheMapNew := make(map[string]tokenCacheModel)
	for k, val := range d.tokenCacheMap {
		tokenCacheMapNew[k] = val
	}
	tokenCacheMapNew[customClaims.KeyId] = tokenCacheModel{
		token:     token,
		cacheTime: exp.Add(time.Duration(-d.bufferTimeInSeconds) * time.Second),
	}

	d.tokenCacheMap = tokenCacheMapNew

	return token, nil
}
