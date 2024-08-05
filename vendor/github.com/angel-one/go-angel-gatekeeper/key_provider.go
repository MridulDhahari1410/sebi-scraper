package gatekeeper

type KeyProvider interface {
	GetPrivateKey(keyName string) (string, error)
	GetPublicKey(issuer string, keyName string) (string, error)
}
