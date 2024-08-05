# angel-gatekeeper

Release : 0.12.7
https://angelbrokingpl.atlassian.net/wiki/spaces/Technology/pages/2815820201/Common+Authenticator

## Sample Repos

### Token generator 
#### with keyProvider
https://github.com/angel-one/amx-margin-service/blob/master/api/middleware/s2s_auth.go#L14

#### without keyProvider (using config client)
https://github.com/angel-one/login-trade/blob/master/utils/jwt_helper/jwt_helper.go#L67


### Authenticator
#### with configClient
https://github.com/angel-one/login-trade/blob/master/api/middleware/auth.go#L29

https://github.com/angel-one/bbe-aggregated-news/blob/master/middleware/authentication.go#L70

#### without configClient 
https://github.com/angel-one/amx-core/blob/master/api/middleware/auth.go#L55



